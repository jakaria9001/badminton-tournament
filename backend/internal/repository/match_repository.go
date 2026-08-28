package repository

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/jakaria9001/badminton-tournament/backend/internal/model"
)

type MatchRepository struct {
	db *pgxpool.Pool
}

func NewMatchRepository(
	db *pgxpool.Pool,
) *MatchRepository {
	return &MatchRepository{
		db: db,
	}
}

func (r *MatchRepository) GetMatchesByEvent(
	ctx context.Context,
	eventID uuid.UUID,
) ([]model.MatchResponse, error) {

	query := `
		SELECT
			m.id,
			m.event_id,
			m.round_id,
			m.round,
			m.match_number,
			m.match_type,
			m.team1_id,
			t1.team_name,
			p1.name,
			p2.name,
			m.team2_id,
			t2.team_name,
			p3.name,
			p4.name,
			m.court_number,
			m.scheduled_at,
			m.status,
			m.winner_team_id,
			m.loser_team_id,
			m.winner_next_match_id,
			m.loser_next_match_id,
			m.team1_source_match_id,
			m.team1_source_type,
			m.team2_source_match_id,
			m.team2_source_type
		FROM matches m

		LEFT JOIN teams t1
			ON t1.id = m.team1_id

		LEFT JOIN players p1
			ON p1.id = t1.player1_id

		LEFT JOIN players p2
			ON p2.id = t1.player2_id

		LEFT JOIN teams t2
			ON t2.id = m.team2_id

		LEFT JOIN players p3
			ON p3.id = t2.player1_id

		LEFT JOIN players p4
			ON p4.id = t2.player2_id

		WHERE m.event_id = $1

		ORDER BY
			CASE m.round
				WHEN 'ROUND_1' THEN 1
				WHEN 'ROUND_2' THEN 2
				WHEN 'ROUND_3' THEN 3
				WHEN 'SEMI_FINAL_1' THEN 4
				WHEN 'SEMI_FINAL_2' THEN 5
				WHEN 'FINAL' THEN 6
				ELSE 99
			END,
			m.match_number
	`

	rows, err := r.db.Query(
		ctx,
		query,
		eventID,
	)

	if err != nil {
		return nil, fmt.Errorf(
			"get matches: %w",
			err,
		)
	}

	defer rows.Close()

	matches := make([]model.MatchResponse, 0)

	for rows.Next() {

		var match model.MatchResponse

		var (
			team1ID   *uuid.UUID
			team1Name *string
			player1   *string
			player2   *string

			team2ID   *uuid.UUID
			team2Name *string
			player3   *string
			player4   *string
		)

		err := rows.Scan(
			&match.ID,
			&match.EventID,
			&match.RoundID,
			&match.Round,
			&match.MatchNumber,
			&match.MatchType,

			&team1ID,
			&team1Name,
			&player1,
			&player2,

			&team2ID,
			&team2Name,
			&player3,
			&player4,

			&match.CourtNumber,
			&match.ScheduledAt,
			&match.Status,
			&match.WinnerTeamID,
			&match.LoserTeamID,
			&match.WinnerNextMatchID,
			&match.LoserNextMatchID,
			&match.Team1SourceMatchID,
			&match.Team1SourceType,
			&match.Team2SourceMatchID,
			&match.Team2SourceType,
		)

		if err != nil {
			return nil, fmt.Errorf(
				"scan match: %w",
				err,
			)
		}

		if team1ID != nil {
			match.Team1 = &model.TeamSummary{
				ID:       team1ID.String(),
				TeamName: valueOrEmpty(team1Name),
				Player1:  valueOrEmpty(player1),
				Player2:  valueOrEmpty(player2),
			}
		}

		if team2ID != nil {
			match.Team2 = &model.TeamSummary{
				ID:       team2ID.String(),
				TeamName: valueOrEmpty(team2Name),
				Player1:  valueOrEmpty(player3),
				Player2:  valueOrEmpty(player4),
			}
		}

		matches = append(matches, match)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf(
			"iterate matches: %w",
			err,
		)
	}

	return matches, nil
}

