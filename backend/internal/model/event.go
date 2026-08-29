package model

type EventResponse struct {
	ID              string `json:"id"`
	Name            string `json:"name"`
	MaxTeams        *int   `json:"maxTeams"`
	RegisteredTeams int    `json:"registeredTeams"`
	Status          string `json:"status"`
}

type EventListResponse []EventResponse

type EventAdminRequest struct {
	Name                 string  `json:"name"`
	Description          string  `json:"description"`
	VenueName            string  `json:"venueName"`
	VenueAddress         string  `json:"venueAddress"`
	StartDate            string  `json:"startDate"`
	EndDate              string  `json:"endDate"`
	RegistrationDeadline *string `json:"registrationDeadline"`
	MaxTeams             *int    `json:"maxTeams"`
	Status               string  `json:"status"`
	AssignedAdminID      *string `json:"assignedAdminId,omitempty"`
}
