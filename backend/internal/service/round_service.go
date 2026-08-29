package service

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"github.com/jakaria9001/badminton-tournament/backend/internal/model"
	"github.com/jakaria9001/badminton-tournament/backend/internal/repository"
)

type RoundService struct {
	repository      *repository.RoundRepository
	matchRepository *repository.MatchRepository
}

func NewRoundService(
	repository *repository.RoundRepository,
	matchRepository *repository.MatchRepository,
) *RoundService {
	return &RoundService{
		repository: repository,
		matchRepository: matchRepository,
	}
}

func (s *RoundService) Create(
	ctx context.Context,
	eventID uuid.UUID,
	roundNumber int,
	roundName string,
	pairingMethod string,
) (*model.TournamentRound, error) {

	if pairingMethod != model.PairingManual &&
		pairingMethod != model.PairingRandom {

		return nil, fmt.Errorf(
			"invalid pairing method",
		)
	}

	rounds, err := s.repository.GetByEvent(ctx, eventID)
	if err != nil {
		return nil, err
	}

	expectedRoundNumber := len(rounds) + 1
	if roundNumber != expectedRoundNumber {
		return nil, fmt.Errorf(
			"next round must be round %d",
			expectedRoundNumber,
		)
	}

	if len(rounds) > 0 && rounds[len(rounds)-1].Status != model.RoundCompleted {
		return nil, fmt.Errorf(
			"round %d must be COMPLETED before creating the next round",
			rounds[len(rounds)-1].RoundNumber,
		)
	}

	teamCount, err := s.matchRepository.CountConfirmedTeamsByEvent(ctx, eventID)
	if err != nil {
		return nil, err
	}
	if teamCount < 2 {
		return nil, fmt.Errorf(
			"at least 2 confirmed teams are required to create rounds",
		)
	}

	maxRounds := 0
	bracketSize := 1
	for bracketSize < teamCount {
		bracketSize *= 2
		maxRounds++
	}
	if roundNumber > maxRounds {
		return nil, fmt.Errorf(
			"cannot create round %d: this event supports a maximum of %d rounds",
			roundNumber,
			maxRounds,
		)
	}

	if roundNumber == maxRounds {
		roundName = "FINAL"
	} else if roundNumber == maxRounds-1 {
		roundName = "SEMIFINAL"
	} else {
		roundName = fmt.Sprintf("ROUND_%d", roundNumber)
	}

	return s.repository.Create(
		ctx,
		eventID,
		roundNumber,
		roundName,
		pairingMethod,
	)
}

func (s *RoundService) GetByEvent(
	ctx context.Context,
	eventID uuid.UUID,
) ([]model.TournamentRound, error) {

	return s.repository.GetByEvent(
		ctx,
		eventID,
	)
}

func (s *RoundService) Lock(
	ctx context.Context,
	roundID uuid.UUID,
) error {
	round, err := s.repository.GetByID(ctx, roundID)
	if err != nil {
		return err
	}
	if !CanLockRound(round.RoundName) {
		return fmt.Errorf("only the final round can be locked")
	}

	return s.repository.Lock(
		ctx,
		roundID,
	)
}

func CanLockRound(roundName string) bool {
	return roundName == model.Final
}

func (s *RoundService) GetAvailableTeams(
	ctx context.Context,
	roundID uuid.UUID,
) ([]model.TeamSummary, error) {

	return s.matchRepository.GetTeamsForRound(
		ctx,
		roundID,
	)
}