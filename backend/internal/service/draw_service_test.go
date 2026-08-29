package service

import (
	"testing"

	"github.com/jakaria9001/badminton-tournament/backend/internal/model"
)

func TestShouldGenerateDraw_WhenRoundIsOpen(t *testing.T) {
	round := &model.TournamentRound{Status: model.RoundOpen}

	if !shouldGenerateDraw(round) {
		t.Fatal("expected OPEN rounds to allow generating a draw")
	}
}

func TestShouldRejectGenerateDraw_WhenRoundIsLocked(t *testing.T) {
	round := &model.TournamentRound{Status: model.RoundLocked}

	if shouldGenerateDraw(round) {
		t.Fatal("expected LOCKED rounds to reject drawing")
	}
}

func TestShouldRejectGenerateDraw_WhenRoundIsCompleted(t *testing.T) {
	round := &model.TournamentRound{Status: model.RoundCompleted}

	if shouldGenerateDraw(round) {
		t.Fatal("expected COMPLETED rounds to reject drawing")
	}
}
