package httpserver

import (
	"bytes"
	"crypto/rand"
	"crypto/sha256"
	"io"
	"net/http"

	rsaKeys "github.com/NStegura/metrics/utils/rsa"
)

func (s *APIServer) decryptMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if s.cryptoKey != nil {
			bodyBytes, err := io.ReadAll(r.Body)
			if err != nil {
				http.Error(w, "failed to read body", http.StatusBadRequest)
				return
			}
			decryptedMessage, err := rsaKeys.DecryptOAEP(
				sha256.New(), rand.Reader, s.cryptoKey, bodyBytes, nil)
			if err != nil {
				http.Error(w, "failed to decrypt body", http.StatusBadRequest)
				s.logger.Error(err)
				return
			}
			r.Body = io.NopCloser(bytes.NewBuffer(decryptedMessage))
		}
		h.ServeHTTP(w, r)
	})
}
