package draw

import (
	"math/rand"

	"github.com/jakaria9001/badminton-tournament/backend/internal/model"
)

type PairingResult struct {
	Team1 *model.TeamSummary
	Team2 *model.TeamSummary
}

func PairTeams(
	teams []model.TeamSummary,
	shuffle bool,
	rng *rand.Rand,
) ([]PairingResult, *model.TeamSummary) {

	working := make(
		[]model.TeamSummary,
		len(teams),
	)

	copy(working, teams)

	if shuffle {
		rng.Shuffle(
			len(working),
			func(i, j int) {
				working[i], working[j] =
					working[j], working[i]
			},
		)
	}

	var bye *model.TeamSummary

	if len(working)%2 != 0 {

		bye, working =
			SelectBye(
				working,
				rng,
			)
	}

	pairings := make(
		[]PairingResult,
		0,
		len(working)/2,
	)

	for i := 0; i < len(working); i += 2 {

		team1 := working[i]
		team2 := working[i+1]

		pairings = append(
			pairings,
			PairingResult{
				Team1: &team1,
				Team2: &team2,
			},
		)
	}

	return pairings, bye
}