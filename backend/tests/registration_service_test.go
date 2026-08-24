package tests

import (
	"strings"
	"testing"

	"github.com/jakaria9001/badminton-tournament/backend/internal/model"
	"github.com/jakaria9001/badminton-tournament/backend/internal/service"
)

func validRegistrationRequest() model.RegistrationRequest {
	return model.RegistrationRequest{
		Player1: model.PlayerInput{Name: "Player One", Phone: "9876543210"},
		Player2: model.PlayerInput{Name: "Player Two", Phone: "8765432109"},
	}
}

func TestValidateRegistrationRequest(t *testing.T) {
	tests := []struct {
		name    string
		request model.RegistrationRequest
		wantErr string
	}{
		{
			name:    "valid request",
			request: validRegistrationRequest(),
		},
		{
			name: "optional player two phone",
			request: func() model.RegistrationRequest {
				request := validRegistrationRequest()
				request.Player2.Phone = ""
				return request
			}(),
		},
		{
			name: "missing player one name",
			request: func() model.RegistrationRequest {
				request := validRegistrationRequest()
				request.Player1.Name = " "
				return request
			}(),
			wantErr: "player 1 name is required",
		},
		{
			name: "missing player one phone",
			request: func() model.RegistrationRequest {
				request := validRegistrationRequest()
				request.Player1.Phone = ""
				return request
			}(),
			wantErr: "player 1 phone is required",
		},
		{
			name: "missing player two name",
			request: func() model.RegistrationRequest {
				request := validRegistrationRequest()
				request.Player2.Name = ""
				return request
			}(),
			wantErr: "player 2 name is required",
		},
		{
			name: "invalid player one phone",
			request: func() model.RegistrationRequest {
				request := validRegistrationRequest()
				request.Player1.Phone = "1234567890"
				return request
			}(),
			wantErr: "player 1 phone must be a valid 10-digit Indian mobile number",
		},
		{
			name: "invalid player two phone",
			request: func() model.RegistrationRequest {
				request := validRegistrationRequest()
				request.Player2.Phone = "987654321"
				return request
			}(),
			wantErr: "player 2 phone must be a valid 10-digit Indian mobile number",
		},
		{
			name: "duplicate phones with whitespace",
			request: func() model.RegistrationRequest {
				request := validRegistrationRequest()
				request.Player1.Phone = " 9876543210 "
				request.Player2.Phone = "9876543210"
				return request
			}(),
			wantErr: "cannot have the same phone number",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := service.ValidateRegistrationRequest(test.request)
			if test.wantErr == "" {
				if err != nil {
					t.Fatalf("expected no error, got %v", err)
				}
				return
			}
			if err == nil || !strings.Contains(err.Error(), test.wantErr) {
				t.Fatalf("expected error containing %q, got %v", test.wantErr, err)
			}
		})
	}
}
