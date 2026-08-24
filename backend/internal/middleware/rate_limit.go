package middleware

import (
	"net"
	"net/http"
	"strconv"
	"sync"
	"time"
)

type rateLimitEntry struct {
	windowStart time.Time
	count       int
}

type RateLimiter struct {
	limit  int
	window time.Duration
	mu     sync.Mutex
	items  map[string]rateLimitEntry
}

func NewRateLimiter(limit int, window time.Duration) *RateLimiter {
	if limit < 1 {
		limit = 1
	}
	if window <= 0 {
		window = time.Minute
	}

	return &RateLimiter{
		limit:  limit,
		window: window,
		items:  make(map[string]rateLimitEntry),
	}
}

func (l *RateLimiter) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		key := clientIP(r)
		now := time.Now()

		l.mu.Lock()
		for itemKey, item := range l.items {
			if now.Sub(item.windowStart) >= l.window {
				delete(l.items, itemKey)
			}
		}

		item := l.items[key]
		if item.windowStart.IsZero() || now.Sub(item.windowStart) >= l.window {
			item = rateLimitEntry{windowStart: now}
		}

		item.count++
		l.items[key] = item
		allowed := item.count <= l.limit
		remaining := l.window - now.Sub(item.windowStart)
		l.mu.Unlock()

		if !allowed {
			retryAfter := int(remaining.Seconds())
			if retryAfter < 1 {
				retryAfter = 1
			}
			w.Header().Set("Retry-After", strconv.Itoa(retryAfter))
			http.Error(w, "rate limit exceeded", http.StatusTooManyRequests)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func clientIP(r *http.Request) string {
	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err == nil {
		return host
	}
	if r.RemoteAddr == "" {
		return "unknown"
	}
	return r.RemoteAddr
}
