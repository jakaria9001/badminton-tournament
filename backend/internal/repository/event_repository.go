package repository

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/jakaria9001/badminton-tournament/backend/internal/model"
)

type EventRepository struct {
	db *pgxpool.Pool
}

func NewEventRepository(db *pgxpool.Pool) *EventRepository {
	return &EventRepository{
		db: db,
	}
}

func (r *EventRepository) GetEventByID(
	ctx context.Context,
	eventID uuid.UUID,
) (*model.EventResponse, error) {

	query := `
		SELECT
			e.id,
			e.name,
			e.max_teams,
			tournament.status,
			COUNT(
				CASE
					WHEN team.status = 'CONFIRMED'
					THEN team.id
				END
			)::int AS registered_teams
		FROM events e
		JOIN tournaments tournament
			ON tournament.id = e.tournament_id
		LEFT JOIN teams team
			ON team.event_id = e.id
		WHERE e.id = $1
		GROUP BY
			e.id,
			e.name,
			e.max_teams,
			tournament.status
	`

	var event model.EventResponse

	err := r.db.QueryRow(
		ctx,
		query,
		eventID,
	).Scan(
		&event.ID,
		&event.Name,
		&event.MaxTeams,
		&event.Status,
		&event.RegisteredTeams,
	)

	if err == pgx.ErrNoRows {
		return nil, model.ErrEventNotFound
	}
	if err != nil {
		return nil, fmt.Errorf(
			"get event: %w",
			err,
		)
	}

	return &event, nil
}
