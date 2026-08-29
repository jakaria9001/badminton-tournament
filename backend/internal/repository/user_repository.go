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
	EventID      uuid.NullUUID
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
			role,
			event_id
		 FROM users
		 WHERE email = $1`,
		email,
	).Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.PasswordHash,
		&user.Role,
		&user.EventID,
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
        SELECT id, name, email, password_hash, role, event_id
        FROM users
        WHERE id = $1
    `, userID).Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.PasswordHash,
		&user.Role,
		&user.EventID,
	)

	if err != nil {
		return nil, fmt.Errorf("get user by ID: %w", err)
	}

	return &user, nil
}

func (r *UserRepository) CreateAdmin(
	ctx context.Context,
	name string,
	email string,
	passwordHash string,
	role string,
	eventID *uuid.UUID,
) error {
	_, err := r.db.Exec(ctx, `
		INSERT INTO users (id, name, email, password_hash, role, event_id)
		VALUES ($1, $2, $3, $4, $5, $6)
	`, uuid.New(), name, email, passwordHash, role, eventID)
	return err
}

func (r *UserRepository) ListAdmins(
	ctx context.Context,
) ([]User, error) {
	rows, err := r.db.Query(ctx, `
		SELECT id, name, email, password_hash, role, event_id
		FROM users
		WHERE role IN ('ADMIN', 'SUPER_ADMIN')
		ORDER BY name ASC
	`)
	if err != nil {
		return nil, fmt.Errorf("list admins: %w", err)
	}
	defer rows.Close()

	admins := make([]User, 0)
	for rows.Next() {
		var user User
		if err := rows.Scan(&user.ID, &user.Name, &user.Email, &user.PasswordHash, &user.Role, &user.EventID); err != nil {
			return nil, fmt.Errorf("scan admin: %w", err)
		}
		admins = append(admins, user)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate admins: %w", err)
	}
	return admins, nil
}
