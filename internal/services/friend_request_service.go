package services

import (
	"converse/internal/models/friends"
	"converse/internal/repositories"
)

type FriendRequestService struct {
	friendRequestRepo *repositories.FriendRequestRepository
}

func NewFriendRequestService() *FriendRequestService {
	return &FriendRequestService{
		friendRequestRepo: repositories.NewFriendRequestRepository(),
	}
}

func (s *FriendRequestService) CreateFriendRequest(requesterID, recipientID string) error {
	friendRequest := &friends.FriendRequest{
		RequesterID: requesterID,
		RecipientID: recipientID,
	}

	return s.friendRequestRepo.Create(friendRequest)
}

//do rest tomorrow