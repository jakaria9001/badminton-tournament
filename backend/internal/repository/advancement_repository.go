package repository

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type AdvancementRepository struct {
	db *pgxpool.Pool
}

func NewAdvancementRepository(
	db *pgxpool.Pool,
) *AdvancementRepository {

	return &AdvancementRepository{
		db: db,
	}
}

func (r *AdvancementRepository) CreateBye(
	ctx context.Context,
	roundID uuid.UUID,
	teamID uuid.UUID,
) error {

	_, err := r.db.Exec(
		ctx,
		`
		INSERT INTO round_advancements (
			id,
			round_id,
			team_id,
			advancement_type
		)
		VALUES (
			$1,
			$2,
			$3,
			'BYE'
		)
		`,
		uuid.New(),
		roundID,
		teamID,
	)

	if err != nil {
		return fmt.Errorf(
			"create bye advancement: %w",
			err,
		)
	}

	return nil
}

func (r *AdvancementRepository) CreateByeTx(
	ctx context.Context,
	tx pgx.Tx,
	roundID uuid.UUID,
	teamID uuid.UUID,
) error {

	_, err := tx.Exec(
		ctx,
		`
		INSERT INTO round_advancements (
			id,
			round_id,
			team_id,
			advancement_type
		)
		VALUES (
			$1,
			$2,
			$3,
			'BYE'
		)
		`,
		uuid.New(),
		roundID,
		teamID,
	)

	if err != nil {
		return fmt.Errorf(
			"create bye advancement tx: %w",
			err,
		)
	}

	return nil
}

func (r *AdvancementRepository) DeleteByRoundTx(
	ctx context.Context,
	tx pgx.Tx,
	roundID uuid.UUID,
) error {

	_, err := tx.Exec(
		ctx,
		`
		DELETE FROM round_advancements
		WHERE round_id = $1
		`,
		roundID,
	)

	if err != nil {
		return fmt.Errorf(
			"delete advancements by round: %w",
			err,
		)
	}

	return nil
}