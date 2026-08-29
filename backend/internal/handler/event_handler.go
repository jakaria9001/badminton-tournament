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

func (h *EventHandler) List(
	w http.ResponseWriter,
	r *http.Request,
) {
	events, err := h.service.ListEvents(r.Context())
	if err != nil {
		http.Error(w, "failed to list events", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(events)
}

func (h *EventHandler) ListAdmin(
	w http.ResponseWriter,
	r *http.Request,
) {
	events, err := h.service.ListAdminEvents(r.Context())
	if err != nil {
		http.Error(w, "failed to list events", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(events)
}

func (h *EventHandler) UpdateRegistrationStatus(
	w http.ResponseWriter,
	r *http.Request,
) {
	eventID, err := uuid.Parse(chi.URLParam(r, "eventID"))
	if err != nil {
		http.Error(w, "invalid event ID", http.StatusBadRequest)
		return
	}

	var request model.UpdateRegistrationStatusRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if err := h.service.UpdateRegistrationStatus(r.Context(), eventID, request.Status); err != nil {
		if errors.Is(err, model.ErrEventNotFound) {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]string{"status": request.Status})
}

func (h *EventHandler) Create(
	w http.ResponseWriter,
	r *http.Request,
) {
	var request model.EventAdminRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}
	eventID, err := h.service.CreateEvent(r.Context(), request)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(map[string]string{"id": eventID.String()})
}

func (h *EventHandler) Delete(
	w http.ResponseWriter,
	r *http.Request,
) {
	eventID, err := uuid.Parse(chi.URLParam(r, "eventID"))
	if err != nil {
		http.Error(w, "invalid event ID", http.StatusBadRequest)
		return
	}
	if err := h.service.DeleteEvent(r.Context(), eventID); err != nil {
		if errors.Is(err, model.ErrEventNotFound) {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *EventHandler) Update(
	w http.ResponseWriter,
	r *http.Request,
) {
	eventID, err := uuid.Parse(chi.URLParam(r, "eventID"))
	if err != nil {
		http.Error(w, "invalid event ID", http.StatusBadRequest)
		return
	}
	var request model.EventAdminRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}
	if err := h.service.UpdateEvent(r.Context(), eventID, request); err != nil {
		if errors.Is(err, model.ErrEventNotFound) {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusNoContent)
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
