package model

type AdminRegistration struct {
	ID           string `json:"id"`
	TeamID       string `json:"teamId"`
	TeamName     string `json:"teamName"`
	Player1Name  string `json:"player1Name"`
	Player1Phone string `json:"player1Phone"`
	Player2Name  string `json:"player2Name"`
	Player2Phone string `json:"player2Phone"`
	Status       string `json:"status"`
	RegisteredAt string `json:"registeredAt"`
}
