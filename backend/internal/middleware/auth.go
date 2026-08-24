package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type contextKey string

const (
	UserIDKey contextKey = "userID"
	RoleKey   contextKey = "role"
)

func RequireAuth(
	jwtSecret string,
) func(http.Handler) http.Handler {

	return func(next http.Handler) http.Handler {

		return http.HandlerFunc(
			func(
				w http.ResponseWriter,
				r *http.Request,
			) {

				authHeader :=
					r.Header.Get("Authorization")

				if authHeader == "" {
					http.Error(
						w,
						"missing authorization header",
						http.StatusUnauthorized,
					)

					return
				}

				parts :=
					strings.SplitN(
						authHeader,
						" ",
						2,
					)

				if len(parts) != 2 ||
					parts[0] != "Bearer" {

					http.Error(
						w,
						"invalid authorization header",
						http.StatusUnauthorized,
					)

					return
				}

				tokenString := parts[1]

				token, err :=
					jwt.Parse(
						tokenString,
						func(token *jwt.Token) (
							interface{},
							error,
						) {

							if token.Method !=
								jwt.SigningMethodHS256 {

								return nil,
									jwt.ErrSignatureInvalid
							}

							return []byte(jwtSecret),
								nil
						},
					)

				if err != nil ||
					!token.Valid {

					http.Error(
						w,
						"invalid or expired token",
						http.StatusUnauthorized,
					)

					return
				}

				claims, ok :=
					token.Claims.(jwt.MapClaims)

				if !ok {
					http.Error(
						w,
						"invalid token claims",
						http.StatusUnauthorized,
					)

					return
				}

				role, ok :=
					claims["role"].(string)

				if !ok {
					http.Error(
						w,
						"invalid token role",
						http.StatusUnauthorized,
					)

					return
				}

				if role != "ADMIN" {
					http.Error(
						w,
						"admin access required",
						http.StatusForbidden,
					)

					return
				}

				userIDString, ok := claims["user_id"].(string)
				if !ok {
					http.Error(w, "invalid token user", http.StatusUnauthorized)
					return
				}

				userID, err := uuid.Parse(userIDString)
				if err != nil {
					http.Error(w, "invalid token user", http.StatusUnauthorized)
					return
				}

				ctx := context.WithValue(
					r.Context(),
					UserIDKey,
					userID,
				)

				ctx = context.WithValue(
					ctx,
					RoleKey,
					role,
				)

				next.ServeHTTP(
					w,
					r.WithContext(ctx),
				)
			},
		)
	}
}
