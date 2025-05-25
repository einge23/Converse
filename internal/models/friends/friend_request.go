package friends

import "time"

type FriendRequest struct {
    ID          uint64    `json:"id" gorm:"primaryKey;autoIncrement"`
    RequesterID string    `json:"requester_id" gorm:"type:char(36);not null;index:idx_friend_requests_requester;uniqueIndex:unique_request;constraint:OnDelete:CASCADE"`
    RecipientID string    `json:"recipient_id" gorm:"type:char(36);not null;index:idx_friend_requests_recipient;uniqueIndex:unique_request;constraint:OnDelete:CASCADE"`
    Status      string    `json:"status" gorm:"type:enum('pending', 'accepted', 'declined');default:'pending'"`
    CreatedAt   time.Time `json:"created_at" gorm:"autoCreateTime"`
    UpdatedAt   time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

func (FriendRequest) TableName() string {
    return "friend_requests"
}