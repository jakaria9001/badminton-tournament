package service

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"github.com/jakaria9001/badminton-tournament/backend/internal/draw"
	"github.com/jakaria9001/badminton-tournament/backend/internal/model"
	"github.com/jakaria9001/badminton-tournament/backend/internal/repository"
)

type DrawService struct {
	roundRepository       *repository.RoundRepository
	matchRepository       *repository.MatchRepository
	advancementRepository *repository.AdvancementRepository
	generator             *draw.Generator
}

func NewDrawService(
	roundRepository *repository.RoundRepository,
	matchRepository *repository.MatchRepository,
	advancementRepository *repository.AdvancementRepository,
	generator *draw.Generator,
) *DrawService {

	return &DrawService{
		roundRepository:       roundRepository,
		matchRepository:       matchRepository,
		advancementRepository: advancementRepository,
		generator:             generator,
	}
}

func (s *DrawService) Generate(
	ctx context.Context,
	roundID uuid.UUID,
) error {

	count, err :=
		s.matchRepository.CountByRound(
			ctx,
			roundID,
		)

	if err != nil {
		return err
	}

	if count > 0 {
		return fmt.Errorf(
			"draw has already been generated",
		)
	}

	round, err :=
		s.roundRepository.GetByID(
			ctx,
			roundID,
		)

	if err != nil {
		return err
	}

	if round.Status != model.RoundOpen {
		return fmt.Errorf(
			"round must be OPEN",
		)
	}

	teams, err :=
		s.matchRepository.GetTeamsForRound(
			ctx,
			roundID,
		)

	if err != nil {
		return err
	}

	if len(teams) < 2 {
		return fmt.Errorf(
			"not enough teams to generate draw",
		)
	}

	// Special three-team stage.
	if len(teams) == 3 {
		return s.generateSpecialSemiFinal(
			ctx,
			round,
			teams,
		)
	}

	if round.PairingMethod != model.PairingRandom {
		return fmt.Errorf(
			"round is not configured for random generation",
		)
	}

	pairings, bye :=
		s.generator.PairTeams(
			teams,
			true,
		)

	if bye != nil {
		err := s.advancementRepository.CreateBye(
			ctx,
			uuid.MustParse(round.ID),
			uuid.MustParse(bye.ID),
		)

		if err != nil {
			return err
		}
	}

	_ = bye

	return s.persistPairings(
		ctx,
		round,
		pairings,
	)
}

func (s *DrawService) generateSpecialSemiFinal(
	ctx context.Context,
	round *model.TournamentRound,
	teams []model.TeamSummary,
) error {
	result, err := s.generator.GenerateThreeTeamSemiFinal(teams)
	if err != nil {
		return err
	}

	if err := s.advancementRepository.CreateBye(
		ctx,
		uuid.MustParse(round.ID),
		uuid.MustParse(result.WaitingTeam.ID),
	); err != nil {
		return err
	}

	return s.persistPairings(
		ctx,
		round,
		[]draw.PairingResult{{
			Team1: result.SF1Team1,
			Team2: result.SF1Team2,
		}},
	)
}

func (s *DrawService) persistPairings(
	ctx context.Context,
	round *model.TournamentRound,
	pairings []draw.PairingResult,
) error {

	for i, pairing := range pairings {

		_, err :=
			s.matchRepository.CreateMatch(
				ctx,
				uuid.MustParse(round.ID),
				i+1,
				uuid.MustParse(pairing.Team1.ID),
				uuid.MustParse(pairing.Team2.ID),
				model.MatchNormal,
			)

		if err != nil {
			return err
		}
	}

	return nil
}
