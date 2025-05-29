package services

import (
	"converse/internal/models"
	"converse/internal/repositories"
	"converse/internal/types"
	"converse/internal/utils"
	"converse/pkg/errors"
	"log"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	userRepo    *repositories.UserRepository
	sessionRepo *repositories.SessionRepository
}

func NewAuthService() *AuthService {
	return &AuthService{
		userRepo:    repositories.NewUserRepository(),
		sessionRepo: repositories.NewSessionRepository(),
	}
}

func (s *AuthService) Register(req types.RegisterRequest, ipAddress, userAgent, deviceID string) (*types.AuthResponse, error) {
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

	// Create a new session
	session := &models.Session{
		SessionID:  uuid.New().String(),
		UserID:     user.UserID,
		Token:      token,
		ExpiresAt:  expiresAt,
		CreatedAt:  time.Now(),
		IsValid:    true,
		IPAddress:  ipAddress,
		UserAgent:  userAgent,
		DeviceID:   deviceID,
	}

	if err := s.sessionRepo.Create(session); err != nil {
		return nil, err
	}

	return &types.AuthResponse{
		Token:     token,
		ExpiresAt: expiresAt,
		SessionID: session.SessionID,
	}, nil
}

func (s *AuthService) Login(req types.LoginRequest, ipAddress, userAgent, deviceID string) (*types.AuthResponse, error) {
	user, err := s.userRepo.FindByUsernameOrEmail(req.Username, req.Email)
	if err != nil {
		return nil, errors.NewUnauthorizedError("Invalid credentials", "User not found")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password))
	if err != nil {
		return nil, errors.NewUnauthorizedError("Invalid credentials", "Invalid password")
	}

	token, expiresAt, err := utils.GenerateJWT(user)
	if err != nil {
		return nil, err
	}

	// Invalidate any existing sessions for this user
	if err := s.sessionRepo.InvalidateAllForUser(user.UserID); err != nil {
		return nil, err
	}

	// Create a new session
	session := &models.Session{
		SessionID:  uuid.New().String(),
		UserID:     user.UserID,
		Token:      token,
		ExpiresAt:  expiresAt,
		CreatedAt:  time.Now(),
		IsValid:    true,
		IPAddress:  ipAddress,
		UserAgent:  userAgent,
		DeviceID:   deviceID,
	}

	if err := s.sessionRepo.Create(session); err != nil {
		return nil, err
	}

	log.Printf("Session created successfully - Session ID: %s, Device ID: %s", session.SessionID, deviceID)

	return &types.AuthResponse{
		Token:     token,
		ExpiresAt: expiresAt,
		SessionID: session.SessionID,
	}, nil
}

func (s *AuthService) ValidateSession(sessionID, ipAddress, userAgent, deviceID string) (*models.User, error) {
	session, err := s.sessionRepo.FindByID(sessionID)
	if err != nil {
		return nil, errors.NewUnauthorizedError("Invalid session", "Session not found or expired")
	}

	if session.ExpiresAt.Before(time.Now()) {
		s.sessionRepo.Invalidate(sessionID)
		return nil, errors.NewUnauthorizedError("Session expired", "Please login again")
	}

	// Verify device ID if provided
	if deviceID != "" && session.DeviceID != deviceID {
		s.sessionRepo.Invalidate(sessionID)
		return nil, errors.NewUnauthorizedError("Invalid device", "Please login again")
	}

	user, err := s.userRepo.FindByID(session.UserID)
	if err != nil {
		return nil, errors.NewUnauthorizedError("Invalid session", "User not found")
	}

	return user, nil
}

func (s *AuthService) Logout(sessionID, userID string) error {
	session, err := s.sessionRepo.FindByID(sessionID)
	if err != nil {
		return err
	}

	// Only allow users to logout their own sessions
	if session.UserID != userID {
		return errors.NewForbiddenError("Cannot logout another user's session")
	}

	return s.sessionRepo.Invalidate(sessionID)
}

func (s *AuthService) Me(userID string) (*models.PublicUser, error) {
	user, err := s.userRepo.FindPublicUserByID(userID)
	if err != nil {
		return nil, errors.NewUnauthorizedError("User not found", "Please login again")
	}
	return user, nil
}
