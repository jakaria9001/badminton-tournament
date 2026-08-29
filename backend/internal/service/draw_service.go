package service

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
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

func shouldGenerateDraw(round *model.TournamentRound) bool {
	return round != nil && round.Status == model.RoundOpen
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

	round, err :=
		s.roundRepository.GetByIDTx(
			ctx,
			tx,
			roundID,
		)

	if err != nil {
		return err
	}

	if !shouldGenerateDraw(round) {
		return fmt.Errorf(
			"round must be OPEN",
		)
	}

	if err := s.matchRepository.DeleteByRoundTx(
		ctx,
		tx,
		roundID,
	); err != nil {
		return err
	}

	if err := s.advancementRepository.DeleteByRoundTx(
		ctx,
		tx,
		roundID,
	); err != nil {
		return err
	}

	teams, err :=
		s.matchRepository.GetTeamsForRoundTx(
			ctx,
			tx,
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
			tx,
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
		err := s.advancementRepository.CreateByeTx(
			ctx,
			tx,
			uuid.MustParse(round.ID),
			uuid.MustParse(bye.ID),
		)

		if err != nil {
			return err
		}
	}

	if err := s.persistPairingsTx(
		ctx,
		tx,
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
	tx pgx.Tx,
	round *model.TournamentRound,
	teams []model.TeamSummary,
) error {
	result, err := s.generator.GenerateThreeTeamSemiFinal(teams)
	if err != nil {
		return err
	}

	if err := s.advancementRepository.CreateByeTx(
		ctx,
		tx,
		uuid.MustParse(round.ID),
		uuid.MustParse(result.WaitingTeam.ID),
	); err != nil {
		return err
	}

	return s.persistPairingsTx(
		ctx,
		tx,
		round,
		[]draw.PairingResult{{
			Team1: result.SF1Team1,
			Team2: result.SF1Team2,
		}},
	)
}

func (s *DrawService) persistPairingsTx(
	ctx context.Context,
	tx pgx.Tx,
	round *model.TournamentRound,
	pairings []draw.PairingResult,
) error {

	for i, pairing := range pairings {
		team1ID := uuid.MustParse(pairing.Team1.ID)
		team2ID := uuid.MustParse(pairing.Team2.ID)

		_, err :=
			s.matchRepository.CreateMatchTx(
				ctx,
				tx,
				uuid.MustParse(round.ID),
				i+1,
				&team1ID,
				&team2ID,
				model.MatchNormal,
			)

		if err != nil {
			return err
		}
	}

	return nil
}
