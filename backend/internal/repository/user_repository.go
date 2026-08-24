package repository

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type User struct {
	ID           uuid.UUID
	Name         string
	Email        string
	PasswordHash string
	Role         string
}

type UserRepository struct {
	db *pgxpool.Pool
}

func NewUserRepository(
	db *pgxpool.Pool,
) *UserRepository {
	return &UserRepository{
		db: db,
	}
}

func (r *UserRepository) GetByEmail(
	ctx context.Context,
	email string,
) (*User, error) {

	var user User

	err := r.db.QueryRow(
		ctx,
		`SELECT
			id,
			name,
			email,
			password_hash,
			role
		 FROM users
		 WHERE email = $1`,
		email,
	).Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.PasswordHash,
		&user.Role,
	)

	if err != nil {
		return nil, fmt.Errorf(
			"get user by email: %w",
			err,
		)
	}

	return &user, nil
}

func (r *UserRepository) GetByID(
	ctx context.Context,
	userID uuid.UUID,
) (*User, error) {
	var user User

	err := r.db.QueryRow(ctx, `
        SELECT id, name, email, password_hash, role
        FROM users
        WHERE id = $1
    `, userID).Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.PasswordHash,
		&user.Role,
	)

	if err != nil {
		return nil, fmt.Errorf("get user by ID: %w", err)
	}

	return &user, nil
}
