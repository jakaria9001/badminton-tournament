package model

type MatchResponse struct {
	ID          string `json:"id"`
	EventID     string `json:"eventId"`
	RoundID     string `json:"roundId"`

	Round       string `json:"round"`
	MatchNumber int    `json:"matchNumber"`
	MatchType   string `json:"matchType"`

	Team1 *TeamSummary `json:"team1"`
	Team2 *TeamSummary `json:"team2"`

	CourtNumber *int    `json:"courtNumber"`
	ScheduledAt *string `json:"scheduledAt"`

	Status string `json:"status"`

	WinnerTeamID *string `json:"winnerTeamId"`
	LoserTeamID  *string `json:"loserTeamId"`

	WinnerNextMatchID *string `json:"winnerNextMatchId"`
	LoserNextMatchID  *string `json:"loserNextMatchId"`

	Team1SourceMatchID *string `json:"team1SourceMatchId"`
	Team1SourceType    *string `json:"team1SourceType"`

	Team2SourceMatchID *string `json:"team2SourceMatchId"`
	Team2SourceType    *string `json:"team2SourceType"`

	Games []GameScore `json:"games"`
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