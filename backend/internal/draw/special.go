package draw

import (
	"fmt"
	"math/rand"

	"github.com/jakaria9001/badminton-tournament/backend/internal/model"
)

type SpecialSemiResult struct {
	SF1Team1 *model.TeamSummary
	SF1Team2 *model.TeamSummary

	WaitingTeam *model.TeamSummary
}

func GenerateThreeTeamSemiFinal(
	teams []model.TeamSummary,
	rng *rand.Rand,
) (*SpecialSemiResult, error) {

	if len(teams) != 3 {
		return nil, fmt.Errorf(
			"special semifinal requires exactly 3 teams",
		)
	}

	working := make(
		[]model.TeamSummary,
		3,
	)

	copy(working, teams)

	rng.Shuffle(
		len(working),
		func(i, j int) {
			working[i], working[j] =
				working[j], working[i]
		},
	)

	return &SpecialSemiResult{
		SF1Team1: &working[0],
		SF1Team2: &working[1],
		WaitingTeam: &working[2],
	}, nil
}