package model

import "errors"

type ErrorResponse struct {
	Error string `json:"error"`
}

var (
	ErrEventNotFound                = errors.New("event not found")
	ErrRegistrationClosed           = errors.New("registration is closed")
	ErrEventFull                    = errors.New("event is full")
	ErrRegistrationDeadlinePassed   = errors.New("registration deadline has passed")
	ErrParticipantAlreadyRegistered = errors.New("participant is already registered in this event")
)
