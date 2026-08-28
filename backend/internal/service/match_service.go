package service

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"github.com/jakaria9001/badminton-tournament/backend/internal/model"
	"github.com/jakaria9001/badminton-tournament/backend/internal/repository"
)

type MatchService struct {
	matchRepository *repository.MatchRepository
	roundRepository *repository.RoundRepository
}

func NewMatchService(
	matchRepository *repository.MatchRepository,
	roundRepository *repository.RoundRepository,
) *MatchService {
	return &MatchService{
		matchRepository: matchRepository,
		roundRepository: roundRepository,
	}
}

func (s *MatchService) GetMatchesByEvent(
	ctx context.Context,
	eventID uuid.UUID,
) ([]model.MatchResponse, error) {

	return s.matchRepository.GetMatchesByEvent(
		ctx,
		eventID,
	)
}

func (s *MatchService) CreateMatch(
	ctx context.Context,
	roundID uuid.UUID,
	team1ID uuid.UUID,
	team2ID uuid.UUID,
) (uuid.UUID, error) {

	if team1ID == team2ID {
		return uuid.Nil, fmt.Errorf(
			"a team cannot play itself",
		)
	}

	round, err :=
		s.roundRepository.GetByID(
			ctx,
			roundID,
		)

	if err != nil {
		return uuid.Nil, err
	}

	if round.Status != model.RoundOpen {
		return uuid.Nil, fmt.Errorf(
			"round is not open",
		)
	}

	// Validate team eligibility here.

	matchNumber, err :=
		s.matchRepository.GetNextMatchNumber(
			ctx,
			roundID,
		)

	if err != nil {
		return uuid.Nil, err
	}

	return s.matchRepository.CreateMatch(
		ctx,
		roundID,
		matchNumber,
		team1ID,
		team2ID,
		model.MatchNormal,
	)
}