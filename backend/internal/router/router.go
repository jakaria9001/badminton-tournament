package router

import (
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/jakaria9001/badminton-tournament/backend/internal/handler"
	"github.com/jakaria9001/badminton-tournament/backend/internal/middleware"
)

func NewRouter(
	db *pgxpool.Pool,
	registrationHandler *handler.RegistrationHandler,
	eventHandler *handler.EventHandler,
	authHandler *handler.AuthHandler,
	matchHandler *handler.MatchHandler,
	roundHandler *handler.RoundHandler,
	resultHandler *handler.ResultHandler,
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

	r.Route("/api/v1", func(r chi.Router) {

		r.Get("/events", eventHandler.List)

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

			r.Get(
				"/matches",
				matchHandler.GetByEvent,
			)

			r.Get(
				"/rounds",
				roundHandler.GetByEvent,
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

		r.Route("/superadmin", func(r chi.Router) {
			r.Use(middleware.RequireRole("SUPER_ADMIN"))
			r.Get("/events", eventHandler.ListAdmin)
			r.Post("/events", eventHandler.Create)
			r.Put("/events/{eventID}", eventHandler.Update)
			r.Delete("/events/{eventID}", eventHandler.Delete)
			r.Get("/admins", authHandler.ListAdmins)
			r.Post("/admins", authHandler.CreateAdmin)
		})

		r.Route("/events/{eventID}", func(r chi.Router) {
			r.Use(middleware.RequireEventAccess(db))
			r.Post("/rounds", roundHandler.Create)
			r.Put("/registration-status", eventHandler.UpdateRegistrationStatus)
			r.Get("/rounds", roundHandler.GetByEvent)
			r.Get("/rounds/{roundID}/available-teams", roundHandler.GetAvailableTeams)
			r.Get("/registrations", registrationHandler.GetRegistrations)
		})

		r.Route("/registrations/{registrationID}", func(r chi.Router) {
			r.Use(middleware.RequireRegistrationAccess(db))
			r.Put("/status", registrationHandler.UpdateStatus)
			r.Put("/withdraw", registrationHandler.WithdrawRegistration)
		})

		r.With(middleware.RequireRoundAccess(db)).Post(
			"/rounds/{roundID}/matches",
			matchHandler.Create,
		)

		r.With(middleware.RequireRoundAccess(db)).Post(
			"/rounds/{roundID}/generate",
			roundHandler.Generate,
		)

		r.With(middleware.RequireRoundAccess(db)).Post(
			"/rounds/{roundID}/lock",
			roundHandler.Lock,
		)

		r.With(middleware.RequireMatchAccess(db)).Post(
			"/matches/{matchID}/result",
			resultHandler.Submit,
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