func valueOrEmpty(value *string) string {
	if value == nil {
		return ""
	}

	return *value
}

// func (r *MatchRepository) Create(
// 	ctx context.Context,
// 	roundID uuid.UUID,
// 	matchNumber int,
// 	team1ID *uuid.UUID,
// 	team2ID *uuid.UUID,
// 	matchType string,
// ) (*model.MatchResponse, error) {

// 	id := uuid.New()

// 	_, err := r.db.Exec(
// 		ctx,
// 		`
// 		INSERT INTO matches (
// 			id,
// 			event_id,
// 			round_id,
// 			round,
// 			match_number,
// 			match_type,
// 			team1_id,
// 			team2_id,
// 			status
// 		)
// 		SELECT
// 			$1,
// 			event_id,
// 			$2,
// 			round_name,
// 			$3,
// 			$4,
// 			$5,
// 			$6,
// 			'SCHEDULED'
// 		FROM tournament_rounds
// 		WHERE id = $2
// 		`,
// 		id,
// 		roundID,
// 		matchNumber,
// 		matchType,
// 		team1ID,
// 		team2ID,
// 	)

// 	if err != nil {
// 		return nil, fmt.Errorf(
// 			"create match: %w",
// 			err,
// 		)
// 	}

// 	return nil, nil
// }

func (r *MatchRepository) GetTeamsForRound(
	ctx context.Context,
	roundID uuid.UUID,
) ([]model.TeamSummary, error) {

	var roundNumber int
	var eventID uuid.UUID

	err := r.db.QueryRow(
		ctx,
		`
		SELECT
			round_number,
			event_id
		FROM tournament_rounds
		WHERE id = $1
		`,
		roundID,
	).Scan(
		&roundNumber,
		&eventID,
	)

	if err != nil {
		return nil, fmt.Errorf(
			"get round: %w",
			err,
		)
	}

	if roundNumber == 1 {
		return r.getConfirmedTeams(
			ctx,
			eventID,
		)
	}

	return r.getPreviousRoundWinners(
		ctx,
		eventID,
		roundNumber,
	)
}

func (r *MatchRepository) getConfirmedTeams(
	ctx context.Context,
	eventID uuid.UUID,
) ([]model.TeamSummary, error) {

	rows, err := r.db.Query(
		ctx,
		`
		SELECT
			t.id,
			COALESCE(t.team_name, ''),
			p1.name,
			p2.name
		FROM teams t
		JOIN players p1
			ON p1.id = t.player1_id
		JOIN players p2
			ON p2.id = t.player2_id
		WHERE t.event_id = $1
		  AND t.status = 'CONFIRMED'
		ORDER BY t.created_at
		`,
		eventID,
	)

	if err != nil {
		return nil, fmt.Errorf(
			"get confirmed teams: %w",
			err,
		)
	}

	defer rows.Close()

	teams := make([]model.TeamSummary, 0)

	for rows.Next() {

		var team model.TeamSummary

		err := rows.Scan(
			&team.ID,
			&team.TeamName,
			&team.Player1,
			&team.Player2,
		)

		if err != nil {
			return nil, fmt.Errorf(
				"scan team: %w",
				err,
			)
		}

		teams = append(
			teams,
			team,
		)
	}

	return teams, rows.Err()
}

