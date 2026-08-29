package model

import "time"

type TournamentRound struct {
	ID            string  `json:"id"`
	EventID       string  `json:"eventId"`
	RoundNumber   int     `json:"roundNumber"`
	RoundName     string  `json:"roundName"`
	PairingMethod string  `json:"pairingMethod"`
	Status        string  `json:"status"`
	LockedAt      *time.Time `json:"lockedAt"`
	CompletedAt   *time.Time `json:"completedAt"`
}