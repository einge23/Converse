package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserStatus string

const (
	StatusOnline       UserStatus = "online"
	StatusOffline      UserStatus = "offline"
	StatusAway         UserStatus = "away"
	StatusDoNotDisturb UserStatus = "do_not_disturb"
)

type User struct {
	UserID       string     `gorm:"column:user_id;type:char(36);primary_key"`
	Username     string     `gorm:"column:username;type:varchar(50);uniqueIndex;not null"`
	Email        string     `gorm:"column:email;type:varchar(255);uniqueIndex;not null"`
	PasswordHash string     `gorm:"column:password_hash;type:varchar(255);not null"`
	DisplayName  string     `gorm:"column:display_name;type:varchar(100)"`
	AvatarURL    string     `gorm:"column:avatar_url;type:text"`
	Status       UserStatus `gorm:"column:status;type:enum('online','offline','away','do_not_disturb');default:offline"`
	LastActiveAt *time.Time `gorm:"column:last_active_at;type:timestamp;null"`
	CreatedAt    time.Time  `gorm:"column:created_at;type:timestamp;not null;default:CURRENT_TIMESTAMP"`
	UpdatedAt    time.Time  `gorm:"column:updated_at;type:timestamp;not null;default:CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP"`
	DeletedAt    *time.Time `gorm:"column:deleted_at;type:timestamp;default:NULL"`
}

func (User) TableName() string {
    return "users"
}

func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
	if u.UserID == "" {
		u.UserID = uuid.New().String()
	}
	return nil
}

func (u *User) BeforeUpdate(tx *gorm.DB) error {
    u.UpdatedAt = time.Now()
    return nil
}

type PublicUser struct {
	UserID       string     `json:"user_id" gorm:"primaryKey"`
	Username     string     `json:"username"`
	Email        string     `json:"email"`
	DisplayName  string     `json:"display_name"`
	AvatarURL    *string    `json:"avatar_url"`
	Status       string     `json:"status"`
	LastActiveAt *time.Time `json:"last_active_at"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
	DeletedAt    *time.Time `json:"deleted_at,omitempty"`
	DMThreadID   *string    `json:"dm_thread_id,omitempty"`
}

func (u *User) ToPublicUser() *PublicUser {
	return &PublicUser{
		UserID:       u.UserID,
		Username:     u.Username,
		Email:        u.Email,
		DisplayName:  u.DisplayName,
		AvatarURL:    &u.AvatarURL,
		Status:       string(u.Status),
		LastActiveAt: u.LastActiveAt,
		CreatedAt:    u.CreatedAt,
		UpdatedAt:    u.UpdatedAt,
		DeletedAt:    u.DeletedAt,
	}
}