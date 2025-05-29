package services

import (
	"converse/internal/models"
	"converse/internal/repositories"
)

type FriendshipService struct {
	friendshipRepo *repositories.FriendRepository
}

func NewFriendshipService() *FriendshipService {
	return &FriendshipService{
		friendshipRepo: repositories.NewFriendRepository(),
	}
}

func (s *FriendshipService) GetFriends(userID string) ([]*models.PublicUser, error) {
	return s.friendshipRepo.GetFriends(userID)
}