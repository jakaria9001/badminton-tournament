package tests

import (
	"math/rand"
	"testing"

	"github.com/jakaria9001/badminton-tournament/backend/internal/draw"
	"github.com/jakaria9001/badminton-tournament/backend/internal/model"
)

func makeTeams(n int) []model.TeamSummary {

	teams := make(
		[]model.TeamSummary,
		n,
	)

	for i := 0; i < n; i++ {

		teams[i] = model.TeamSummary{
			ID:       string(rune('A' + i)),
			TeamName: "Team " + string(rune('A'+i)),
		}
	}

	return teams
}

func TestPairTeamsEven(t *testing.T) {

	teams := makeTeams(20)

	rng := rand.New(
		rand.NewSource(42),
	)

	pairings, bye :=
		draw.PairTeams(
			teams,
			true,
			rng,
		)

	if len(pairings) != 10 {
		t.Fatalf(
			"expected 10 pairings, got %d",
			len(pairings),
		)
	}

	if bye != nil {
		t.Fatal(
			"expected no bye",
		)
	}
}

func TestPairTeamsOdd(t *testing.T) {

	teams := makeTeams(5)

	rng := rand.New(
		rand.NewSource(42),
	)

	pairings, bye :=
		draw.PairTeams(
			teams,
			true,
			rng,
		)

	if len(pairings) != 2 {
		t.Fatalf(
			"expected 2 pairings, got %d",
			len(pairings),
		)
	}

	if bye == nil {
		t.Fatal(
			"expected a bye",
		)
	}
}

func TestThreeTeamSemiFinal(t *testing.T) {

	teams := makeTeams(3)

	rng := rand.New(
		rand.NewSource(42),
	)

	result, err :=
		draw.GenerateThreeTeamSemiFinal(
			teams,
			rng,
		)

	if err != nil {
		t.Fatal(err)
	}

	if result.SF1Team1 == nil ||
		result.SF1Team2 == nil ||
		result.WaitingTeam == nil {

		t.Fatal(
			"all three teams must be assigned",
		)
	}
}