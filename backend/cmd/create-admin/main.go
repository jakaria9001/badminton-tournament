package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"

	"github.com/jakaria9001/badminton-tournament/backend/internal/database"
)

func main() {
	password := os.Getenv("ADMIN_PASSWORD")

	if password == "" {
		log.Fatal("ADMIN_PASSWORD is not set")
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
			role
		)
		VALUES ($1, $2, $3, $4, $5)`,
		userID,
		"Tournament Admin",
		"admin@badminton.local",
		string(passwordHash),
		"ADMIN",
	)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Admin created successfully")
}
