package middleware

import (
	"net/http"
)

func RequireTrustedOrigin(allowedOrigins map[string]struct{}) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Method == http.MethodGet || r.Method == http.MethodHead || r.Method == http.MethodOptions {
				next.ServeHTTP(w, r)
				return
			}

			origin := r.Header.Get("Origin")
			if _, allowed := allowedOrigins[origin]; !allowed {
				http.Error(w, "untrusted request origin", http.StatusForbidden)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}