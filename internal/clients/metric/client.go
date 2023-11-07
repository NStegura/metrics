package metric

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
)

type client struct {
	client *http.Client
	URL    string
}

func New() *client {
	return &client{
		client: &http.Client{},
		URL:    "http://localhost:8080",
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

func (c *client) UpdateGaugeMetric(name string, value float64) error {
	resp, err := c.client.Post(
		fmt.Sprintf("%s/update/gauge/%s/%v", c.URL, name, value),
		"text/plain",
		nil,
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

func (c *client) UpdateCounterMetric(name string, value int64) error {
	resp, err := c.client.Post(
		fmt.Sprintf("%s/update/counter/%s/%v", c.URL, name, value),
		"text/plain",
		nil,
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
