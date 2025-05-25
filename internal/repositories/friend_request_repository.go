package repositories

import (
	"converse/internal/db"
	"converse/internal/models/friends"
	"errors"
	"time"

	"gorm.io/gorm"
)

type FriendRequestRepository struct {
	db *gorm.DB
}

func NewFriendRequestRepository() *FriendRequestRepository {
	return &FriendRequestRepository{
		db: db.GetDB(),
	}
}

func (r *FriendRequestRepository) Create(friendRequest *friends.FriendRequest) error {
	return r.db.Create(friendRequest).Error
}

func (r *FriendRequestRepository) UpdateStatus(friendRequestID uint64, status string) error {
	return r.db.Model(&friends.FriendRequest{}).
		Where("friend_request_id = ?", friendRequestID).
		Update("status", status).Error
}

func (r *FriendRequestRepository) GetUserFriendRequests(userID string) ([]*friends.FriendRequest, error) {
	var requests []*friends.FriendRequest
	err := r.db.Where("recipient_id = ? OR requester_id = ?", userID, userID).
		Order("created_at DESC").
		Find(&requests).Error
	if err != nil {
		return nil, err
	}
	return requests, nil
}

func (r *FriendRequestRepository) AcceptFriendRequest(friendRequestID uint64, userID string) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		// First, verify that the user is the recipient of this friend request
		var friendRequest friends.FriendRequest
		err := tx.Where("id = ? AND recipient_id = ? AND status = ?", friendRequestID, userID, "pending").
			First(&friendRequest).Error
		if err != nil {
			if err == gorm.ErrRecordNotFound {
				return errors.New("friend request not found or you are not authorized to accept it")
			}
			return err
		}

		// Update friend request status to accepted
		err = tx.Model(&friendRequest).
			Updates(map[string]interface{}{
				"status":     "accepted",
				"updated_at": time.Now(),
			}).Error
		if err != nil {
			return err
		}

		// Retrieve the updated friend request to get requester and recipient IDs
		err = tx.Where("friend_request_id = ?", friendRequestID).First(&friendRequest).Error
		if err != nil {
			return err
		}

		// Ensure user1_id < user2_id for the CHECK constraint
		var user1ID, user2ID string
		if friendRequest.RequesterID < friendRequest.RecipientID {
			user1ID = friendRequest.RequesterID
			user2ID = friendRequest.RecipientID
		} else {
			user1ID = friendRequest.RecipientID
			user2ID = friendRequest.RequesterID
		}

		// Create new friendship record
		friendship := &friends.Friendship{
			User1ID:   user1ID,
			User2ID:   user2ID,
			CreatedAt: time.Now(),
		}

		err = tx.Create(friendship).Error
		if err != nil {
			return err
		}

		return nil
	})
}