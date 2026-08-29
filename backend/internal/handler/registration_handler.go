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

type RegistrationHandler struct {
	service *service.RegistrationService
}

func NewRegistrationHandler(
	service *service.RegistrationService,
) *RegistrationHandler {
	return &RegistrationHandler{
		service: service,
	}
}

func (h *RegistrationHandler) Create(
	w http.ResponseWriter,
	r *http.Request,
) {
	eventIDString := chi.URLParam(r, "eventID")

	eventID, err := uuid.Parse(eventIDString)
	if err != nil {
		http.Error(
			w,
			"invalid event ID",
			http.StatusBadRequest,
		)
		return
	}

	var req model.RegistrationRequest

	if !decodeJSON(w, r, &req) {
		return
	}

	registrationID, err := h.service.RegisterTeam(
		r.Context(),
		eventID,
		req,
	)

	if err != nil {
		status := RegistrationErrorStatus(err)
		http.Error(
			w,
			err.Error(),
			status,
		)
		return
	}

	response := model.RegistrationResponse{
		RegistrationID: registrationID.String(),
		Status:         "PENDING",
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	json.NewEncoder(w).Encode(response)
}

func RegistrationErrorStatus(err error) int {
	switch {
	case errors.Is(err, model.ErrEventNotFound):
		return http.StatusNotFound
	case errors.Is(err, model.ErrRegistrationClosed),
		errors.Is(err, model.ErrRegistrationDeadlinePassed),
		errors.Is(err, model.ErrEventFull):
		return http.StatusConflict
	case errors.Is(err, model.ErrParticipantAlreadyRegistered):
		return http.StatusBadRequest
	default:
		return http.StatusInternalServerError
	}
}

func (h *RegistrationHandler) GetTeams(
	w http.ResponseWriter,
	r *http.Request,
) {
	eventIDString := chi.URLParam(r, "eventID")

	eventID, err := uuid.Parse(eventIDString)
	if err != nil {
		http.Error(
			w,
			"invalid event ID",
			http.StatusBadRequest,
		)
		return
	}

	teams, err := h.service.GetTeamsByEvent(
		r.Context(),
		eventID,
	)

	if err != nil {
		http.Error(
			w,
			"failed to get teams",
			http.StatusInternalServerError,
		)
		return
	}

	response := model.TeamsResponse{
		Teams: teams,
	}

	w.Header().Set(
		"Content-Type",
		"application/json",
	)

	json.NewEncoder(w).Encode(response)
}

func (h *RegistrationHandler) GetRegistrations(
	w http.ResponseWriter,
	r *http.Request,
) {
	eventIDString := chi.URLParam(
		r,
		"eventID",
	)

	eventID, err := uuid.Parse(
		eventIDString,
	)

	if err != nil {
		http.Error(
			w,
			"invalid event ID",
			http.StatusBadRequest,
		)
		return
	}

	registrations, err :=
		h.service.GetRegistrationsByEvent(
			r.Context(),
			eventID,
		)

	if err != nil {
		http.Error(
			w,
			"failed to get registrations",
			http.StatusInternalServerError,
		)
		return
	}

	w.Header().Set(
		"Content-Type",
		"application/json",
	)

	json.NewEncoder(w).Encode(
		registrations,
	)
}

func (h *RegistrationHandler) UpdateStatus(
	w http.ResponseWriter,
	r *http.Request,
) {
	idString := chi.URLParam(
		r,
		"registrationID",
	)

	registrationID, err := uuid.Parse(
		idString,
	)

	if err != nil {
		http.Error(
			w,
			"invalid registration ID",
			http.StatusBadRequest,
		)
		return
	}

	var req model.UpdateRegistrationStatusRequest

	if !decodeJSON(w, r, &req) {
		return
	}

	err = h.service.UpdateStatus(
		r.Context(),
		registrationID,
		req.Status,
	)

	if err != nil {
		http.Error(
			w,
			err.Error(),
			http.StatusBadRequest,
		)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *RegistrationHandler) WithdrawRegistration(
	w http.ResponseWriter,
	r *http.Request,
) {
	idString := chi.URLParam(
		r,
		"registrationID",
	)

	registrationID, err := uuid.Parse(
		idString,
	)

	if err != nil {
		http.Error(
			w,
			"invalid registration ID",
			http.StatusBadRequest,
		)
		return
	}

	err = h.service.WithdrawRegistration(
		r.Context(),
		registrationID,
	)

	if err != nil {
		http.Error(
			w,
			err.Error(),
			http.StatusBadRequest,
		)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
