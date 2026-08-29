package handler

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"

	"github.com/jakaria9001/badminton-tournament/backend/internal/model"
	"github.com/jakaria9001/badminton-tournament/backend/internal/service"
)

type ResultHandler struct {
	service *service.ResultService
}

func NewResultHandler(resultService *service.ResultService) *ResultHandler {
	return &ResultHandler{service: resultService}
}

func (h *ResultHandler) Submit(w http.ResponseWriter, r *http.Request) {
	matchID, err := uuid.Parse(chi.URLParam(r, "matchID"))
	if err != nil {
		http.Error(w, "invalid match ID", http.StatusBadRequest)
		return
	}

	var request model.SubmitMatchResultRequest
	if !decodeJSON(w, r, &request) {
		return
	}

	if err := h.service.Submit(r.Context(), matchID, request); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(map[string]string{"status": "completed"})
}
