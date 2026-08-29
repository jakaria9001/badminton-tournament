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

func (r *MatchRepository) CreateThreeTeamPlayoffTx(
	ctx context.Context,
	tx pgx.Tx,
	roundID uuid.UUID,
	sourceMatchID uuid.UUID,
	loserID uuid.UUID,
) error {
	var waitingTeamID uuid.UUID
	err := tx.QueryRow(ctx, `
		SELECT team_id FROM round_advancements
		WHERE round_id = $1 AND advancement_type = 'BYE'
	`, roundID).Scan(&waitingTeamID)
	if err == pgx.ErrNoRows {
		return nil
	}
	if err != nil {
		return fmt.Errorf("get three-team waiting team: %w", err)
	}

	_, err = tx.Exec(ctx, `
		INSERT INTO matches (
			id, event_id, round_id, round, match_number, match_type,
			team1_id, team2_id, team1_source_match_id, team1_source_type,
			team2_source_type, status
		)
		SELECT $1, event_id, id, round_name, 2, 'NORMAL',
			$2, $3, $4, 'LOSER', 'BYE', 'SCHEDULED'
		FROM tournament_rounds
		WHERE id = $5 AND round_name LIKE 'SEMIFINAL%'
		ON CONFLICT (event_id, round, match_number) DO NOTHING
	`, uuid.New(), loserID, waitingTeamID, sourceMatchID, roundID)
	if err != nil {
		return fmt.Errorf("create three-team playoff: %w", err)
	}
	return nil
}

func (r *MatchRepository) CreateFinalFixtureIfReadyTx(
	ctx context.Context,
	tx pgx.Tx,
	roundID uuid.UUID,
) error {
	var eventID uuid.UUID
	var roundNumber int
	var roundName string
	err := tx.QueryRow(ctx, `
		SELECT event_id, round_number, round_name
		FROM tournament_rounds WHERE id = $1 FOR UPDATE
	`, roundID).Scan(&eventID, &roundNumber, &roundName)
	if err != nil {
		return fmt.Errorf("lock semifinal round: %w", err)
	}
	if roundName != "SEMIFINAL" && roundName != "SEMIFINAL_1" && roundName != "SEMIFINAL_2" {
		return nil
	}

	var matchCount, incompleteCount int
	if err := tx.QueryRow(ctx, `
		SELECT COUNT(*), COUNT(*) FILTER (WHERE status <> 'COMPLETED')
		FROM matches WHERE round_id = $1
	`, roundID).Scan(&matchCount, &incompleteCount); err != nil {
		return fmt.Errorf("check semifinal completion: %w", err)
	}
	if matchCount < 2 || incompleteCount != 0 {
		return nil
	}

	rows, err := tx.Query(ctx, `
		SELECT winner_team_id FROM matches
		WHERE round_id = $1 AND status = 'COMPLETED'
		ORDER BY match_number
	`, roundID)
	if err != nil {
		return fmt.Errorf("get semifinal winners: %w", err)
	}
	defer rows.Close()

	winners := make([]uuid.UUID, 0, 2)
	for rows.Next() {
		var winnerID uuid.UUID
		if err := rows.Scan(&winnerID); err != nil {
			return fmt.Errorf("scan semifinal winner: %w", err)
		}
		winners = append(winners, winnerID)
	}
	if err := rows.Err(); err != nil {
		return fmt.Errorf("iterate semifinal winners: %w", err)
	}
	if len(winners) != 2 {
		return fmt.Errorf("a final requires exactly two semifinal winners")
	}

	finalRoundID := uuid.New()
	err = tx.QueryRow(ctx, `
		INSERT INTO tournament_rounds (id, event_id, round_number, round_name, pairing_method, status)
		VALUES ($1, $2, $3, 'FINAL', 'MANUAL', 'OPEN')
		ON CONFLICT (event_id, round_number) DO UPDATE SET id = tournament_rounds.id
		RETURNING id
	`, finalRoundID, eventID, roundNumber+1).Scan(&finalRoundID)
	if err != nil {
		return fmt.Errorf("create final round: %w", err)
	}

	_, err = tx.Exec(ctx, `
		INSERT INTO matches (
			id, event_id, round_id, round, match_number, match_type,
			team1_id, team2_id, status
		)
		VALUES ($1, $2, $3, 'FINAL', 1, 'NORMAL', $4, $5, 'SCHEDULED')
		ON CONFLICT (event_id, round, match_number) DO NOTHING
	`, uuid.New(), eventID, finalRoundID, winners[0], winners[1])
	if err != nil {
		return fmt.Errorf("create final fixture: %w", err)
	}
	return nil
}
