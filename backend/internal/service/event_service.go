package service

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"github.com/jakaria9001/badminton-tournament/backend/internal/model"
	"github.com/jakaria9001/badminton-tournament/backend/internal/repository"
)

type EventService struct {
	repository *repository.EventRepository
}

func NewEventService(
	repository *repository.EventRepository,
) *EventService {
	return &EventService{
		repository: repository,
	}
}

func (s *EventService) GetEventByID(
	ctx context.Context,
	eventID uuid.UUID,
) (*model.EventResponse, error) {

	return s.repository.GetEventByID(
		ctx,
		eventID,
	)
}

func (s *EventService) ListEvents(
	ctx context.Context,
) ([]model.EventResponse, error) {
	return s.repository.ListEvents(ctx)
}

func (s *EventService) ListAdminEvents(ctx context.Context) ([]model.EventResponse, error) {
	return s.repository.ListAdminEvents(ctx)
}

func (s *EventService) UpdateRegistrationStatus(
	ctx context.Context,
	eventID uuid.UUID,
	status string,
) error {
	if status != "REGISTRATION_OPEN" && status != "REGISTRATION_CLOSED" {
		return fmt.Errorf("invalid registration status")
	}
	return s.repository.UpdateRegistrationStatus(ctx, eventID, status)
}

func (s *EventService) CreateEvent(ctx context.Context, request model.EventAdminRequest) (uuid.UUID, error) {
	return s.repository.CreateEvent(ctx, request)
}

func (s *EventService) DeleteEvent(ctx context.Context, eventID uuid.UUID) error {
	return s.repository.DeleteEvent(ctx, eventID)
}

func (s *EventService) UpdateEvent(ctx context.Context, eventID uuid.UUID, request model.EventAdminRequest) error {
	return s.repository.UpdateEvent(ctx, eventID, request)
}
