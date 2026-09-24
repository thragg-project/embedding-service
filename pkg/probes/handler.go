package probes

import (
	"net/http"
	"sync/atomic"
)

func (p *Probes) Handler() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("/livez", probeHandler(&p.isAlive))
	mux.HandleFunc("/readyz", probeHandler(&p.isReady))

	return mux
}

func probeHandler(state *atomic.Bool) http.HandlerFunc {
	return func(w http.ResponseWriter, _ *http.Request) {
		if state.Load() {
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte("1\n"))
			return
		}

		w.WriteHeader(http.StatusServiceUnavailable)
		_, _ = w.Write([]byte("0\n"))
	}
}
