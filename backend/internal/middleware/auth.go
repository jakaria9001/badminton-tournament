package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/jakaria9001/badminton-tournament/backend/internal/model"
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

				if role != model.RoleAdmin && role != model.RoleSuperAdmin {
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

func RequireRole(roles ...string) func(http.Handler) http.Handler {
	allowed := make(map[string]struct{}, len(roles))
	for _, role := range roles {
		allowed[role] = struct{}{}
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			role, ok := r.Context().Value(RoleKey).(string)
			if !ok {
				http.Error(w, "authenticated role required", http.StatusForbidden)
				return
			}
			if _, ok := allowed[role]; !ok {
				http.Error(w, "insufficient permissions", http.StatusForbidden)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

func RequireEventAccess(db *pgxpool.Pool) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			role, _ := r.Context().Value(RoleKey).(string)
			if role == model.RoleSuperAdmin {
				next.ServeHTTP(w, r)
				return
			}

			if role != model.RoleAdmin {
				http.Error(w, "insufficient permissions", http.StatusForbidden)
				return
			}

			eventIDString := chi.URLParam(r, "eventID")
			if eventIDString == "" {
				http.Error(w, "event ID required", http.StatusBadRequest)
				return
			}

			eventID, err := uuid.Parse(eventIDString)
			if err != nil {
				http.Error(w, "invalid event ID", http.StatusBadRequest)
				return
			}

			userID, ok := r.Context().Value(UserIDKey).(uuid.UUID)
			if !ok {
				http.Error(w, "authenticated user required", http.StatusUnauthorized)
				return
			}

			var assignedEventID uuid.NullUUID
			err = db.QueryRow(r.Context(), `SELECT event_id FROM users WHERE id = $1`, userID).Scan(&assignedEventID)
			if err != nil {
				http.Error(w, "event access not mapped", http.StatusForbidden)
				return
			}

			if !assignedEventID.Valid || assignedEventID.UUID != eventID {
				http.Error(w, "event access denied", http.StatusForbidden)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

func RequireRegistrationAccess(db *pgxpool.Pool) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			role, _ := r.Context().Value(RoleKey).(string)
			if role == model.RoleSuperAdmin {
				next.ServeHTTP(w, r)
				return
			}

			if role != model.RoleAdmin {
				http.Error(w, "insufficient permissions", http.StatusForbidden)
				return
			}

			userID, ok := r.Context().Value(UserIDKey).(uuid.UUID)
			if !ok {
				http.Error(w, "authenticated user required", http.StatusUnauthorized)
				return
			}

			registrationIDString := chi.URLParam(r, "registrationID")
			if registrationIDString == "" {
				http.Error(w, "registration ID required", http.StatusBadRequest)
				return
			}

			registrationID, err := uuid.Parse(registrationIDString)
			if err != nil {
				http.Error(w, "invalid registration ID", http.StatusBadRequest)
				return
			}

			var mappedEventID uuid.NullUUID
			err = db.QueryRow(r.Context(), `
				SELECT t.event_id
				FROM registrations r
				JOIN teams t ON t.id = r.team_id
				WHERE r.id = $1
			`, registrationID).Scan(&mappedEventID)
			if err != nil {
				http.Error(w, "registration access denied", http.StatusForbidden)
				return
			}

			if !mappedEventID.Valid {
				http.Error(w, "registration access denied", http.StatusForbidden)
				return
			}

			var assignedEventID uuid.NullUUID
			err = db.QueryRow(r.Context(), `SELECT event_id FROM users WHERE id = $1`, userID).Scan(&assignedEventID)
			if err != nil || !assignedEventID.Valid || assignedEventID.UUID != mappedEventID.UUID {
				http.Error(w, "registration access denied", http.StatusForbidden)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

func RequireRoundAccess(db *pgxpool.Pool) func(http.Handler) http.Handler {
	return requireResourceEventAccess(db, "roundID", `SELECT event_id FROM tournament_rounds WHERE id = $1`)
}

func RequireMatchAccess(db *pgxpool.Pool) func(http.Handler) http.Handler {
	return requireResourceEventAccess(db, "matchID", `SELECT event_id FROM matches WHERE id = $1`)
}

func requireResourceEventAccess(db *pgxpool.Pool, parameter, query string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			role, _ := r.Context().Value(RoleKey).(string)
			if role == model.RoleSuperAdmin {
				next.ServeHTTP(w, r)
				return
			}
			if role != model.RoleAdmin {
				http.Error(w, "insufficient permissions", http.StatusForbidden)
				return
			}

			resourceID, err := uuid.Parse(chi.URLParam(r, parameter))
			if err != nil {
				http.Error(w, "invalid resource ID", http.StatusBadRequest)
				return
			}
			userID, ok := r.Context().Value(UserIDKey).(uuid.UUID)
			if !ok {
				http.Error(w, "authenticated user required", http.StatusUnauthorized)
				return
			}

			var resourceEventID, assignedEventID uuid.NullUUID
			err = db.QueryRow(r.Context(), query, resourceID).Scan(&resourceEventID)
			if err != nil || !resourceEventID.Valid {
				http.Error(w, "event access denied", http.StatusForbidden)
				return
			}
			err = db.QueryRow(r.Context(), `SELECT event_id FROM users WHERE id = $1`, userID).Scan(&assignedEventID)
			if err != nil || !assignedEventID.Valid || assignedEventID.UUID != resourceEventID.UUID {
				http.Error(w, "event access denied", http.StatusForbidden)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
