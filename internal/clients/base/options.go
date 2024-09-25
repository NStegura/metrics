package base

import (
	"crypto/rsa"
	"fmt"
	"net/http"
	"time"

	rsaKeys "github.com/NStegura/metrics/internal/utils/rsa"

	"github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Option func(*BaseClient) error

// WithLogger Опция для настройки логгера.
func WithLogger(logger *logrus.Logger) Option {
	return func(c *BaseClient) error {
		c.Logger = logger
		return nil
	}
}

// WithRetryPolicy Опция для настройки политики ретраев.
func WithRetryPolicy(
	backoffs []time.Duration,
	isRetryable func(result any, err error) bool,
) Option {
	return func(c *BaseClient) error {
		c.retryPolicy = backoffs
		c.isRetryable = isRetryable
		return nil
	}
}

// WithBodyHashKey Опция для настройки ключа хэширования тела.
func WithBodyHashKey(key string) Option {
	return func(c *BaseClient) error {
		c.BodyHashKey = key
		return nil
	}
}

// WithCryptoKey Опция для настройки ключа хэширования тела.
func WithCryptoKey(key string) Option {
	return func(c *BaseClient) error {
		var (
			cryptoKey *rsa.PublicKey
			err       error
		)

		if key == "" {
			c.CryptoKey = nil
			return nil
		}

		cryptoKey, err = rsaKeys.ReadPublicKey(key)
		if err != nil {
			return fmt.Errorf("failed to read public key: %w", err)
		}
		c.CryptoKey = cryptoKey
		return nil
	}
}

// WithCompressType Опция для настройки типа сжатия.
func WithCompressType(compressType string) Option {
	return func(c *BaseClient) error {
		c.CompressType = compressType
		return nil
	}
}

func IsRetryableHTTPRequest(result any, _ error) bool {
	resp, ok := result.(*http.Response)
	if !ok {
		logrus.Errorf("failed to check type for %v", result)
		return false
	}
	if resp.StatusCode >= http.StatusInternalServerError {
		return true
	}
	return false
}

func IsRetryableGRPCRequest(_ any, err error) bool {
	if e, ok := status.FromError(err); ok {
		if e.Code() == codes.Internal || e.Code() == codes.Unavailable {
			return true
		}
	}
	return false
}
