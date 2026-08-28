package service

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/jakaria9001/badminton-tournament/backend/internal/draw"
	"github.com/jakaria9001/badminton-tournament/backend/internal/model"
	"github.com/jakaria9001/badminton-tournament/backend/internal/repository"
)

type DrawService struct {
	db                    *pgxpool.Pool
	roundRepository       *repository.RoundRepository
	matchRepository       *repository.MatchRepository
	advancementRepository *repository.AdvancementRepository
	generator             *draw.Generator
}

func NewDrawService(
	db *pgxpool.Pool,
	roundRepository *repository.RoundRepository,
	matchRepository *repository.MatchRepository,
	advancementRepository *repository.AdvancementRepository,
	generator *draw.Generator,
) *DrawService {

	return &DrawService{
		db:                    db,
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

	tx, err := s.db.Begin(ctx)

	if err != nil {
		return err
	}

	defer tx.Rollback(ctx)

	count, err :=
		s.matchRepository.CountByRoundTx(
			ctx,
			tx,
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
		s.roundRepository.GetByIDTx(
			ctx,
			tx,
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
		err := s.generateSpecialSemiFinal(
			ctx,
			round,
			teams,
		)
		if err != nil {
			return err
		}

		return tx.Commit(ctx)
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

	if err := s.persistPairings(
		ctx,
		round,
		pairings,
	); err != nil {
		return err
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf(
			"commit draw: %w",
			err,
		)
	}

	return nil
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
