package models

import "time"

type DirectMessageThread struct {
	DirectMessageThreadID string    `json:"thread_id" gorm:"column:thread_id;type:char(36);primaryKey"`
	User1ID               string    `json:"user1_id" gorm:"column:user1_id;type:char(36);not null;index:idx_dm_threads_user1;constraint:OnDelete:CASCADE"`
	User2ID               string    `json:"user2_id" gorm:"column:user2_id;type:char(36);not null;index:idx_dm_threads_user2;constraint:OnDelete:CASCADE"`
	CreatedAt             time.Time `json:"created_at" gorm:"column:created_at;autoCreateTime"`
	UpdatedAt             time.Time `json:"updated_at" gorm:"column:updated_at;autoUpdateTime"`
	LastMessageAt		 *time.Time `json:"last_message_at" gorm:"column:last_message_at;type:timestamp;null"`
	User1LastSeenMessageID string    `json:"user1_last_seen_message_id" gorm:"column:user1_last_seen_message_id;type:char(36);constraint:OnDelete:SET NULL"`
	User2LastSeenMessageID string    `json:"user2_last_seen_message_id" gorm:"column:user2_last_seen_message_id;type:char(36);constraint:OnDelete:SET NULL"`
}

func (DirectMessageThread) TableName() string {
	return "direct_message_threads"
}