package service

import (
	"context"
	"fmt"
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

	return model.AdminProfileResponse{
		ID:    user.ID.String(),
		Name:  user.Name,
		Email: user.Email,
		Role:  user.Role,
	}, nil
}
