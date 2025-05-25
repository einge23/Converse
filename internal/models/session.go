package models

import (
	"time"
)

type Session struct {
	SessionID    string    `json:"session_id" gorm:"primaryKey;type:char(36)"`
	UserID       string    `json:"user_id" gorm:"index;type:char(36)"`
	Token        string    `json:"token" gorm:"index;type:varchar(512)"`
	ExpiresAt    time.Time `json:"expires_at" gorm:"index"`
	CreatedAt    time.Time `json:"created_at"`
	IsValid      bool      `json:"is_valid" gorm:"index;default:true"`
	IPAddress    string    `json:"ip_address" gorm:"type:varchar(45)"`
	UserAgent    string    `json:"user_agent" gorm:"type:text"`
	DeviceID     string    `json:"device_id" gorm:"index;type:char(36)"`
}

// Add composite indexes for common query patterns
func (Session) TableName() string {
	return "sessions"
} 