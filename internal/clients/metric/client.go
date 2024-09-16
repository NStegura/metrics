package metric

import (
	"bytes"
	"compress/gzip"
	"crypto/hmac"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/NStegura/metrics/utils/ip"

	rsaKeys "github.com/NStegura/metrics/utils/rsa"

	"github.com/sirupsen/logrus"
)

// Client - клиент к хранению метрик.
type Client struct {
	client       *http.Client
	logger       *logrus.Logger
	URL          string
	bodyHashKey  string
	cryptoKey    *rsa.PublicKey
	compressType string
}

func New(
	addr string,
	bodyHashKey string,
	cryptoKeyPath string,
	logger *logrus.Logger,
) (*Client, error) {
	var (
		cryptoKey *rsa.PublicKey
		err       error
	)
	if !strings.HasPrefix(addr, "http") {
		addr, err = url.JoinPath("http:", addr)
		if err != nil {
			return nil, fmt.Errorf("failed to init client, %w", err)
		}
	}
	if cryptoKeyPath != "" {
		cryptoKey, err = rsaKeys.ReadPublicKey(cryptoKeyPath)
		if err != nil {
			return nil, fmt.Errorf("failed to read public key: %w", err)
		}
	}
	return &Client{
		client:       &http.Client{},
		URL:          addr,
		bodyHashKey:  bodyHashKey,
		cryptoKey:    cryptoKey,
		logger:       logger,
		compressType: "gzip",
	}, nil
}

type RequestError struct {
	URL        *url.URL
	Body       []byte
	StatusCode int
}

func (e RequestError) Error() string {
	return fmt.Sprintf(
		"Metric request error: url=%s, code=%v, body=%s",
		e.URL, e.StatusCode, e.Body,
	)
}

func NewRequestError(response *http.Response) error {
	body, err := io.ReadAll(response.Body)
	if err != nil {
		return fmt.Errorf("failed to read metric client resp: %w", err)
	}
	return &RequestError{response.Request.URL, body, response.StatusCode}
}

// UpdateGaugeMetric обновляет gauge метрику.
func (c *Client) UpdateGaugeMetric(name string, value float64) error {
	resp, err := c.post(
		fmt.Sprintf("%s/update/gauge/%s/%v", c.URL, name, value),
		"text/plain",
		nil,
	)
	if err != nil {
		return err
	}
	defer func() {
		err = resp.Body.Close()
		if err != nil {
			c.logger.Error(err)
		}
	}()

	if resp.StatusCode != http.StatusOK {
		return NewRequestError(resp)
	}
	return nil
}

// UpdateCounterMetric обновляет counter метрику.
func (c *Client) UpdateCounterMetric(name string, value int64) error {
	resp, err := c.post(
		fmt.Sprintf("%s/update/counter/%s/%v", c.URL, name, value),
		"text/plain",
		nil,
	)
	if err != nil {
		return err
	}
	defer func() {
		err = resp.Body.Close()
		if err != nil {
			c.logger.Error(err)
		}
	}()

	if resp.StatusCode != http.StatusOK {
		return NewRequestError(resp)
	}
	return nil
}

// UpdateMetric обновляет метрику.
func (c *Client) UpdateMetric(jsonBody []byte) error {
	resp, err := c.post(
		fmt.Sprintf("%s/update/", c.URL),
		"application/json",
		jsonBody,
	)
	if err != nil {
		return err
	}
	defer func() {
		err = resp.Body.Close()
		if err != nil {
			c.logger.Error(err)
		}
	}()

	if resp.StatusCode != http.StatusOK {
		return NewRequestError(resp)
	}
	return nil
}

// UpdateMetrics обновляет набор метрик.
func (c *Client) UpdateMetrics(metrics []Metrics) error {
	if len(metrics) == 0 {
		c.logger.Info("Empty metric result")
		return nil
	}

	jsonBody, err := json.Marshal(metrics)
	if err != nil {
		return fmt.Errorf("failed to decode metrics, err %w", err)
	}

	resp, err := c.post(
		fmt.Sprintf("%s/updates/", c.URL),
		"application/json",
		jsonBody,
	)
	if err != nil {
		return err
	}
	defer func() {
		err = resp.Body.Close()
		if err != nil {
			c.logger.Error(err)
		}
	}()

	if resp.StatusCode != http.StatusOK {
		return NewRequestError(resp)
	}
	return nil
}

func (c *Client) post(
	url string,
	contentType string,
	body []byte,
) (resp *http.Response, err error) {
	headers := make(map[string]string)

	if c.bodyHashKey != "" {
		h := hmac.New(sha256.New, []byte(c.bodyHashKey))
		_, err = h.Write(body)
		if err != nil {
			return nil, fmt.Errorf("failed to write body hash: %w", err)
		}
		headers["HashSHA256"] = hex.EncodeToString(h.Sum(nil))
	}

	if c.cryptoKey != nil {
		body, err = rsaKeys.EncryptOAEP(sha256.New(), rand.Reader, c.cryptoKey, body, nil)
		if err != nil {
			return nil, fmt.Errorf("failed to encrypt key, %w", err)
		}
	}

	if c.compressType == "gzip" {
		body, err = compress(body)
		if err != nil {
			return nil, err
		}
		headers["Accept-Encoding"] = c.compressType
		headers["Content-Encoding"] = c.compressType
	}
	headers["Content-Type"] = contentType

	selfIP, err := ip.GetIP()
	if err != nil {
		return nil, fmt.Errorf("failed to get ip: %w", err)
	}
	headers["X-Real-IP"] = selfIP

	bodyReader := bytes.NewReader(body)
	req, err := http.NewRequest(http.MethodPost, url, bodyReader)

	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	for h, v := range headers {
		req.Header.Set(h, v)
	}

	resp, err = c.doWithRetry(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	return resp, nil
}

func (c *Client) doWithRetry(req *http.Request) (resp *http.Response, err error) {
	for _, backoff := range c.scheduleBackoffAttempts() {
		resp, err = c.doWithLog(req)
		if err != nil {
			c.logger.Warningf("Request failed %v, err: %+v", req, err)
			return
		}

		if resp.StatusCode < http.StatusInternalServerError {
			break
		}
		time.Sleep(backoff)
		c.logger.Warningf("Retrying in %v, Request error: %+v", backoff, err)
	}
	return
}

func (c *Client) doWithLog(req *http.Request) (resp *http.Response, err error) {
	start := time.Now()
	resp, err = c.client.Do(req)
	duration := time.Since(start)

	if err != nil {
		c.logger.WithFields(logrus.Fields{
			"uri":      req.URL.Path,
			"method":   req.Method,
			"duration": duration,
			"err":      err,
		}).Warning()
		return
	}
	c.logger.WithFields(logrus.Fields{
		"uri":      req.URL.Path,
		"method":   req.Method,
		"status":   resp.Status,
		"duration": duration,
	}).Info()
	return
}

func (c *Client) scheduleBackoffAttempts() []time.Duration {
	return []time.Duration{
		1 * time.Second,
		3 * time.Second,
		5 * time.Second,
	}
}

func compress(data []byte) ([]byte, error) {
	var b bytes.Buffer
	w := gzip.NewWriter(&b)
	_, err := w.Write(data)
	if err != nil {
		return nil, fmt.Errorf("failed write data to compress buffer: %w", err)
	}
	err = w.Close()
	if err != nil {
		return nil, fmt.Errorf("failed compress data: %w", err)
	}
	return b.Bytes(), nil
}
