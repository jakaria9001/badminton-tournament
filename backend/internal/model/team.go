package model

type TeamPlayer struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type TeamResponse struct {
	ID        string     `json:"id"`
	TeamName  string     `json:"teamName"`
	Player1   TeamPlayer `json:"player1"`
	Player2   TeamPlayer `json:"player2"`
	Status    string     `json:"status"`
	CreatedAt string     `json:"createdAt"`
}

type TeamsResponse struct {
	Teams []TeamResponse `json:"teams"`
}