func (r *MatchRepository) getPreviousRoundWinners(
	ctx context.Context,
	eventID uuid.UUID,
	roundNumber int,
) ([]model.TeamSummary, error) {

	rows, err := r.db.Query(
		ctx,
		`
		SELECT
			t.id,
			COALESCE(t.team_name, ''),
			p1.name,
			p2.name
		FROM matches m
		JOIN teams t
			ON t.id = m.winner_team_id
		JOIN players p1
			ON p1.id = t.player1_id
		JOIN players p2
			ON p2.id = t.player2_id
		WHERE m.event_id = $1
		  AND m.round = (
		      SELECT round_name
		      FROM tournament_rounds
		      WHERE event_id = $1
		        AND round_number = $2 - 1
		  )
		  AND m.status = 'COMPLETED'
		  AND m.winner_team_id IS NOT NULL
		ORDER BY m.match_number
		`,
		eventID,
		roundNumber,
	)

	if err != nil {
		return nil, fmt.Errorf(
			"get previous round winners: %w",
			err,
		)
	}

	defer rows.Close()

	teams := make([]model.TeamSummary, 0)

	for rows.Next() {

		var team model.TeamSummary

		err := rows.Scan(
			&team.ID,
			&team.TeamName,
			&team.Player1,
			&team.Player2,
		)

		if err != nil {
			return nil, fmt.Errorf(
				"scan winner: %w",
				err,
			)
		}

		teams = append(
			teams,
			team,
		)
	}

	return teams, rows.Err()
}

func (r *MatchRepository) CreateMatch(
	ctx context.Context,
	roundID uuid.UUID,
	matchNumber int,
	team1ID uuid.UUID,
	team2ID uuid.UUID,
	matchType string,
) (uuid.UUID, error) {

	id := uuid.New()

	var createdID uuid.UUID

	err := r.db.QueryRow(
		ctx,
		`
		INSERT INTO matches (
			id,
			event_id,
			round_id,
			round,
			match_number,
			match_type,
			team1_id,
			team2_id,
			status
		)
		SELECT
			$1,
			tournament_rounds.event_id,
			$2,
			round_name,
			$3,
			$4,
			$5,
			$6,
			'SCHEDULED'
		FROM tournament_rounds
		JOIN teams t1
			ON t1.id = $5
		JOIN teams t2
			ON t2.id = $6
		WHERE tournament_rounds.id = $2
		  AND tournament_rounds.status = 'OPEN'
		  AND t1.event_id = tournament_rounds.event_id
		  AND t2.event_id = tournament_rounds.event_id
		  AND t1.status = 'CONFIRMED'
		  AND t2.status = 'CONFIRMED'
		  AND NOT EXISTS (
				SELECT 1
				FROM matches existing_match
				WHERE existing_match.round_id = $2
				  AND (
						existing_match.team1_id IN ($5, $6)
						OR existing_match.team2_id IN ($5, $6)
				  )
			)
		RETURNING id
		`,
		id,
		roundID,
		matchNumber,
		matchType,
		team1ID,
		team2ID,
	).Scan(&createdID)

	if err != nil {
		if err == pgx.ErrNoRows {
			return uuid.Nil, fmt.Errorf(
				"teams must be confirmed, belong to the round event, and not already participate in this round",
			)
		}

		return uuid.Nil, fmt.Errorf(
			"create match: %w",
			err,
		)
	}

	return createdID, nil
}

func (r *MatchRepository) CreateMatchTx(
	ctx context.Context,
	tx pgx.Tx,
	roundID uuid.UUID,
	matchNumber int,
	team1ID *uuid.UUID,
	team2ID *uuid.UUID,
	matchType string,
) (uuid.UUID, error) {

	id := uuid.New()

	_, err := tx.Exec(
		ctx,
		`
		INSERT INTO matches (
			id,
			event_id,
			round_id,
			round,
			match_number,
			match_type,
			team1_id,
			team2_id,
			status
		)
		SELECT
			$1,
			event_id,
			$2,
			round_name,
			$3,
			$4,
			$5,
			$6,
			'SCHEDULED'
		FROM tournament_rounds
		WHERE id = $2
		`,
		id,
		roundID,
		matchNumber,
		matchType,
		team1ID,
		team2ID,
	)

	if err != nil {
		return uuid.Nil, fmt.Errorf(
			"create match tx: %w",
			err,
		)
	}

	return id, nil
}

