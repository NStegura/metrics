package base

import (
	"bytes"
	"compress/gzip"
	"crypto/hmac"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"time"

	rsaKeys "github.com/NStegura/metrics/internal/utils/rsa"

	"github.com/sirupsen/logrus"
)

//nolint:govet // unexpected?
type BaseClient struct {
	BodyHashKey  string
	CompressType string
	retryPolicy  []time.Duration
	isRetryable  func(result any, err error) bool
	CryptoKey    *rsa.PublicKey
	Logger       *logrus.Logger
}

func NewBaseClient(options ...Option) (*BaseClient, error) {
	c := &BaseClient{
		Logger:       logrus.New(),
		CompressType: "gzip",
		retryPolicy:  nil,
	}
	for _, opt := range options {
		err := opt(c)
		if err != nil {
			return nil, err
		}
	}
	return c, nil
}

func (c *BaseClient) Compress(data []byte) ([]byte, bool, error) {
	if c.CompressType == "" {
		return nil, false, nil
	}

	var b bytes.Buffer
	w := gzip.NewWriter(&b)
	_, err := w.Write(data)
	if err != nil {
		return nil, true, fmt.Errorf("failed to write data to compress buffer: %w", err)
	}
	err = w.Close()
	if err != nil {
		return nil, true, fmt.Errorf("failed to compress data: %w", err)
	}
	return b.Bytes(), true, nil
}

func (c *BaseClient) GenerateHMAC(body []byte) (string, bool, error) {
	if c.BodyHashKey == "" {
		return "", false, nil
	}

	h := hmac.New(sha256.New, []byte(c.BodyHashKey))
	_, err := h.Write(body)
	if err != nil {
		return "", true, fmt.Errorf("failed to write body hash: %w", err)
	}
	return hex.EncodeToString(h.Sum(nil)), true, nil
}

func (c *BaseClient) Encrypt(body []byte) ([]byte, bool, error) {
	if c.CryptoKey == nil {
		return nil, false, nil
	}

	body, err := rsaKeys.EncryptOAEP(sha256.New(), rand.Reader, c.CryptoKey, body, nil)
	if err != nil {
		return nil, true, fmt.Errorf("failed to encrypt key, %w", err)
	}
	return body, true, nil
}

func (c *BaseClient) Execute(
	doFunc func() (any, error),
	path string,
	method string,
) (any, error) {
	if c.retryPolicy != nil {
		return c.executeWithRetry(doFunc, path, method)
	}
	return c.execute(doFunc, path, method)
}

func (c *BaseClient) execute(
	doFunc func() (any, error),
	path string,
	method string,
) (result any, err error) {
	start := time.Now()
	result, err = doFunc()
	duration := time.Since(start)
	if err != nil {
		c.Logger.WithFields(logrus.Fields{
			"uri":      path,
			"method":   method,
			"duration": duration,
			"err":      err,
		}).Warning()
		return
	}
	c.Logger.WithFields(logrus.Fields{
		"uri":      path,
		"method":   method,
		"duration": duration,
		"status":   "Ok",
	}).Info()
	return
}

func (c *BaseClient) executeWithRetry(
	doFunc func() (any, error),
	path string,
	method string,
) (result any, err error) {
	for _, backoff := range c.retryPolicy {
		result, err = c.execute(doFunc, path, method)
		if err == nil || !c.isRetryable(result, err) {
			return
		}
		c.Logger.Warningf("Retrying in %v, error: %+v", backoff, err)
		time.Sleep(backoff)
	}
	return result, fmt.Errorf("failed after retries: %w", err)
}
