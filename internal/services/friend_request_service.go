package services

import (
	"converse/internal/models/friends"
	"converse/internal/repositories"
	"converse/internal/types"
)

type FriendRequestService struct {
	friendRequestRepo *repositories.FriendRequestRepository
}

func NewFriendRequestService() *FriendRequestService {
	return &FriendRequestService{
		friendRequestRepo: repositories.NewFriendRequestRepository(),
	}
}

func (s *FriendRequestService) CreateFriendRequest(request types.CreateFriendRequest, requesterID string) error {
	friendRequest := &friends.FriendRequest{
		RequesterID: requesterID,
		RecipientID: request.RecipientID,
	}

	return s.friendRequestRepo.Create(friendRequest)
}

func (s *FriendRequestService) DeclineFriendRequest(friendRequestID uint64) error {
	return s.friendRequestRepo.DeclineFriendRequest(friendRequestID)
}

func (s *FriendRequestService) GetUserFriendRequests(userID string) ([]*friends.FriendRequest, error) {
	return s.friendRequestRepo.GetUserFriendRequests(userID)
}

func (s *FriendRequestService) AcceptFriendRequest(friendRequestID uint64, userID string) error {
	return s.friendRequestRepo.AcceptFriendRequest(friendRequestID, userID)
}