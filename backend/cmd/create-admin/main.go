package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"

	"github.com/jakaria9001/badminton-tournament/backend/internal/database"
)

func main() {
	password := os.Getenv("ADMIN_PASSWORD")

	if password == "" {
		log.Fatal("ADMIN_PASSWORD is not set")
	}

	email := envOrDefault("ADMIN_EMAIL", "admin@badminton.local")
	name := envOrDefault("ADMIN_NAME", "Tournament Admin")
	role := envOrDefault("ADMIN_ROLE", "ADMIN")
	if role != "ADMIN" && role != "SUPER_ADMIN" {
		log.Fatal("ADMIN_ROLE must be ADMIN or SUPER_ADMIN")
	}

	var eventID *uuid.UUID
	if role == "ADMIN" {
		configuredEventID := strings.TrimSpace(os.Getenv("ADMIN_EVENT_ID"))
		if configuredEventID == "" {
			log.Fatal("ADMIN_EVENT_ID is required for ADMIN role")
		}
		parsedEventID, err := uuid.Parse(configuredEventID)
		if err != nil {
			log.Fatal("ADMIN_EVENT_ID must be a valid UUID")
		}
		eventID = &parsedEventID
	}

	passwordHash, err := bcrypt.GenerateFromPassword(
		[]byte(password),
		bcrypt.DefaultCost,
	)

	if err != nil {
		log.Fatal(err)
	}

	ctx := context.Background()

	databaseURL := os.Getenv("DATABASE_URL")

	if databaseURL == "" {
		log.Fatal("DATABASE_URL is not set")
	}

	db, err := database.NewPostgresPool(
		ctx,
		databaseURL,
	)

	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()

	userID := uuid.New()

	_, err = db.Exec(
		ctx,
		`INSERT INTO users (
			id,
			name,
			email,
			password_hash,
			role,
			event_id
		)
		VALUES ($1, $2, $3, $4, $5, $6)`,
		userID,
		name,
		email,
		string(passwordHash),
		role,
		eventID,
	)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%s created successfully: %s\n", role, email)
}

func envOrDefault(name string, fallback string) string {
	value := os.Getenv(name)
	if value == "" {
		return fallback
	}
	return value
}
