package main

import (
	"context"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"os"

	"github.com/jakaria9001/badminton-tournament/backend/internal/database"
	"github.com/jakaria9001/badminton-tournament/backend/internal/draw"
	"github.com/jakaria9001/badminton-tournament/backend/internal/handler"
	"github.com/jakaria9001/badminton-tournament/backend/internal/repository"
	"github.com/jakaria9001/badminton-tournament/backend/internal/router"
	"github.com/jakaria9001/badminton-tournament/backend/internal/service"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: .env file not found")
	}

	ctx := context.Background()

	databaseURL := os.Getenv("DATABASE_URL")

	if databaseURL == "" {
		log.Fatal("DATABASE_URL is not set")
	}

	// AUTH
	jwtSecret := os.Getenv("JWT_SECRET")

	if jwtSecret == "" {
		log.Fatal("JWT_SECRET is not set")
	}

	// Database
	db, err := database.NewPostgresPool(
		ctx,
		databaseURL,
	)

	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()

	log.Println("Connected to PostgreSQL")

	// generator
	generator :=
		draw.NewGenerator()

	// Repository
	registrationRepository :=
		repository.NewRegistrationRepository(db)

	eventRepository :=
		repository.NewEventRepository(db)

	userRepository :=
		repository.NewUserRepository(db)

	matchRepository :=
		repository.NewMatchRepository(db)

	roundRepository :=
		repository.NewRoundRepository(db)

	advancementRepository :=
		repository.NewAdvancementRepository(db)

	// Service
	registrationService :=
		service.NewRegistrationService(
			registrationRepository,
			eventRepository,
		)

	eventService :=
		service.NewEventService(
			eventRepository,
		)

	authService :=
		service.NewAuthService(
			userRepository,
			jwtSecret,
		)

	matchService :=
		service.NewMatchService(
			matchRepository,
			roundRepository,
		)

	roundService :=
		service.NewRoundService(
			roundRepository,
			matchRepository,
		)

	resultService :=
		service.NewResultService(
			db,
			matchRepository,
		)

	drawService :=
		service.NewDrawService(
			db,
			roundRepository,
			matchRepository,
			advancementRepository,
			generator,
		)

	// Handler
	registrationHandler :=
		handler.NewRegistrationHandler(
			registrationService,
		)

	eventHandler :=
		handler.NewEventHandler(
			eventService,
		)

	authHandler :=
		handler.NewAuthHandler(
			authService,
		)

	matchHandler :=
		handler.NewMatchHandler(
			matchService,
		)

	roundHandler :=
		handler.NewRoundHandler(
			roundService,
			drawService,
		)

	resultHandler :=
		handler.NewResultHandler(
			resultService,
		)

	// Router
	r := router.NewRouter(
		registrationHandler,
		eventHandler,
		authHandler,
		matchHandler,
		roundHandler,
		resultHandler,
		jwtSecret,
	)

	port := os.Getenv("PORT")

	if port == "" {
		port = "8080"
	}

	log.Printf(
		"Server running on http://localhost:%s",
		port,
	)

	if err := http.ListenAndServe(
		":"+port,
		r,
	); err != nil {
		log.Fatal(err)
	}
}
