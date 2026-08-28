package model

type CreateMatchRequest struct {
	Team1ID string `json:"team1Id"`
	Team2ID string `json:"team2Id"`
}