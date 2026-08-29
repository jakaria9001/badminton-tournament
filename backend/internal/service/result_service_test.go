package service

import (
	"testing"

	"github.com/google/uuid"

	"github.com/jakaria9001/badminton-tournament/backend/internal/model"
	"github.com/jakaria9001/badminton-tournament/backend/internal/repository"
)

func TestValidateGames(t *testing.T) {
	tests := []struct {
		name    string
		games   []model.GameResult
		wantErr bool
	}{
		{
			name: "valid best of three",
			games: []model.GameResult{
				{GameNumber: 1, Team1Score: 21, Team2Score: 19},
				{GameNumber: 2, Team1Score: 18, Team2Score: 21},
				{GameNumber: 3, Team1Score: 21, Team2Score: 15},
			},
		},
		{name: "no games", wantErr: true},
		{
			name:    "non-sequential games",
			games:   []model.GameResult{{GameNumber: 2, Team1Score: 21, Team2Score: 19}},
			wantErr: true,
		},
		{
			name:    "negative score",
			games:   []model.GameResult{{GameNumber: 1, Team1Score: -1, Team2Score: 0}},
			wantErr: true,
		},
		{
			name:    "tied game",
			games:   []model.GameResult{{GameNumber: 1, Team1Score: 21, Team2Score: 21}},
			wantErr: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := validateGames(test.games)
			if (err != nil) != test.wantErr {
				t.Fatalf("validateGames() error = %v, wantErr %v", err, test.wantErr)
			}
		})
	}
}

func TestDetermineWinner(t *testing.T) {
	team1 := uuid.New()
	team2 := uuid.New()
	scoring := &repository.EventScoring{
		BestOf:        3,
		WinningPoints: 21,
		MaximumPoints: 30,
	}

	winner, loser, err := determineWinner(team1, team2, []model.GameResult{
		{GameNumber: 1, Team1Score: 21, Team2Score: 19},
		{GameNumber: 2, Team1Score: 20, Team2Score: 22},
		{GameNumber: 3, Team1Score: 22, Team2Score: 20},
	}, scoring)
	if err != nil {
		t.Fatal(err)
	}
	if winner != team1 || loser != team2 {
		t.Fatalf("winner/loser = %v/%v, want %v/%v", winner, loser, team1, team2)
	}
}

func TestValidateGameScoreRejectsInvalidDeuceScores(t *testing.T) {
	scoring := &repository.EventScoring{
		BestOf:        3,
		WinningPoints: 21,
		MaximumPoints: 30,
	}

	for _, scores := range [][2]int{{22, 0}, {30, 28}} {
		if err := validateGameScore(scores[0], scores[1], scoring); err == nil {
			t.Fatalf("validateGameScore(%d, %d) accepted invalid score", scores[0], scores[1])
		}
	}
}
