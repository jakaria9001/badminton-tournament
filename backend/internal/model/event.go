package model

type EventResponse struct {
	ID              string `json:"id"`
	Name            string `json:"name"`
	MaxTeams        *int   `json:"maxTeams"`
	RegisteredTeams int    `json:"registeredTeams"`
	Status          string `json:"status"`
}
