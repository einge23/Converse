package repositories

import (
	"converse/internal/db"
	"converse/internal/models/friends"
	"converse/pkg/errors"
	stderrors "errors"
	"net/http"
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

func (r *FriendRequestRepository) Create(username string, requesterID string) error {
	// First, find the recipient user by username
	userRepo := NewUserRepository()
	user, err := userRepo.FindByUsername(username)
	if err != nil {
		if stderrors.Is(err, gorm.ErrRecordNotFound) {
			return &errors.AppError{
				Code:    http.StatusBadRequest,
				Message: "User not found",
				Details: "No user exists with the provided username",
			}
		}
		return err
	}

	// Check if user is trying to send a friend request to themselves
	if user.UserID == requesterID {
		return &errors.AppError{
			Code:    http.StatusBadRequest,
			Message: "Cannot send friend request to yourself",
		}
	}

	// Check if there's already a pending or accepted friend request between these users
	var existingRequest friends.FriendRequest
	err = r.db.Where(
		"((requester_id = ? AND recipient_id = ?) OR (requester_id = ? AND recipient_id = ?)) AND status IN ('pending', 'accepted')",
		requesterID, user.UserID, user.UserID, requesterID,
	).First(&existingRequest).Error

	if err == nil {
		// Friend request already exists
		if existingRequest.Status == "accepted" {
			return &errors.AppError{
				Code:    http.StatusBadRequest,
				Message: "Users are already friends",
				Details: "A friendship already exists between these users",
			}
		}
		return &errors.AppError{
			Code:    http.StatusBadRequest,
			Message: "Friend request already pending",
			Details: "A friend request is already pending between these users",
		}
	} else if !stderrors.Is(err, gorm.ErrRecordNotFound) {
		// Some other database error occurred
		return err
	}

	// Create the friend request
	friendRequest := &friends.FriendRequest{
		RequesterID: requesterID,
		RecipientID: user.UserID,
		Status:      "pending",
	}

	return r.db.Create(friendRequest).Error
}

func (r *FriendRequestRepository) DeclineFriendRequest(friendRequestID uint64) error {
	return r.db.Model(&friends.FriendRequest{}).
		Where("friend_request_id = ?", friendRequestID).
		Update("status", "declined").Error
}

func (r *FriendRequestRepository) GetUserFriendRequests(userID string) ([]*friends.FriendRequestWithUser, error) {
	var requests []*friends.FriendRequestWithUser
	err := r.db.Table("friend_requests fr").
		Select(`fr.friend_request_id, fr.requester_id, fr.recipient_id, fr.status, fr.created_at, fr.updated_at,
			u.user_id as user_user_id, u.username as user_username, u.email as user_email, 
			u.display_name as user_display_name, u.avatar_url as user_avatar_url, 
			u.status as user_status, u.last_active_at as user_last_active_at, 
			u.created_at as user_created_at, u.updated_at as user_updated_at`).
		Joins("JOIN users u ON fr.requester_id = u.user_id").
		Where("fr.recipient_id = ? AND fr.status = ?", userID, "pending").
		Order("fr.created_at DESC").
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
		err := tx.Where("friend_request_id = ? AND recipient_id = ? AND status = ?", friendRequestID, userID, "pending").
			First(&friendRequest).Error
		if err != nil {
			if stderrors.Is(err, gorm.ErrRecordNotFound) {
				return &errors.AppError{
					Code:    http.StatusBadRequest,
					Message: "Friend request not found",
					Details: "Friend request not found or you are not authorized to accept it",
				}
			}
			return err
		}

		// Update friend request status to accepted
		err = tx.Model(&friendRequest).
			Updates(map[string]any{
				"status":     "accepted",
				"updated_at": time.Now(),
			}).Error
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
			return &errors.AppError{
				Code:    http.StatusInternalServerError,
				Message: "Failed to create friendship",
				Details: err.Error(),
			}
		}

		// Create direct message thread between the two users
		dmRepo := NewDirectMessageRepository()
		_, err = dmRepo.CreateDirectMessageThread(friendRequest.RequesterID, friendRequest.RecipientID)
		if err != nil {
			return &errors.AppError{
				Code:    http.StatusInternalServerError,
				Message: "Failed to create direct message thread",
				Details: err.Error(),
			}
		}

		return nil
	})
}