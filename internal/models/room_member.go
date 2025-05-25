package models

import (
	"time"
)

// RoomMember represents a member of a room in the chat application.
type RoomMember struct {
	RoomMemberID       string    `gorm:"column:room_member_id;type:char(36);primaryKey"`
	RoomID             string    `gorm:"column:room_id;type:char(36);not null;index:idx_room_members_room_id;constraint:OnDelete:CASCADE"`
	UserID             string    `gorm:"column:user_id;type:char(36);not null;index:idx_room_members_user_id;constraint:OnDelete:CASCADE"`
	Role               string    `gorm:"column:role;type:enum('member', 'admin', 'owner');default:'member';not null"`
	JoinedAt           time.Time `gorm:"column:joined_at;autoCreateTime"`
	LastSeenMessageID  string    `gorm:"column:last_seen_message_id;type:char(36);index:idx_room_members_last_seen_message_id;constraint:OnDelete:SET NULL"`
}

func (RoomMember) TableName() string {
	return "room_members"
}