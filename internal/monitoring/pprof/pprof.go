package pprof

import (
	"fmt"
	"net/http"
	"net/http/pprof"
)

// Start запускает профайлер на отдельном порту.
func Start() error {
	r := http.NewServeMux()

	r.HandleFunc("/debug/pprof/", pprof.Index)
	r.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
	r.HandleFunc("/debug/pprof/profile", pprof.Profile)
	r.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
	r.HandleFunc("/debug/pprof/trace", pprof.Trace)

	err := http.ListenAndServe(":8081", r)
	if err != nil {
		return fmt.Errorf("failed to start pprof server: %w", err)
	}
	return nil
}
