package model

type SubmitMatchResultRequest struct {
	WinnerTeamID string       `json:"winnerTeamId"`
	Games        []GameResult `json:"games"`
}

type GameResult struct {
	GameNumber int `json:"gameNumber"`
	Team1Score int `json:"team1Score"`
	Team2Score int `json:"team2Score"`
}