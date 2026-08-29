package repository

import (
	"context"
	"fmt"
	"regexp"
	"strings"
	"time"

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

func validateEventRequest(request model.EventAdminRequest) error {
	if strings.TrimSpace(request.Name) == "" || strings.TrimSpace(request.VenueName) == "" {
		return fmt.Errorf("event name and venue are required")
	}
	if request.MaxTeams != nil && *request.MaxTeams < 1 {
		return fmt.Errorf("max teams must be greater than 0")
	}
	startDate, err := parseDate(request.StartDate)
	if err != nil {
		return fmt.Errorf("invalid start date: %w", err)
	}
	endDate, err := parseDate(request.EndDate)
	if err != nil {
		return fmt.Errorf("invalid end date: %w", err)
	}
	if endDate.Before(startDate) {
		return fmt.Errorf("end date must be on or after start date")
	}
	if request.RegistrationDeadline != nil {
		deadline, err := time.Parse(time.RFC3339, *request.RegistrationDeadline)
		if err != nil {
			return fmt.Errorf("invalid registration deadline: %w", err)
		}
		if !deadline.Before(startDate) {
			return fmt.Errorf("registration deadline must be before start date")
		}
	}
	if request.Status != "" && request.Status != "DRAFT" && request.Status != "REGISTRATION_OPEN" && request.Status != "REGISTRATION_CLOSED" {
		return fmt.Errorf("invalid event status")
	}
	return nil
}

func parseDate(value string) (time.Time, error) {
	if value == "" {
		return time.Time{}, fmt.Errorf("date is required")
	}
	return time.Parse("2006-01-02", value)
}

func buildSlug(value string) string {
	sanitized := strings.TrimSpace(value)
	sanitized = strings.ToLower(sanitized)
	sanitized = strings.ReplaceAll(sanitized, "&", " and ")
	sanitized = regexp.MustCompile(`[^a-z0-9]+`).ReplaceAllString(sanitized, "-")
	sanitized = strings.Trim(sanitized, "-")
	if sanitized == "" {
		return "event-" + strings.ReplaceAll(uuid.NewString(), "-", "")[:12]
	}
	return sanitized + "-" + strings.ReplaceAll(uuid.NewString(), "-", "")[:8]
}

func (r *EventRepository) CreateEvent(ctx context.Context, request model.EventAdminRequest) (uuid.UUID, error) {
	if err := validateEventRequest(request); err != nil {
		return uuid.Nil, err
	}
	status := request.Status
	if status == "" {
		status = "DRAFT"
	}
	tournamentID := uuid.New()
	eventID := uuid.New()
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return uuid.Nil, fmt.Errorf("begin create event: %w", err)
	}
	defer tx.Rollback(ctx)
	_, err = tx.Exec(ctx, `
		INSERT INTO tournaments (id, name, slug, description, venue_name, venue_address, start_date, end_date, registration_deadline, status)
		VALUES ($1, $2, $3, NULLIF($4, ''), $5, NULLIF($6, ''), $7, $8, $9, $10)
	`, tournamentID, request.Name, buildSlug(request.Name), request.Description, request.VenueName, request.VenueAddress, request.StartDate, request.EndDate, request.RegistrationDeadline, status)
	if err != nil {
		return uuid.Nil, fmt.Errorf("create tournament: %w", err)
	}
	_, err = tx.Exec(ctx, `
		INSERT INTO events (id, tournament_id, name, event_type, max_teams)
		VALUES ($1, $2, $3, 'MENS_DOUBLES', $4)
	`, eventID, tournamentID, request.Name, request.MaxTeams)
	if err != nil {
		return uuid.Nil, fmt.Errorf("create event: %w", err)
	}
	if err := assignAdminToEvent(ctx, tx, eventID, request.AssignedAdminID); err != nil {
		return uuid.Nil, err
	}
	if err := tx.Commit(ctx); err != nil {
		return uuid.Nil, fmt.Errorf("commit event: %w", err)
	}
	return eventID, nil
}

func assignAdminToEvent(ctx context.Context, tx pgx.Tx, eventID uuid.UUID, assignedAdminID *string) error {
	if assignedAdminID == nil || strings.TrimSpace(*assignedAdminID) == "" {
		return nil
	}

	adminID, err := uuid.Parse(*assignedAdminID)
	if err != nil {
		return fmt.Errorf("invalid assigned admin ID")
	}
	var exists bool
	err = tx.QueryRow(ctx, `SELECT EXISTS(SELECT 1 FROM users WHERE id = $1 AND role = 'ADMIN')`, adminID).Scan(&exists)
	if err != nil {
		return fmt.Errorf("check assigned admin: %w", err)
	}
	if !exists {
		return fmt.Errorf("admin not found or not eligible for assignment")
	}

	if _, err := tx.Exec(ctx, `UPDATE users SET event_id = NULL WHERE event_id = $1 AND role = 'ADMIN'`, eventID); err != nil {
		return fmt.Errorf("clear event admin assignment: %w", err)
	}
	if _, err := tx.Exec(ctx, `UPDATE users SET event_id = $1 WHERE id = $2`, eventID, adminID); err != nil {
		return fmt.Errorf("assign admin to event: %w", err)
	}
	return nil
}

