package models

import (
	"time"
)

type Session struct {
    SessionID        string    `json:"session_id" gorm:"primaryKey;type:char(36)"`
    UserID           string    `json:"user_id" gorm:"not null;type:char(36);index"`
    Token            string    `json:"token" gorm:"not null;type:varchar(512)"`
    RefreshToken     string    `json:"refresh_token" gorm:"not null;type:varchar(512)"`
    ExpiresAt        time.Time `json:"expires_at" gorm:"not null"`
    RefreshExpiresAt time.Time `json:"refresh_expires_at" gorm:"not null"`
    CreatedAt        time.Time `json:"created_at" gorm:"autoCreateTime"`
    IsValid          bool      `json:"is_valid" gorm:"default:true"`
    IPAddress        string    `json:"ip_address" gorm:"type:varchar(45)"`
    UserAgent        string    `json:"user_agent" gorm:"type:text"`
    DeviceID         string    `json:"device_id" gorm:"type:varchar(255)"`
}

// Add composite indexes for common query patterns
func (Session) TableName() string {
	return "sessions"
} 