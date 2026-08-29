package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRequireTrustedOrigin(t *testing.T) {
	middleware := RequireTrustedOrigin(map[string]struct{}{"http://localhost:5173": {}})
	handler := middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	}))

	for _, test := range []struct {
		name   string
		method string
		origin string
		want   int
	}{
		{name: "allows safe request without origin", method: http.MethodGet, want: http.StatusNoContent},
		{name: "allows trusted mutation", method: http.MethodPost, origin: "http://localhost:5173", want: http.StatusNoContent},
		{name: "rejects missing origin mutation", method: http.MethodPost, want: http.StatusForbidden},
		{name: "rejects untrusted origin", method: http.MethodPost, origin: "https://attacker.example", want: http.StatusForbidden},
	} {
		t.Run(test.name, func(t *testing.T) {
			recorder := httptest.NewRecorder()
			request := httptest.NewRequest(test.method, "/", nil)
			if test.origin != "" {
				request.Header.Set("Origin", test.origin)
			}
			handler.ServeHTTP(recorder, request)
			if recorder.Code != test.want {
				t.Fatalf("status = %d, want %d", recorder.Code, test.want)
			}
		})
	}
}