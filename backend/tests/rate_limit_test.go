package tests

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/jakaria9001/badminton-tournament/backend/internal/middleware"
)

func TestRateLimiterBlocksAfterLimit(t *testing.T) {
	limiter := middleware.NewRateLimiter(2, time.Minute)
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	})
	handler := limiter.Middleware(next)

	for requestNumber := 1; requestNumber <= 3; requestNumber++ {
		recording := httptest.NewRecorder()
		request := httptest.NewRequest(http.MethodPost, "/", nil)
		request.RemoteAddr = "192.0.2.10:1234"

		handler.ServeHTTP(recording, request)

		if requestNumber <= 2 && recording.Code != http.StatusNoContent {
			t.Fatalf("request %d: expected 204, got %d", requestNumber, recording.Code)
		}
		if requestNumber == 3 {
			if recording.Code != http.StatusTooManyRequests {
				t.Fatalf("expected 429, got %d", recording.Code)
			}
			if recording.Header().Get("Retry-After") == "" {
				t.Fatal("expected Retry-After header")
			}
		}
	}
}

func TestRateLimiterSeparatesClientIPs(t *testing.T) {
	limiter := middleware.NewRateLimiter(1, time.Minute)
	handler := limiter.Middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	}))

	for _, address := range []string{"192.0.2.11:1234", "192.0.2.12:1234"} {
		recording := httptest.NewRecorder()
		request := httptest.NewRequest(http.MethodPost, "/", nil)
		request.RemoteAddr = address
		handler.ServeHTTP(recording, request)
		if recording.Code != http.StatusNoContent {
			t.Fatalf("%s: expected 204, got %d", address, recording.Code)
		}
	}
}
