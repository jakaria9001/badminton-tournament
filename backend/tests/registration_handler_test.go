package tests

import (
	"errors"
	"net/http"
	"testing"

	"github.com/jakaria9001/badminton-tournament/backend/internal/handler"
	"github.com/jakaria9001/badminton-tournament/backend/internal/model"
)

func TestRegistrationErrorStatus(t *testing.T) {
	tests := []struct {
		name string
		err  error
		want int
	}{
		{name: "event not found", err: model.ErrEventNotFound, want: http.StatusNotFound},
		{name: "registration closed", err: model.ErrRegistrationClosed, want: http.StatusConflict},
		{name: "deadline passed", err: model.ErrRegistrationDeadlinePassed, want: http.StatusConflict},
		{name: "event full", err: model.ErrEventFull, want: http.StatusConflict},
		{name: "invalid participant", err: model.ErrParticipantAlreadyRegistered, want: http.StatusBadRequest},
		{name: "unexpected error", err: errors.New("database unavailable"), want: http.StatusInternalServerError},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got := handler.RegistrationErrorStatus(test.err); got != test.want {
				t.Fatalf("expected status %d, got %d", test.want, got)
			}
		})
	}
}
