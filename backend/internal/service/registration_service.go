package service

import (
	"context"
	"fmt"
	"regexp"
	"strings"

	"github.com/google/uuid"

	"github.com/jakaria9001/badminton-tournament/backend/internal/model"
	"github.com/jakaria9001/badminton-tournament/backend/internal/repository"
)

type RegistrationService struct {
	repository      *repository.RegistrationRepository
	eventRepository *repository.EventRepository
}

var indianMobileNumber = regexp.MustCompile(`^[6-9][0-9]{9}$`)

func NewRegistrationService(
	repository *repository.RegistrationRepository,
	eventRepository *repository.EventRepository,
) *RegistrationService {

	return &RegistrationService{
		repository:      repository,
		eventRepository: eventRepository,
	}
}

func (s *RegistrationService) RegisterTeam(
	ctx context.Context,
	eventID uuid.UUID,
	req model.RegistrationRequest,
) (uuid.UUID, error) {
	if err := ValidateRegistrationRequest(req); err != nil {
		return uuid.Nil, err
	}

	req.Player1.Name = strings.TrimSpace(req.Player1.Name)
	req.Player1.Phone = strings.TrimSpace(req.Player1.Phone)
	req.Player2.Name = strings.TrimSpace(req.Player2.Name)
	req.Player2.Phone = strings.TrimSpace(req.Player2.Phone)
	req.TeamName = strings.TrimSpace(req.TeamName)

	return s.repository.CreateRegistration(
		ctx,
		eventID,
		req,
	)
}

func ValidateRegistrationRequest(req model.RegistrationRequest) error {
	if strings.TrimSpace(req.Player1.Name) == "" {
		return fmt.Errorf("player 1 name is required")
	}
	if strings.TrimSpace(req.Player1.Phone) == "" {
		return fmt.Errorf("player 1 phone is required")
	}
	if !indianMobileNumber.MatchString(strings.TrimSpace(req.Player1.Phone)) {
		return fmt.Errorf("player 1 phone must be a valid 10-digit Indian mobile number")
	}
	if strings.TrimSpace(req.Player2.Name) == "" {
		return fmt.Errorf("player 2 name is required")
	}
	player1Phone := strings.TrimSpace(req.Player1.Phone)
	player2Phone := strings.TrimSpace(req.Player2.Phone)
	if player2Phone != "" && !indianMobileNumber.MatchString(player2Phone) {
		return fmt.Errorf("player 2 phone must be a valid 10-digit Indian mobile number")
	}
	if player2Phone != "" && player1Phone == player2Phone {
		return fmt.Errorf("player 1 and player 2 cannot have the same phone number")
	}
	return nil
}

func (s *RegistrationService) GetTeamsByEvent(
	ctx context.Context,
	eventID uuid.UUID,
) ([]model.TeamResponse, error) {

	return s.repository.GetTeamsByEvent(
		ctx,
		eventID,
	)
}

func (s *RegistrationService) GetRegistrationsByEvent(
	ctx context.Context,
	eventID uuid.UUID,
) ([]model.AdminRegistration, error) {

	return s.repository.GetRegistrationsByEvent(
		ctx,
		eventID,
	)
}

func (s *RegistrationService) UpdateStatus(
	ctx context.Context,
	registrationID uuid.UUID,
	status string,
) error {

	if status != "CONFIRMED" && status != "REJECTED" && status != "WITHDRAWN" {
		return fmt.Errorf(
			"invalid registration status",
		)
	}

	return s.repository.TransitionStatus(
		ctx,
		registrationID,
		status,
	)
}

func (s *RegistrationService) WithdrawRegistration(
	ctx context.Context,
	registrationID uuid.UUID,
) error {
	return s.repository.TransitionStatus(
		ctx,
		registrationID,
		"WITHDRAWN",
	)
}

func ValidRegistrationTransition(current, next string) bool {
	switch current {
	case "PENDING":
		return next == "CONFIRMED" || next == "REJECTED" || next == "WITHDRAWN"
	case "CONFIRMED":
		return next == "WITHDRAWN"
	default:
		return false
	}
}
