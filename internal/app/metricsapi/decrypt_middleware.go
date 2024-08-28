package metricsapi

import (
	"bytes"
	"crypto/rand"
	"crypto/sha256"
	"fmt"
	"io"
	"net/http"

	rsaKeys "github.com/NStegura/metrics/utils/rsa"
)

func (s *APIServer) decryptMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if s.config.PrivateCryptoKey != nil {
			bodyBytes, err := io.ReadAll(r.Body)
			if err != nil {
				http.Error(w, "failed to read body", http.StatusBadRequest)
				return
			}
			decryptedMessage, err := rsaKeys.DecryptOAEP(
				sha256.New(), rand.Reader, s.config.PrivateCryptoKey, bodyBytes, nil)
			if err != nil {
				http.Error(w, "failed to decrypt body", http.StatusBadRequest)
				fmt.Println(err)
				return
			}
			r.Body = io.NopCloser(bytes.NewBuffer(decryptedMessage))
		}
		h.ServeHTTP(w, r)
	})
}
