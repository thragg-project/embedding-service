package heartbeat

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/fr33dman/go-template/pkg/probes"
)

func TestProbesServerRoutes(t *testing.T) {
	appProbes := probes.NewProbes()
	server := NewProbesServer("127.0.0.1", 0, appProbes)

	for _, tc := range []struct {
		name     string
		alive    bool
		ready    bool
		method   string
		path     string
		wantCode int
		wantBody string
	}{
		{"alive", true, false, http.MethodGet, "/livez", http.StatusOK, "1\n"},
		{"not ready", true, false, http.MethodGet, "/readyz", http.StatusServiceUnavailable, "0\n"},
		{"ready", true, true, http.MethodGet, "/readyz", http.StatusOK, "1\n"},
		{"not alive", false, true, http.MethodGet, "/livez", http.StatusServiceUnavailable, "0\n"},
		{"post alive", true, false, http.MethodPost, "/livez", http.StatusOK, "1\n"},
		{"post not ready", true, false, http.MethodPost, "/readyz", http.StatusServiceUnavailable, "0\n"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			appProbes.SetLiveness(tc.alive)
			appProbes.SetReadiness(tc.ready)
			req := httptest.NewRequestWithContext(
				context.Background(), tc.method, tc.path, nil,
			)
			rec := httptest.NewRecorder()
			server.Handler.ServeHTTP(rec, req)
			if rec.Code != tc.wantCode || rec.Body.String() != tc.wantBody {
				t.Fatalf(
					"got status %d and body %q, want %d and %q",
					rec.Code, rec.Body.String(), tc.wantCode, tc.wantBody,
				)
			}
		})
	}
}
