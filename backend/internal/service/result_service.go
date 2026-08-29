package service

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/jakaria9001/badminton-tournament/backend/internal/model"
	"github.com/jakaria9001/badminton-tournament/backend/internal/repository"
)

type ResultService struct {
	db              *pgxpool.Pool
	matchRepository *repository.MatchRepository
}

func NewResultService(
	db *pgxpool.Pool,
	matchRepository *repository.MatchRepository,
) *ResultService {
	return &ResultService{db: db, matchRepository: matchRepository}
}

func (s *ResultService) Submit(
	ctx context.Context,
	matchID uuid.UUID,
	request model.SubmitMatchResultRequest,
) error {
	if err := validateGames(request.Games); err != nil {
		return err
	}

	tx, err := s.db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("begin result transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	target, err := s.matchRepository.GetMatchResultTargetTx(ctx, tx, matchID)
	if err != nil {
		return err
	}
	if target.Status != model.MatchScheduled && target.Status != model.MatchInProgress {
		return fmt.Errorf("match cannot accept a result")
	}

	scoring, err := s.matchRepository.GetEventScoringTx(ctx, tx, target.EventID)
	if err != nil {
		return err
	}
	winner, loser, err := determineWinner(target.Team1ID, target.Team2ID, request.Games, scoring)
	if err != nil {
		return err
	}

	if request.WinnerTeamID != "" {
		suppliedWinner, err := uuid.Parse(request.WinnerTeamID)
		if err != nil || suppliedWinner != winner {
			return fmt.Errorf("winnerTeamId does not match the game scores")
		}
	}

	for _, game := range request.Games {
		if err := s.matchRepository.InsertGameTx(ctx, tx, matchID, game); err != nil {
			return err
		}
	}

	if err := s.matchRepository.CompleteMatchTx(ctx, tx, matchID, winner, loser); err != nil {
		return err
	}
	if err := s.matchRepository.CreateThreeTeamPlayoffTx(ctx, tx, target.RoundID, matchID, loser); err != nil {
		return err
	}
	if err := s.matchRepository.CreateFinalFixtureIfReadyTx(ctx, tx, target.RoundID); err != nil {
		return err
	}

	if err := s.matchRepository.CompleteRoundIfFinishedTx(ctx, tx, target.RoundID); err != nil {
		return err
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("commit match result: %w", err)
	}
	return nil
}

func validateGames(games []model.GameResult) error {
	if len(games) == 0 {
		return fmt.Errorf("at least one game result is required")
	}

	seen := make(map[int]bool, len(games))
	for index, game := range games {
		if game.GameNumber != index+1 || seen[game.GameNumber] {
			return fmt.Errorf("game numbers must be sequential starting at 1")
		}
		seen[game.GameNumber] = true
		if game.Team1Score < 0 || game.Team2Score < 0 {
			return fmt.Errorf("scores cannot be negative")
		}
		if game.Team1Score == game.Team2Score {
			return fmt.Errorf("a game cannot end in a tie")
		}
	}
	return nil
}

func determineWinner(
	team1ID uuid.UUID,
	team2ID uuid.UUID,
	games []model.GameResult,
	scoring *repository.EventScoring,
) (uuid.UUID, uuid.UUID, error) {
	if scoring.BestOf != 1 && scoring.BestOf != 3 {
		return uuid.Nil, uuid.Nil, fmt.Errorf("invalid event best_of")
	}

	team1Wins, team2Wins := 0, 0
	for _, game := range games {
		if game.Team1Score > game.Team2Score {
			if err := validateGameScore(game.Team1Score, game.Team2Score, scoring); err != nil {
				return uuid.Nil, uuid.Nil, err
			}
			team1Wins++
		} else {
			if err := validateGameScore(game.Team2Score, game.Team1Score, scoring); err != nil {
				return uuid.Nil, uuid.Nil, err
			}
			team2Wins++
		}
	}

	winsNeeded := scoring.BestOf/2 + 1
	if (team1Wins >= winsNeeded && team2Wins >= winsNeeded) || team1Wins < winsNeeded && team2Wins < winsNeeded {
		return uuid.Nil, uuid.Nil, fmt.Errorf("games do not determine a match winner")
	}
	if len(games) > scoring.BestOf {
		return uuid.Nil, uuid.Nil, fmt.Errorf("too many games for best-of-%d match", scoring.BestOf)
	}
	if team1Wins >= winsNeeded {
		return team1ID, team2ID, nil
	}
	return team2ID, team1ID, nil
}

func validateGameScore(winnerScore, loserScore int, scoring *repository.EventScoring) error {
	if winnerScore < scoring.WinningPoints || winnerScore > scoring.MaximumPoints {
		return fmt.Errorf("winning score must be between %d and %d", scoring.WinningPoints, scoring.MaximumPoints)
	}
	if winnerScore-loserScore < 2 {
		return fmt.Errorf("a game must be won by two points unless the maximum score is reached")
	}
	if winnerScore > scoring.WinningPoints && loserScore < scoring.WinningPoints-1 {
		return fmt.Errorf("scores above the winning points must follow deuce scoring")
	}
	if winnerScore == scoring.MaximumPoints && loserScore != scoring.MaximumPoints-1 {
		return fmt.Errorf("a game at the maximum score must be won %d-%d", scoring.MaximumPoints, scoring.MaximumPoints-1)
	}
	if loserScore >= scoring.MaximumPoints || loserScore >= winnerScore {
		return fmt.Errorf("invalid losing score")
	}
	return nil
}
