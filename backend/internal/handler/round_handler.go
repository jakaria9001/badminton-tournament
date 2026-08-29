package handler

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"

	"github.com/jakaria9001/badminton-tournament/backend/internal/service"
)

type RoundHandler struct {
	roundService *service.RoundService
	drawService *service.DrawService
}

func NewRoundHandler(
	roundService *service.RoundService,
	drawService *service.DrawService,
) *RoundHandler {
	return &RoundHandler{
		roundService: roundService,
		drawService: drawService,
	}
}

type CreateRoundRequest struct {
	RoundNumber   int    `json:"roundNumber"`
	RoundName     string `json:"roundName"`
	PairingMethod string `json:"pairingMethod"`
}

func (h *RoundHandler) Create(
	w http.ResponseWriter,
	r *http.Request,
) {

	eventID, err := uuid.Parse(
		chi.URLParam(r, "eventID"),
	)

	if err != nil {
		http.Error(
			w,
			"invalid event ID",
			http.StatusBadRequest,
		)
		return
	}

	var req CreateRoundRequest

	if !decodeJSON(w, r, &req) {
		return
	}

	round, err := h.roundService.Create(
		r.Context(),
		eventID,
		req.RoundNumber,
		req.RoundName,
		req.PairingMethod,
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

	json.NewEncoder(w).Encode(round)
}

func (h *RoundHandler) GetByEvent(
	w http.ResponseWriter,
	r *http.Request,
) {

	eventID, err := uuid.Parse(
		chi.URLParam(r, "eventID"),
	)

	if err != nil {
		http.Error(
			w,
			"invalid event ID",
			http.StatusBadRequest,
		)
		return
	}

	rounds, err :=
		h.roundService.GetByEvent(
			r.Context(),
			eventID,
		)

	if err != nil {
		http.Error(
			w,
			"failed to get rounds",
			http.StatusInternalServerError,
		)
		return
	}

	w.Header().Set(
		"Content-Type",
		"application/json",
	)

	json.NewEncoder(w).Encode(rounds)
}

func (h *RoundHandler) Lock(
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

	err = h.roundService.Lock(
		r.Context(),
		roundID,
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

	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(map[string]string{
		"status": "locked",
	})
}

func (h *RoundHandler) GetAvailableTeams(
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

	teams, err :=
		h.roundService.GetAvailableTeams(
			r.Context(),
			roundID,
		)

	if err != nil {
		http.Error(
			w,
			"failed to get available teams",
			http.StatusInternalServerError,
		)
		return
	}

	w.Header().Set(
		"Content-Type",
		"application/json",
	)

	json.NewEncoder(w).Encode(teams)
}

func (h *RoundHandler) Generate(
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

	err = h.drawService.Generate(
		r.Context(),
		roundID,
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