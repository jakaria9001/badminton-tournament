package draw

import (
	"math/rand"
	"time"

	"github.com/jakaria9001/badminton-tournament/backend/internal/model"
)

type Generator struct {
	random *rand.Rand
}

func NewGenerator() *Generator {
	return &Generator{
		random: rand.New(
			rand.NewSource(time.Now().UnixNano()),
		),
	}
}

func (g *Generator) ShuffleTeams(
	teams []model.TeamSummary,
) []model.TeamSummary {

	result := make(
		[]model.TeamSummary,
		len(teams),
	)

	copy(result, teams)

	g.random.Shuffle(
		len(result),
		func(i, j int) {
			result[i], result[j] =
				result[j], result[i]
		},
	)

	return result
}

func (g *Generator) PairTeams(
	teams []model.TeamSummary,
	shuffle bool,
) ([]PairingResult, *model.TeamSummary) {

	return PairTeams(
		teams,
		shuffle,
		g.random,
	)
}

func (g *Generator) GenerateThreeTeamSemiFinal(
	teams []model.TeamSummary,
) (*SpecialSemiResult, error) {
	return GenerateThreeTeamSemiFinal(
		teams,
		g.random,
	)
}
