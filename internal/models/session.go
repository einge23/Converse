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
	UpdatedAt    time.Time `json:"updated_at"`
	LastActive   time.Time `json:"last_active"`
	IsValid      bool      `json:"is_valid" gorm:"default:true"`
	IPAddress    string    `json:"ip_address"`
	UserAgent    string    `json:"user_agent"`
	DeviceID     string    `json:"device_id"`
	LoginCount   int       `json:"login_count" gorm:"default:1"`
	LastIPChange time.Time `json:"last_ip_change"`
} 