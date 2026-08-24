package repository

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/jakaria9001/badminton-tournament/backend/internal/model"
)

type RegistrationRepository struct {
	db *pgxpool.Pool
}

func NewRegistrationRepository(db *pgxpool.Pool) *RegistrationRepository {
	return &RegistrationRepository{
		db: db,
	}
}

func (r *RegistrationRepository) CreateRegistration(
	ctx context.Context,
	eventID uuid.UUID,
	req model.RegistrationRequest,
) (uuid.UUID, error) {

	tx, err := r.db.Begin(ctx)
	if err != nil {
		return uuid.Nil, fmt.Errorf("begin transaction: %w", err)
	}

	defer func() {
		_ = tx.Rollback(ctx)
	}()

	var registrationStatus string
	var registrationDeadline *time.Time
	var maxTeams *int

	err = tx.QueryRow(ctx, `
		SELECT tournament.status, tournament.registration_deadline, e.max_teams
		FROM events e
		JOIN tournaments tournament ON tournament.id = e.tournament_id
		WHERE e.id = $1
		FOR UPDATE OF e, tournament
	`, eventID).Scan(&registrationStatus, &registrationDeadline, &maxTeams)
	if err == pgx.ErrNoRows {
		return uuid.Nil, model.ErrEventNotFound
	}
	if err != nil {
		return uuid.Nil, fmt.Errorf("check event registration: %w", err)
	}
	if registrationStatus != "REGISTRATION_OPEN" {
		return uuid.Nil, model.ErrRegistrationClosed
	}
	if registrationDeadline != nil && time.Now().After(*registrationDeadline) {
		return uuid.Nil, model.ErrRegistrationDeadlinePassed
	}
	if maxTeams != nil {
		var registeredTeams int
		err = tx.QueryRow(ctx, `
			SELECT COUNT(*)
			FROM teams
			WHERE event_id = $1 AND status NOT IN ('REJECTED', 'WITHDRAWN')
		`, eventID).Scan(&registeredTeams)
		if err != nil {
			return uuid.Nil, fmt.Errorf("check event capacity: %w", err)
		}
		if registeredTeams >= *maxTeams {
			return uuid.Nil, model.ErrEventFull
		}
	}

	// Create Player 1
	var player1ID uuid.UUID

	err = tx.QueryRow(
		ctx,
		`SELECT id
		FROM players
		WHERE phone = $1`,
		req.Player1.Phone,
	).Scan(&player1ID)

	if err == pgx.ErrNoRows {
		player1ID = uuid.New()

		_, err = tx.Exec(
			ctx,
			`INSERT INTO players (
				id,
				name,
				phone
			)
			VALUES ($1, $2, $3)`,
			player1ID,
			req.Player1.Name,
			req.Player1.Phone,
		)

		if err != nil {
			return uuid.Nil, fmt.Errorf(
				"create player 1: %w",
				err,
			)
		}
	} else if err != nil {
		return uuid.Nil, fmt.Errorf(
			"find player 1: %w",
			err,
		)
	}

	// Create Player 2
	var player2ID uuid.UUID

	player2Phone := strings.TrimSpace(req.Player2.Phone)
	if player2Phone != "" {
		err = tx.QueryRow(
			ctx,
			`SELECT id
			FROM players
			WHERE phone = $1`,
			player2Phone,
		).Scan(&player2ID)
	} else {
		err = pgx.ErrNoRows
	}

	if err == pgx.ErrNoRows {
		player2ID = uuid.New()

		_, err = tx.Exec(
			ctx,
			`INSERT INTO players (
				id,
				name,
				phone
			)
			VALUES ($1, $2, NULLIF($3, ''))`,
			player2ID,
			req.Player2.Name,
			player2Phone,
		)

		if err != nil {
			return uuid.Nil, fmt.Errorf(
				"create player 2: %w",
				err,
			)
		}
	} else if err != nil {
		return uuid.Nil, fmt.Errorf(
			"find player 2: %w",
			err,
		)
	}

	// check for player 1's existing registration in the same event
	var existingTeamID uuid.UUID

	err = tx.QueryRow(
		ctx,
		`SELECT id
		FROM teams
		WHERE event_id = $1
		AND (
				player1_id = $2
				OR player2_id = $2
		)`,
		eventID,
		player1ID,
	).Scan(&existingTeamID)

	if err == nil {
		return uuid.Nil, fmt.Errorf(
			"%w: player 1", model.ErrParticipantAlreadyRegistered,
		)
	}

	if err != pgx.ErrNoRows {
		return uuid.Nil, fmt.Errorf(
			"check player 1 registration: %w",
			err,
		)
	}

	// check for player 2's existing registration in the same event
	err = tx.QueryRow(
		ctx,
		`SELECT id
		FROM teams
		WHERE event_id = $1
		AND (
				player1_id = $2
				OR player2_id = $2
		)`,
		eventID,
		player2ID,
	).Scan(&existingTeamID)

	if err == nil {
		return uuid.Nil, fmt.Errorf(
			"%w: player 2", model.ErrParticipantAlreadyRegistered,
		)
	}

	if err != pgx.ErrNoRows {
		return uuid.Nil, fmt.Errorf(
			"check player 2 registration: %w",
			err,
		)
	}

	// Create team
	teamID := uuid.New()

	pairKey := generatePairKey(
		player1ID,
		player2ID,
	)

	_, err = tx.Exec(
		ctx,
		`INSERT INTO teams (
			id,
			event_id,
			player1_id,
			player2_id,
			team_name,
			player_pair_key
		)
		VALUES ($1, $2, $3, $4, NULLIF($5, ''), $6)`,
		teamID,
		eventID,
		player1ID,
		player2ID,
		req.TeamName,
		pairKey,
	)

	if err != nil {
		return uuid.Nil, fmt.Errorf("create team: %w", err)
	}

	// Create registration
	registrationID := uuid.New()

	_, err = tx.Exec(
		ctx,
		`INSERT INTO registrations (
			id,
			team_id
		)
		VALUES ($1, $2)`,
		registrationID,
		teamID,
	)

	if err != nil {
		return uuid.Nil, fmt.Errorf("create registration: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return uuid.Nil, fmt.Errorf("commit transaction: %w", err)
	}

	return registrationID, nil
}

func (r *RegistrationRepository) GetTeamsByEvent(
	ctx context.Context,
	eventID uuid.UUID,
) ([]model.TeamResponse, error) {

	query := `
		SELECT
			t.id,
			COALESCE(NULLIF(t.team_name, ''), p1.name || ' / ' || p2.name) AS team_name,
			p1.id,
			p1.name,
			p2.id,
			p2.name,
			t.status,
			t.created_at
		FROM teams t
		JOIN players p1
			ON p1.id = t.player1_id
		JOIN players p2
			ON p2.id = t.player2_id
		WHERE t.event_id = $1
			AND t.status = 'CONFIRMED'
		ORDER BY t.created_at ASC
	`

	rows, err := r.db.Query(ctx, query, eventID)
	if err != nil {
		return nil, fmt.Errorf("get teams: %w", err)
	}
	defer rows.Close()

	teams := make([]model.TeamResponse, 0)

	for rows.Next() {
		var team model.TeamResponse
		var (
			teamID    uuid.UUID
			player1ID uuid.UUID
			player2ID uuid.UUID
			createdAt time.Time
		)

		err := rows.Scan(
			&teamID,
			&team.TeamName,
			&player1ID,
			&team.Player1.Name,
			&player2ID,
			&team.Player2.Name,
			&team.Status,
			&createdAt,
		)
		if err != nil {
			return nil, fmt.Errorf("scan team: %w", err)
		}

		team.ID = teamID.String()
		team.Player1.ID = player1ID.String()
		team.Player2.ID = player2ID.String()
		team.CreatedAt = createdAt.Format(time.RFC3339)

		teams = append(teams, team)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate teams: %w", err)
	}

	return teams, nil
}

func (r *RegistrationRepository) GetRegistrationsByEvent(
	ctx context.Context,
	eventID uuid.UUID,
) ([]model.AdminRegistration, error) {

	query := `
		SELECT
			r.id,
			t.id,
			COALESCE(
				NULLIF(t.team_name, ''),
				p1.name || ' / ' || p2.name
			),
			p1.name,
			p1.phone,
			p2.name,
			COALESCE(p2.phone, ''),
			r.status,
			r.registered_at
		FROM registrations r
		JOIN teams t
			ON t.id = r.team_id
		JOIN players p1
			ON p1.id = t.player1_id
		JOIN players p2
			ON p2.id = t.player2_id
		WHERE t.event_id = $1
		ORDER BY r.registered_at ASC
	`

	rows, err := r.db.Query(
		ctx,
		query,
		eventID,
	)

	if err != nil {
		return nil, fmt.Errorf(
			"get registrations: %w",
			err,
		)
	}

	defer rows.Close()

	registrations := make([]model.AdminRegistration, 0)

	for rows.Next() {
		var registration model.AdminRegistration
		var (
			registrationID uuid.UUID
			teamID         uuid.UUID
			registeredAt   time.Time
		)

		err := rows.Scan(
			&registrationID,
			&teamID,
			&registration.TeamName,
			&registration.Player1Name,
			&registration.Player1Phone,
			&registration.Player2Name,
			&registration.Player2Phone,
			&registration.Status,
			&registeredAt,
		)

		if err != nil {
			return nil, fmt.Errorf(
				"scan registration: %w",
				err,
			)
		}

		registration.ID = registrationID.String()
		registration.TeamID = teamID.String()
		registration.RegisteredAt = registeredAt.Format(time.RFC3339)

		registrations =
			append(
				registrations,
				registration,
			)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf(
			"iterate registrations: %w",
			err,
		)
	}

	return registrations, nil
}

func generatePairKey(
	player1ID uuid.UUID,
	player2ID uuid.UUID,
) string {
	if player1ID.String() < player2ID.String() {
		return player1ID.String() + ":" + player2ID.String()
	}

	return player2ID.String() + ":" + player1ID.String()
}

func (r *RegistrationRepository) UpdateStatus(
	ctx context.Context,
	registrationID uuid.UUID,
	status string,
) error {

	query := `
		WITH updated_registration AS (
			UPDATE registrations
			SET
				status = $1::varchar,
				confirmed_at = CASE
					WHEN $1::varchar = 'CONFIRMED' THEN NOW()
					ELSE confirmed_at
				END
			WHERE id = $2
			RETURNING team_id
		)
		UPDATE teams t
		SET
			status = $1::varchar,
			updated_at = NOW()
		FROM updated_registration r
		WHERE t.id = r.team_id
	`

	result, err := r.db.Exec(
		ctx,
		query,
		status,
		registrationID,
	)

	if err != nil {
		return fmt.Errorf(
			"update registration status: %w",
			err,
		)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf(
			"registration not found",
		)
	}

	return nil
}
