package repository

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/jakaria9001/badminton-tournament/backend/internal/model"
)

type RoundRepository struct {
	db *pgxpool.Pool
}

func NewRoundRepository(
	db *pgxpool.Pool,
) *RoundRepository {
	return &RoundRepository{
		db: db,
	}
}

func (r *RoundRepository) Create(
	ctx context.Context,
	eventID uuid.UUID,
	roundNumber int,
	roundName string,
	pairingMethod string,
) (*model.TournamentRound, error) {

	id := uuid.New()

	var round model.TournamentRound

	err := r.db.QueryRow(
		ctx,
		`
		INSERT INTO tournament_rounds (
			id,
			event_id,
			round_number,
			round_name,
			pairing_method,
			status
		)
		VALUES ($1, $2, $3, $4, $5, 'OPEN')
		RETURNING
			id,
			event_id,
			round_number,
			round_name,
			pairing_method,
			status
		`,
		id,
		eventID,
		roundNumber,
		roundName,
		pairingMethod,
	).Scan(
		&round.ID,
		&round.EventID,
		&round.RoundNumber,
		&round.RoundName,
		&round.PairingMethod,
		&round.Status,
	)

	if err != nil {
		return nil, fmt.Errorf(
			"create round: %w",
			err,
		)
	}

	return &round, nil
}

func (r *RoundRepository) GetByEvent(
	ctx context.Context,
	eventID uuid.UUID,
) ([]model.TournamentRound, error) {

	rows, err := r.db.Query(
		ctx,
		`
		SELECT
			id,
			event_id,
			round_number,
			round_name,
			pairing_method,
			status,
			locked_at,
			completed_at
		FROM tournament_rounds
		WHERE event_id = $1
		ORDER BY round_number
		`,
		eventID,
	)

	if err != nil {
		return nil, fmt.Errorf(
			"get rounds: %w",
			err,
		)
	}

	defer rows.Close()

	rounds := make([]model.TournamentRound, 0)

	for rows.Next() {

		var round model.TournamentRound

		err := rows.Scan(
			&round.ID,
			&round.EventID,
			&round.RoundNumber,
			&round.RoundName,
			&round.PairingMethod,
			&round.Status,
			&round.LockedAt,
			&round.CompletedAt,
		)

		if err != nil {
			return nil, fmt.Errorf(
				"scan round: %w",
				err,
			)
		}

		rounds = append(rounds, round)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf(
			"iterate rounds: %w",
			err,
		)
	}

	return rounds, nil
}

func (r *RoundRepository) Lock(
	ctx context.Context,
	roundID uuid.UUID,
) error {

	result, err := r.db.Exec(
		ctx,
		`
		UPDATE tournament_rounds
		SET
			status = 'LOCKED',
			locked_at = NOW()
		WHERE id = $1
		  AND status = 'OPEN'
		`,
		roundID,
	)

	if err != nil {
		return fmt.Errorf(
			"lock round: %w",
			err,
		)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf(
			"round cannot be locked",
		)
	}

	return nil
}

func (r *RoundRepository) GetByID(
	ctx context.Context,
	id uuid.UUID,
) (*model.TournamentRound, error) {

	var round model.TournamentRound

	err := r.db.QueryRow(
		ctx,
		`
		SELECT
			id,
			event_id,
			round_number,
			round_name,
			pairing_method,
			status,
			locked_at,
			completed_at
		FROM tournament_rounds
		WHERE id = $1
		`,
		id,
	).Scan(
		&round.ID,
		&round.EventID,
		&round.RoundNumber,
		&round.RoundName,
		&round.PairingMethod,
		&round.Status,
		&round.LockedAt,
		&round.CompletedAt,
	)

	if err != nil {
		return nil, fmt.Errorf(
			"get round: %w",
			err,
		)
	}

	return &round, nil
}