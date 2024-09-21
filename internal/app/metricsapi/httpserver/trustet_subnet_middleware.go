package httpserver

import (
	"net"
	"net/http"
)

func (s *APIServer) trustedSubnetMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if s.cfg.TrustedSubnet == "" {
			h.ServeHTTP(w, r)
			return
		}

		realIP := r.Header.Get("X-Real-IP")
		if realIP == "" {
			http.Error(w, "X-Real-IP header is missing", http.StatusBadRequest)
			return
		}

		ip := net.ParseIP(realIP)
		if ip == nil || !s.trustedSubnet.Contains(ip) {
			w.WriteHeader(http.StatusForbidden)
			return
		}

		h.ServeHTTP(w, r)
	})
}
