package router

import (
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"

	"github.com/jakaria9001/badminton-tournament/backend/internal/handler"
	"github.com/jakaria9001/badminton-tournament/backend/internal/middleware"
)

func NewRouter(
	registrationHandler *handler.RegistrationHandler,
	eventHandler *handler.EventHandler,
	authHandler *handler.AuthHandler,
	jwtSecret string,
) http.Handler {
	allowedOrigins := strings.Split(os.Getenv("CORS_ALLOWED_ORIGINS"), ",")
	if allowedOrigins[0] == "" {
		allowedOrigins = []string{"http://localhost:5173"}
	}

	r := chi.NewRouter()
	loginLimiter := middleware.NewRateLimiter(
		envInt("LOGIN_RATE_LIMIT_MAX_REQUESTS", 5),
		envDurationSeconds("LOGIN_RATE_LIMIT_WINDOW_SECONDS", 60),
	)
	registrationLimiter := middleware.NewRateLimiter(
		envInt("REGISTRATION_RATE_LIMIT_MAX_REQUESTS", 10),
		envDurationSeconds("REGISTRATION_RATE_LIMIT_WINDOW_SECONDS", 60),
	)

	r.Use(cors.Handler(cors.Options{
		AllowedOrigins: allowedOrigins,
		AllowedMethods: []string{
			"GET",
			"POST",
			"PUT",
			"DELETE",
			"OPTIONS",
		},
		AllowedHeaders: []string{
			"Accept",
			"Authorization",
			"Content-Type",
		},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	r.Get("/health", func(
		w http.ResponseWriter,
		r *http.Request,
	) {
		w.Header().Set(
			"Content-Type",
			"application/json",
		)

		w.WriteHeader(http.StatusOK)

		w.Write([]byte(`{"status":"ok"}`))
	})

	r.Route("/api/v1", func(r chi.Router) {

		r.Route("/events/{eventID}", func(r chi.Router) {

			r.Get(
				"/",
				eventHandler.GetByID,
			)

			r.With(registrationLimiter.Middleware).Post(
				"/registrations",
				registrationHandler.Create,
			)

			r.Get(
				"/teams",
				registrationHandler.GetTeams,
			)
		})

		r.With(loginLimiter.Middleware).Post(
			"/auth/login",
			authHandler.Login,
		)
	})

	r.Route("/api/v1/admin", func(r chi.Router) {

		r.Use(
			middleware.RequireAuth(jwtSecret),
		)

		r.Get(
			"/events/{eventID}/registrations",
			registrationHandler.GetRegistrations,
		)

		r.Put(
			"/registrations/{registrationID}/status",
			registrationHandler.UpdateStatus,
		)

		r.Put(
			"/registrations/{registrationID}/withdraw",
			registrationHandler.WithdrawRegistration,
		)

		r.Get("/me", authHandler.Me)

	})

	return r
}

func envInt(name string, fallback int) int {
	value, err := strconv.Atoi(os.Getenv(name))
	if err != nil || value < 1 {
		return fallback
	}
	return value
}

func envDurationSeconds(name string, fallback int) time.Duration {
	return time.Duration(envInt(name, fallback)) * time.Second
}
