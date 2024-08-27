package metricsapi

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"fmt"
	"io"
	"net/http"

	"github.com/NStegura/metrics/utils/pem"
)

type DecryptingReader struct {
	io.Reader
	privateKey *rsa.PrivateKey
}

func (dr *DecryptingReader) Read(p []byte) (n int, err error) {
	n, err = dr.Reader.Read(p)
	if err != nil {
		return n, fmt.Errorf("failed to read message, err %w", err)
	}

	decryptedMessage, err := rsa.DecryptOAEP(sha256.New(), rand.Reader, dr.privateKey, p[:n], nil)
	if err != nil {
		return 0, fmt.Errorf("failed to decrypt message, %w", err)
	}
	copy(p, decryptedMessage)
	return len(decryptedMessage), nil
}

func (s *APIServer) decryptMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if s.config.PrivateCryptoKey != "" {
			privateKey, err := pem.ReadPrivateKey(s.config.PrivateCryptoKey)
			if err != nil {
				http.Error(w, "Ошибка загрузки приватного ключа", http.StatusInternalServerError)
				return
			}
			// Создание кастомного Reader
			decryptingReader := &DecryptingReader{
				Reader:     r.Body,
				privateKey: privateKey,
			}

			// Заменяем r.Body на наш кастомный Reader
			r.Body = io.NopCloser(decryptingReader)
		}
		h.ServeHTTP(w, r)
	})
}
