package service

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"

	"github.com/jakaria9001/badminton-tournament/backend/internal/model"
	"github.com/jakaria9001/badminton-tournament/backend/internal/repository"
)

type AuthService struct {
	userRepository *repository.UserRepository
	jwtSecret      string
}

func NewAuthService(
	userRepository *repository.UserRepository,
	jwtSecret string,
) *AuthService {
	return &AuthService{
		userRepository: userRepository,
		jwtSecret:      jwtSecret,
	}
}

func (s *AuthService) Login(
	ctx context.Context,
	email string,
	password string,
) (string, error) {

	user, err := s.userRepository.GetByEmail(
		ctx,
		email,
	)

	if err != nil {
		return "", fmt.Errorf(
			"invalid email or password",
		)
	}

	err = bcrypt.CompareHashAndPassword(
		[]byte(user.PasswordHash),
		[]byte(password),
	)

	if err != nil {
		return "", fmt.Errorf(
			"invalid email or password",
		)
	}

	claims := jwt.MapClaims{
		"user_id": user.ID.String(),
		"role":    user.Role,
		"exp": time.Now().
			Add(24 * time.Hour).
			Unix(),
		"iat": time.Now().Unix(),
	}

	token := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		claims,
	)

	signedToken, err := token.SignedString(
		[]byte(s.jwtSecret),
	)

	if err != nil {
		return "", fmt.Errorf(
			"create token: %w",
			err,
		)
	}

	return signedToken, nil
}

func (s *AuthService) GetProfile(
	ctx context.Context,
	userID uuid.UUID,
) (model.AdminProfileResponse, error) {
	user, err := s.userRepository.GetByID(ctx, userID)
	if err != nil {
		return model.AdminProfileResponse{}, err
	}

	profile := model.AdminProfileResponse{
		ID:    user.ID.String(),
		Name:  user.Name,
		Email: user.Email,
		Role:  user.Role,
	}
	if user.EventID.Valid {
		eventID := user.EventID.UUID.String()
		profile.EventID = &eventID
	}

	return profile, nil
}

func (s *AuthService) CreateAdmin(
	ctx context.Context,
	request model.CreateAdminRequest,
) error {
	if strings.TrimSpace(request.Name) == "" || strings.TrimSpace(request.Email) == "" || strings.TrimSpace(request.Password) == "" {
		return fmt.Errorf("name, email, and password are required")
	}
	if request.Role != model.RoleAdmin && request.Role != model.RoleSuperAdmin {
		return fmt.Errorf("role must be ADMIN or SUPER_ADMIN")
	}
	if request.Role == model.RoleAdmin {
		if strings.TrimSpace(request.EventID) == "" {
			return fmt.Errorf("event assignment is required for admin users")
		}
		if _, err := uuid.Parse(request.EventID); err != nil {
			return fmt.Errorf("invalid event ID")
		}
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("hash password: %w", err)
	}

	var eventID *uuid.UUID
	if request.Role == model.RoleAdmin {
		parsedEventID, err := uuid.Parse(request.EventID)
		if err != nil {
			return fmt.Errorf("invalid event ID")
		}
		eventID = &parsedEventID
	}

	return s.userRepository.CreateAdmin(ctx, request.Name, strings.ToLower(request.Email), string(passwordHash), request.Role, eventID)
}

func (s *AuthService) ListAdmins(
	ctx context.Context,
) ([]model.AdminUserResponse, error) {
	users, err := s.userRepository.ListAdmins(ctx)
	if err != nil {
		return nil, err
	}

	admins := make([]model.AdminUserResponse, 0, len(users))
	for _, user := range users {
		response := model.AdminUserResponse{
			ID:    user.ID.String(),
			Name:  user.Name,
			Email: user.Email,
			Role:  user.Role,
		}
		if user.EventID.Valid {
			eventID := user.EventID.UUID.String()
			response.EventID = &eventID
		}
		admins = append(admins, response)
	}

	return admins, nil
}