func (r *MatchRepository) CreateMatchSlotTx(
	ctx context.Context,
	tx pgx.Tx,
	roundID uuid.UUID,
	matchNumber int,
	matchType string,
	team1ID *uuid.UUID,
	team2ID *uuid.UUID,
	team1SourceMatchID *uuid.UUID,
	team1SourceType *string,
	team2SourceMatchID *uuid.UUID,
	team2SourceType *string,
) (uuid.UUID, error) {
	id := uuid.New()

	var createdID uuid.UUID

	err := tx.QueryRow(
		ctx,
		`
		INSERT INTO matches (
			id,
			event_id,
			round_id,
			round,
			match_number,
			match_type,
			team1_id,
			team2_id,
			team1_source_match_id,
			team1_source_type,
			team2_source_match_id,
			team2_source_type,
			status
		)
		SELECT
			$1,
			event_id,
			$2,
			round_name,
			$3,
			$4,
			$5,
			$6,
			$7,
			$8,
			$9,
			$10,
			'SCHEDULED'
		FROM tournament_rounds
		WHERE id = $2
		RETURNING id
		`,
		id,
		roundID,
		matchNumber,
		matchType,
		team1ID,
		team2ID,
		team1SourceMatchID,
		team1SourceType,
		team2SourceMatchID,
		team2SourceType,
	).Scan(&createdID)

	if err != nil {
		return uuid.Nil, fmt.Errorf(
			"create match slot tx: %w",
			err,
		)
	}

	return createdID, nil
}

func (r *MatchRepository) CountByRoundTx(
	ctx context.Context,
	tx pgx.Tx,
	roundID uuid.UUID,
) (int, error) {

	var count int

	err := tx.QueryRow(
		ctx,
		`
		SELECT COUNT(*)
		FROM matches
		WHERE round_id = $1
		`,
		roundID,
	).Scan(&count)

	return count, err
}

func (r *MatchRepository) CompleteMatchTx(
	ctx context.Context,
	tx pgx.Tx,
	matchID uuid.UUID,
	winnerID uuid.UUID,
	loserID uuid.UUID,
) error {

	result, err := tx.Exec(
		ctx,
		`
		UPDATE matches
		SET
			status = 'COMPLETED',
			winner_team_id = $2,
			loser_team_id = $3,
			updated_at = NOW()
		WHERE id = $1
			AND status IN ('SCHEDULED', 'IN_PROGRESS')
		`,
		matchID,
		winnerID,
		loserID,
	)

	if err != nil {
		return err
	}

	if result.RowsAffected() != 1 {
		return fmt.Errorf(
			"match cannot be completed",
		)
	}

	return nil
}

func (r *MatchRepository) GetNextMatchNumber(
	ctx context.Context,
	roundID uuid.UUID,
) (int, error) {

	var number int

	err := r.db.QueryRow(
		ctx,
		`
		SELECT COALESCE(
			MAX(match_number),
			0
		) + 1
		FROM matches
		WHERE round_id = $1
		`,
		roundID,
	).Scan(&number)

	return number, err
}

func (r *MatchRepository) CountByRound(
	ctx context.Context,
	roundID uuid.UUID,
) (int, error) {

	var count int

	err := r.db.QueryRow(
		ctx,
		`
		SELECT COUNT(*)
		FROM matches
		WHERE round_id = $1
		`,
		roundID,
	).Scan(&count)

	return count, err
}

func (r *MatchRepository) CountConfirmedTeamsByEvent(
	ctx context.Context,
	eventID uuid.UUID,
) (int, error) {
	var count int

	err := r.db.QueryRow(
		ctx,
		`
		SELECT COUNT(*)
		FROM teams
		WHERE event_id = $1
		  AND status = 'CONFIRMED'
		`,
		eventID,
	).Scan(&count)

	if err != nil {
		return 0, fmt.Errorf(
			"count confirmed teams: %w",
			err,
		)
	}

	return count, nil
}
