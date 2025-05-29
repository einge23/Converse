package friends

import (
	"converse/internal/models"
	"time"
)

type FriendRequest struct {
    FriendRequestID          uint64    `json:"friend_request_id" gorm:"primaryKey;autoIncrement"`
    RequesterID string    `json:"requester_id" gorm:"type:char(36);not null;index:idx_friend_requests_requester;uniqueIndex:unique_request;constraint:OnDelete:CASCADE"`
    RecipientID string    `json:"recipient_id" gorm:"type:char(36);not null;index:idx_friend_requests_recipient;uniqueIndex:unique_request;constraint:OnDelete:CASCADE"`
    Status      string    `json:"status" gorm:"type:enum('pending', 'accepted', 'declined');default:'pending'"`
    CreatedAt   time.Time `json:"created_at" gorm:"autoCreateTime"`
    UpdatedAt   time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

// FriendRequestWithUser combines friend request data with requester's public user info
type FriendRequestWithUser struct {
    FriendRequestID uint64               `json:"friend_request_id"`
    RequesterID     string               `json:"requester_id"`
    RecipientID     string               `json:"recipient_id"`
    Status          string               `json:"status"`
    CreatedAt       time.Time            `json:"created_at"`
    UpdatedAt       time.Time            `json:"updated_at"`
    Requester       models.PublicUser    `json:"requester" gorm:"embedded;embeddedPrefix:user_"`
}

func (FriendRequest) TableName() string {
    return "friend_requests"
}