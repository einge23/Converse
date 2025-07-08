package services

import (
	"converse/internal/models"
	"converse/internal/repositories"
)

type UserService struct {
	userRepo *repositories.UserRepository
}

func NewUserService() *UserService {
	return &UserService{
		userRepo: repositories.NewUserRepository(),
	}
}

func (s *UserService) UpdateUserStatus(userID string, status models.UserStatus) error {
    return s.userRepo.UpdateStatus(userID, status)
}