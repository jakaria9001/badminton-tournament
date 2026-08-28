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

	return s.repository.Lock(
		ctx,
		roundID,
	)
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