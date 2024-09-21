package metric

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/NStegura/metrics/internal/clients/base"
	"github.com/NStegura/metrics/utils/ip"
)

// Client - клиент к хранению метрик.
type Client struct {
	*base.BaseClient
	client *http.Client
	URL    string
}

func NewHTTPClient(addr string, options ...base.Option) (*Client, error) {
	var err error
	if !strings.HasPrefix(addr, "http") {
		addr, err = url.JoinPath("http:", addr)
		if err != nil {
			return nil, fmt.Errorf("failed to init client, %w", err)
		}
	}
	bc, err := base.NewBaseClient(options...)
	if err != nil {
		return nil, fmt.Errorf("failed to init client: %w", err)
	}
	return &Client{
		BaseClient: bc,
		client:     &http.Client{},
		URL:        addr,
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
			c.Logger.Error(err)
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
			c.Logger.Error(err)
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
			c.Logger.Error(err)
		}
	}()

	if resp.StatusCode != http.StatusOK {
		return NewRequestError(resp)
	}
	return nil
}

// UpdateMetrics обновляет набор метрик.
func (c *Client) UpdateMetrics(_ context.Context, metrics []Metrics) error {
	if len(metrics) == 0 {
		c.Logger.Info("Empty metric result")
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
			c.Logger.Error(err)
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

	body, headers, err := c.prepareRequest(contentType, body)
	if err != nil {
		return nil, err
	}

	bodyReader := bytes.NewReader(body)
	req, err := http.NewRequest(http.MethodPost, url, bodyReader)

	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	for h, v := range headers {
		req.Header.Set(h, v)
	}

	execute, err := c.Execute(
		func() (any, error) {
			return c.client.Do(req)
		},
		req.URL.Path,
		req.Method,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	resp = execute.(*http.Response)
	return resp, nil
}

func (c *Client) prepareRequest(
	contentType string,
	body []byte,
) ([]byte, map[string]string, error) {
	var (
		included      bool
		hash          string
		encryptedBody []byte
		err           error
	)

	headers := make(map[string]string)
	headers["Content-Type"] = contentType

	if hash, included, err = c.GenerateHMAC(body); included {
		if err != nil {
			return nil, nil, err
		}
		headers["HashSHA256"] = hash
	}

	if encryptedBody, included, err = c.Encrypt(body); included {
		if err != nil {
			return nil, nil, err
		}
		body = encryptedBody
	}

	if body, included, err = c.Compress(body); included {
		if err != nil {
			return nil, nil, err
		}
		headers["Accept-Encoding"] = c.CompressType
		headers["Content-Encoding"] = c.CompressType
	}

	selfIP, err := ip.GetIP()
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get ip: %w", err)
	}
	headers["X-Real-IP"] = selfIP
	return body, headers, nil
}
