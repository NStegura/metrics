package metricsapi

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"io"
	"net/http"
)

func (s *APIServer) hashValidation(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if s.config.RequestKey != "" {
			hm := hmac.New(sha256.New, []byte(s.config.RequestKey))
			b, err := io.ReadAll(r.Body)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			err = r.Body.Close()
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			r.Body = io.NopCloser(bytes.NewBuffer(b))

			hm.Write(b)
			calcHash := hm.Sum(nil)
			hashRequest, err := hex.DecodeString(r.Header.Get("HashSHA256"))
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			if !hmac.Equal(calcHash, hashRequest) {
				w.WriteHeader(http.StatusBadRequest)
				return
			}
		}
		h.ServeHTTP(w, r)
	})
}
