package metric

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

type client struct {
	client *http.Client
	URL    string
}

func New(url string) *client {
	return &client{
		client: &http.Client{},
		URL:    fmt.Sprintf("http://%s", url),
	}
}

type RequestError struct {
	URL        *url.URL
	StatusCode int
	Body       []byte
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
		return err
	}
	return &RequestError{response.Request.URL, response.StatusCode, body}
}

func (c *client) UpdateGaugeMetric(name string, value float64, compressType string) error {
	resp, err := c.Post(
		fmt.Sprintf("%s/update/gauge/%s/%v", c.URL, name, value),
		"text/plain",
		nil,
		compressType,
	)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return NewRequestError(resp)
	}
	return nil
}

func (c *client) UpdateCounterMetric(name string, value int64, compressType string) error {
	resp, err := c.Post(
		fmt.Sprintf("%s/update/counter/%s/%v", c.URL, name, value),
		"text/plain",
		nil,
		compressType,
	)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return NewRequestError(resp)
	}
	return nil
}

func (c *client) UpdateMetric(jsonBody []byte, compressType string) error {
	resp, err := c.Post(
		fmt.Sprintf("%s/update/", c.URL),
		"application/json",
		jsonBody,
		compressType,
	)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return NewRequestError(resp)
	}
	return nil
}

func (c *client) Post(url string, contentType string, body []byte, compressType string) (resp *http.Response, err error) {
	headers := make(map[string]string, 10)

	if compressType == "gzip" {
		body, err = compress(body)
		if err != nil {
			return nil, err
		}
		headers["Accept-Encoding"] = compressType
		headers["Content-Encoding"] = compressType
	}
	headers["Content-Type"] = contentType

	bodyReader := bytes.NewReader(body)
	req, err := http.NewRequest("POST", url, bodyReader)

	if err != nil {
		return nil, err
	}

	for h, v := range headers {
		req.Header.Set(h, v)
	}
	return c.client.Do(req)
}

func compress(data []byte) ([]byte, error) {
	var b bytes.Buffer
	w := gzip.NewWriter(&b)
	_, err := w.Write(data)
	if err != nil {
		return nil, fmt.Errorf("failed write data to compress buffer: %v", err)
	}
	err = w.Close()
	if err != nil {
		return nil, fmt.Errorf("failed compress data: %v", err)
	}
	return b.Bytes(), nil
}
