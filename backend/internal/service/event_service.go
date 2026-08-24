package service

import (
	"context"

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
