package repository

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"

	"github.com/jakaria9001/badminton-tournament/backend/internal/model"
)

type MatchResultTarget struct {
	EventID uuid.UUID
	RoundID uuid.UUID
	Team1ID uuid.UUID
	Team2ID uuid.UUID
	Status  string
}

type EventScoring struct {
	BestOf        int
	WinningPoints int
	MaximumPoints int
}

func (r *MatchRepository) GetMatchResultTargetTx(
	ctx context.Context,
	tx pgx.Tx,
	matchID uuid.UUID,
) (*MatchResultTarget, error) {
	var target MatchResultTarget

	err := tx.QueryRow(ctx, `
		SELECT event_id, round_id, team1_id, team2_id, status
		FROM matches
		WHERE id = $1
		FOR UPDATE
	`, matchID).Scan(
		&target.EventID,
		&target.RoundID,
		&target.Team1ID,
		&target.Team2ID,
		&target.Status,
	)
	if err != nil {
		return nil, fmt.Errorf("get match for result: %w", err)
	}

	return &target, nil
}

func (r *MatchRepository) CompleteRoundIfFinishedTx(
	ctx context.Context,
	tx pgx.Tx,
	roundID uuid.UUID,
) error {
	_, err := tx.Exec(ctx, `
		UPDATE tournament_rounds
		SET status = 'COMPLETED', completed_at = NOW()
		WHERE id = $1
		  AND status IN ('OPEN', 'LOCKED')
		  AND EXISTS (
				SELECT 1 FROM matches
				WHERE round_id = $1
			)
		  AND NOT EXISTS (
				SELECT 1 FROM matches
				WHERE round_id = $1
				  AND status NOT IN ('COMPLETED', 'CANCELLED')
			)
	`, roundID)
	if err != nil {
		return fmt.Errorf("complete round: %w", err)
	}

	return nil
}

func (r *MatchRepository) GetEventScoringTx(
	ctx context.Context,
	tx pgx.Tx,
	eventID uuid.UUID,
) (*EventScoring, error) {
	var scoring EventScoring

	err := tx.QueryRow(ctx, `
		SELECT best_of, winning_points, maximum_points
		FROM events
		WHERE id = $1
	`, eventID).Scan(
		&scoring.BestOf,
		&scoring.WinningPoints,
		&scoring.MaximumPoints,
	)
	if err != nil {
		return nil, fmt.Errorf("get event scoring: %w", err)
	}

	return &scoring, nil
}

func (r *MatchRepository) InsertGameTx(
	ctx context.Context,
	tx pgx.Tx,
	matchID uuid.UUID,
	game model.GameResult,
) error {
	_, err := tx.Exec(ctx, `
		INSERT INTO match_games (
			id, match_id, game_number, team1_score, team2_score
		)
		VALUES ($1, $2, $3, $4, $5)
	`, uuid.New(), matchID, game.GameNumber, game.Team1Score, game.Team2Score)
	if err != nil {
		return fmt.Errorf("insert game %d: %w", game.GameNumber, err)
	}

	return nil
}
