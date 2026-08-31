package handler

import (
	"encoding/json"
	"net/http"
	"os"
	"time"

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

	if !decodeJSON(w, r, &req) {
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

	secure := secureSessionCookie()
	sameSite := http.SameSiteLaxMode
	if secure {
		sameSite = http.SameSiteNoneMode
	}

	cookie := &http.Cookie{
		Name:     "shuttlehub_session",
		Value:    token,
		Path:     "/",
		HttpOnly: true,
		Secure:   secure,
		SameSite: sameSite,
		MaxAge:   int((24 * time.Hour).Seconds()),
	}
	if secure {
		cookie.Domain = getCookieDomain()
	}

	http.SetCookie(w, cookie)
	w.WriteHeader(http.StatusNoContent)
}

func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	secure := secureSessionCookie()
	sameSite := http.SameSiteLaxMode
	if secure {
		sameSite = http.SameSiteNoneMode
	}

	cookie := &http.Cookie{
		Name:     "shuttlehub_session",
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		Secure:   secure,
		SameSite: sameSite,
		MaxAge:   -1,
		Expires:  time.Unix(0, 0),
	}
	if secure {
		cookie.Domain = getCookieDomain()
	}

	http.SetCookie(w, cookie)
	w.WriteHeader(http.StatusNoContent)
}

func secureSessionCookie() bool {
	return os.Getenv("COOKIE_SECURE") == "true"
}

func getCookieDomain() string {
	return os.Getenv("COOKIE_DOMAIN")
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

func (h *AuthHandler) CreateAdmin(
	w http.ResponseWriter,
	r *http.Request,
) {
	var req model.CreateAdminRequest
	if !decodeJSON(w, r, &req) {
		return
	}

	if err := h.service.CreateAdmin(r.Context(), req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"status": "created"})
}

func (h *AuthHandler) ListAdmins(
	w http.ResponseWriter,
	r *http.Request,
) {
	admins, err := h.service.ListAdmins(r.Context())
	if err != nil {
		http.Error(w, "failed to list admins", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(admins)
}
