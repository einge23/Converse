package services

import (
	"converse/internal/models"
	"converse/internal/repositories"
	"converse/internal/types"
	"converse/internal/utils"
	"converse/pkg/errors"

	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	userRepo *repositories.UserRepository
}

func NewAuthService() *AuthService {
	return &AuthService{
		userRepo: repositories.NewUserRepository(),
	}
}

func (s *AuthService) Register(req types.RegisterRequest) (*types.AuthResponse, error) {
existingUser, _ := s.userRepo.FindByUsername(req.Username)
    if existingUser != nil {
        return nil, errors.NewConflictError("Username already exists")
    }

    // Check if email exists
    existingUser, _ = s.userRepo.FindByEmail(req.Email)
    if existingUser != nil {
        return nil, errors.NewConflictError("Email already exists")
    }

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
    if err != nil {
        return nil, err
    }

	user := &models.User{
        Username:     req.Username,
        Email:        req.Email,
        PasswordHash: string(hashedPassword),
        DisplayName:  req.DisplayName,
        Status:       models.StatusOffline,
    }

	if err := s.userRepo.Create(user); err != nil {
        return nil, err
    }

	token, expiresAt, err := utils.GenerateJWT(user)
    if err != nil {
        return nil, err
    }

	return &types.AuthResponse{
        UserID:      user.UserID,
        Username:    user.Username,
        Email:       user.Email,
        DisplayName: user.DisplayName,
        Token:       token,
        ExpiresAt:   expiresAt,
    }, nil
}
