package models

import "time"

type Room struct {
	RoomID      string     `json:"room_id" gorm:"type:char(36):column:room_id;primary_key;autoIncrement"`
	Name        string     `json:"name" gorm:"column:name;type:varchar(100);not null;index:idx_rooms_name"`
	Description string     `json:"description" gorm:"column:description;type:text"`
	IsPrivate   bool       `json:"is_private" gorm:"column:is_private;default:false"`
	CreatedBy   uint64     `json:"created_by" gorm:"column:created_by;not null;index:idx_rooms_created_by_user_id;constraint:OnDelete:CASCADE"`
	CreatedAt   time.Time  `json:"created_at" gorm:"column:created_at;autoCreateTime"`
	UpdatedAt   time.Time  `json:"updated_at" gorm:"column:updated_at;autoUpdateTime"`
	LastMessageAt *time.Time `json:"last_message_at" gorm:"column:last_message_at;type:timestamp;null"`
}

func (Room) TableName() string {
	return "rooms"
}