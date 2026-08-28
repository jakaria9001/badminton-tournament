package draw

import (
	"math/rand"

	"github.com/jakaria9001/badminton-tournament/backend/internal/model"
)

func SelectBye(
	teams []model.TeamSummary,
	rng *rand.Rand,
) (*model.TeamSummary, []model.TeamSummary) {

	if len(teams) == 0 {
		return nil, teams
	}

	index := rng.Intn(len(teams))

	bye := teams[index]

	remaining := make(
		[]model.TeamSummary,
		0,
		len(teams)-1,
	)

	remaining = append(
		remaining,
		teams[:index]...,
	)

	remaining = append(
		remaining,
		teams[index+1:]...,
	)

	return &bye, remaining
}