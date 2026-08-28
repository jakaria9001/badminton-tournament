package handler

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"

	"github.com/jakaria9001/badminton-tournament/backend/internal/model"
	"github.com/jakaria9001/badminton-tournament/backend/internal/service"
)

type MatchHandler struct {
	service *service.MatchService
}

func NewMatchHandler(
	service *service.MatchService,
) *MatchHandler {
	return &MatchHandler{
		service: service,
	}
}

func (h *MatchHandler) GetByEvent(
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

	matches, err :=
		h.service.GetMatchesByEvent(
			r.Context(),
			eventID,
		)

	if err != nil {
		http.Error(
			w,
			"failed to get matches",
			http.StatusInternalServerError,
		)
		return
	}

	w.Header().Set(
		"Content-Type",
		"application/json",
	)

	json.NewEncoder(w).Encode(
		matches,
	)
}

func (h *MatchHandler) Create(
	w http.ResponseWriter,
	r *http.Request,
) {

	roundID, err := uuid.Parse(
		chi.URLParam(r, "roundID"),
	)

	if err != nil {
		http.Error(
			w,
			"invalid round ID",
			http.StatusBadRequest,
		)
		return
	}

	var req model.CreateMatchRequest

	if err := json.NewDecoder(
		r.Body,
	).Decode(&req); err != nil {

		http.Error(
			w,
			"invalid request body",
			http.StatusBadRequest,
		)

		return
	}

	team1ID, err := uuid.Parse(req.Team1ID)

	if err != nil {
		http.Error(
			w,
			"invalid team1 ID",
			http.StatusBadRequest,
		)
		return
	}

	team2ID, err := uuid.Parse(req.Team2ID)

	if err != nil {
		http.Error(
			w,
			"invalid team2 ID",
			http.StatusBadRequest,
		)
		return
	}

	match, err := h.service.CreateMatch(
		r.Context(),
		roundID,
		team1ID,
		team2ID,
	)

	if err != nil {
		http.Error(
			w,
			err.Error(),
			http.StatusBadRequest,
		)
		return
	}

	w.Header().Set(
		"Content-Type",
		"application/json",
	)

	w.WriteHeader(http.StatusCreated)

	json.NewEncoder(w).Encode(map[string]string{
		"id": match.String(),
	})
}