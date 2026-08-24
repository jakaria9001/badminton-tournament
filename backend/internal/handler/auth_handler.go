package handler

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"

	"github.com/jakaria9001/badminton-tournament/backend/internal/middleware"
	"github.com/jakaria9001/badminton-tournament/backend/internal/model"
	"github.com/jakaria9001/badminton-tournament/backend/internal/service"
)

type AuthHandler struct {
	service *service.AuthService
}

func NewAuthHandler(
	service *service.AuthService,
) *AuthHandler {
	return &AuthHandler{
		service: service,
	}
}

func (h *AuthHandler) Login(
	w http.ResponseWriter,
	r *http.Request,
) {

	var req model.LoginRequest

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

	token, err := h.service.Login(
		r.Context(),
		req.Email,
		req.Password,
	)

	if err != nil {
		http.Error(
			w,
			"invalid email or password",
			http.StatusUnauthorized,
		)

		return
	}

	response := model.LoginResponse{
		Token: token,
	}

	w.Header().Set(
		"Content-Type",
		"application/json",
	)

	json.NewEncoder(w).Encode(response)
}

func (h *AuthHandler) Me(
	w http.ResponseWriter,
	r *http.Request,
) {
	userID, ok := r.Context().
		Value(middleware.UserIDKey).(uuid.UUID)

	if !ok {
		http.Error(
			w,
			"invalid authenticated user",
			http.StatusUnauthorized,
		)
		return
	}

	profile, err := h.service.GetProfile(
		r.Context(),
		userID,
	)
	if err != nil {
		http.Error(
			w,
			"failed to get admin profile",
			http.StatusInternalServerError,
		)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(profile)
}
