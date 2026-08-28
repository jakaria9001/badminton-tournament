package model

type MatchResponse struct {
	ID            string       `json:"id"`
	EventID       string       `json:"eventId"`
	Round         string       `json:"round"`
	MatchNumber   int          `json:"matchNumber"`
	MatchType     string       `json:"matchType"`
	Team1         *TeamSummary `json:"team1"`
	Team2         *TeamSummary `json:"team2"`
	CourtNumber   *int         `json:"courtNumber"`
	ScheduledAt   *string      `json:"scheduledAt"`
	Status        string       `json:"status"`
	WinnerTeamID  *string      `json:"winnerTeamId"`
	LoserTeamID   *string      `json:"loserTeamId"`
	NextMatchID   *string      `json:"nextMatchId"`
	Games         []GameScore  `json:"games"`
}

type TeamSummary struct {
	ID       string `json:"id"`
	TeamName string `json:"teamName"`
	Player1  string `json:"player1"`
	Player2  string `json:"player2"`
}

type GameScore struct {
	GameNumber int `json:"gameNumber"`
	Team1Score int `json:"team1Score"`
	Team2Score int `json:"team2Score"`
}