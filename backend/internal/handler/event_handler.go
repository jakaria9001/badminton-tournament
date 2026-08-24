package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"

	"github.com/jakaria9001/badminton-tournament/backend/internal/model"
	"github.com/jakaria9001/badminton-tournament/backend/internal/service"
)

type EventHandler struct {
	service *service.EventService
}

func NewEventHandler(
	service *service.EventService,
) *EventHandler {
	return &EventHandler{
		service: service,
	}
}

func (h *EventHandler) GetByID(
	w http.ResponseWriter,
	r *http.Request,
) {
	eventIDString := chi.URLParam(
		r,
		"eventID",
	)

	eventID, err := uuid.Parse(eventIDString)

	if err != nil {
		if errors.Is(err, model.ErrEventNotFound) {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		http.Error(
			w,
			"invalid event ID",
			http.StatusBadRequest,
		)
		return
	}

	event, err := h.service.GetEventByID(
		r.Context(),
		eventID,
	)

	if err != nil {
		http.Error(
			w,
			"failed to get event",
			http.StatusInternalServerError,
		)
		return
	}

	w.Header().Set(
		"Content-Type",
		"application/json",
	)

	json.NewEncoder(w).Encode(event)
}
