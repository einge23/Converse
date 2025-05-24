package models

import (
	"time"
)

type Session struct {
	SessionID    string    `json:"session_id" gorm:"primaryKey"`
	UserID       string    `json:"user_id" gorm:"index"`
	Token        string    `json:"token"`
	ExpiresAt    time.Time `json:"expires_at"`
	CreatedAt    time.Time `json:"created_at"`
	IsValid      bool      `json:"is_valid" gorm:"default:true"`
	IPAddress    string    `json:"ip_address"`
	UserAgent    string    `json:"user_agent"`
	DeviceID     string    `json:"device_id"`
} 