func (r *EventRepository) DeleteEvent(ctx context.Context, eventID uuid.UUID) error {
	result, err := r.db.Exec(ctx, `DELETE FROM tournaments WHERE id = (SELECT tournament_id FROM events WHERE id = $1)`, eventID)
	if err != nil {
		return fmt.Errorf("delete event: %w", err)
	}
	if result.RowsAffected() == 0 {
		return model.ErrEventNotFound
	}
	return nil
}

func (r *EventRepository) UpdateEvent(ctx context.Context, eventID uuid.UUID, request model.EventAdminRequest) error {
	if err := validateEventRequest(request); err != nil {
		return err
	}
	status := request.Status
	if status == "" {
		status = "DRAFT"
	}
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("begin update event: %w", err)
	}
	defer tx.Rollback(ctx)
	result, err := tx.Exec(ctx, `
		UPDATE tournaments
		SET name = $1, slug = $2, description = NULLIF($3, ''), venue_name = $4, venue_address = NULLIF($5, ''),
			start_date = $6, end_date = $7, registration_deadline = $8, status = $9, updated_at = NOW()
		WHERE id = (SELECT tournament_id FROM events WHERE id = $10)
	`, request.Name, buildSlug(request.Name), request.Description, request.VenueName, request.VenueAddress, request.StartDate, request.EndDate, request.RegistrationDeadline, status, eventID)
	if err != nil {
		return fmt.Errorf("update tournament: %w", err)
	}
	if result.RowsAffected() == 0 {
		return model.ErrEventNotFound
	}
	_, err = tx.Exec(ctx, `UPDATE events SET name = $1, max_teams = $2 WHERE id = $3`, request.Name, request.MaxTeams, eventID)
	if err != nil {
		return fmt.Errorf("update event: %w", err)
	}
	if err := assignAdminToEvent(ctx, tx, eventID, request.AssignedAdminID); err != nil {
		return err
	}
	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("commit event update: %w", err)
	}
	return nil
}

func (r *EventRepository) UpdateRegistrationStatus(
	ctx context.Context,
	eventID uuid.UUID,
	status string,
) error {
	result, err := r.db.Exec(ctx, `
		UPDATE tournaments
		SET status = $1, updated_at = NOW()
		WHERE id = (
			SELECT tournament_id
			FROM events
			WHERE id = $2
		)
	`, status, eventID)
	if err != nil {
		return fmt.Errorf("update registration status: %w", err)
	}
	if result.RowsAffected() == 0 {
		return model.ErrEventNotFound
	}
	return nil
}

func (r *EventRepository) ListEvents(ctx context.Context) ([]model.EventResponse, error) {
	return r.listEvents(ctx, false)
}

func (r *EventRepository) ListAdminEvents(ctx context.Context) ([]model.EventResponse, error) {
	return r.listEvents(ctx, true)
}

func (r *EventRepository) listEvents(
	ctx context.Context,
	includeDrafts bool,
) ([]model.EventResponse, error) {
	rows, err := r.db.Query(ctx, `
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
		WHERE ($1 OR tournament.status IN ('REGISTRATION_OPEN', 'REGISTRATION_CLOSED'))
		GROUP BY e.id, e.name, e.max_teams, tournament.status, e.created_at
		ORDER BY e.created_at DESC, e.id DESC
	`, includeDrafts)
	if err != nil {
		return nil, fmt.Errorf("list events: %w", err)
	}
	defer rows.Close()

	events := make([]model.EventResponse, 0)
	for rows.Next() {
		var event model.EventResponse
		if err := rows.Scan(
			&event.ID,
			&event.Name,
			&event.MaxTeams,
			&event.Status,
			&event.RegisteredTeams,
		); err != nil {
			return nil, fmt.Errorf("scan event: %w", err)
		}
		events = append(events, event)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate events: %w", err)
	}

	return events, nil
}

func (r *EventRepository) GetEventByID(
	ctx context.Context,
	eventID uuid.UUID,
) (*model.EventResponse, error) {
	return r.getEventByID(ctx, eventID, false)
}

func (r *EventRepository) getEventByID(
	ctx context.Context,
	eventID uuid.UUID,
	includeDrafts bool,
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
			AND ($2 OR tournament.status IN ('REGISTRATION_OPEN', 'REGISTRATION_CLOSED'))
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
		includeDrafts,
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
