package model

type PlayerInput struct {
	Name  string `json:"name"`
	Phone string `json:"phone"`
}

type RegistrationRequest struct {
	Player1  PlayerInput `json:"player1"`
	Player2  PlayerInput `json:"player2"`
	TeamName string      `json:"teamName"`
}

type RegistrationResponse struct {
	RegistrationID string `json:"registrationId"`
	Status         string `json:"status"`
}
