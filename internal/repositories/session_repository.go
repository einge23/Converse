package repositories

import (
	"converse/internal/db"
	"converse/internal/models"
	"time"

	"gorm.io/gorm"
)

type SessionRepository struct {
	db *gorm.DB
}

func NewSessionRepository() *SessionRepository {
	return &SessionRepository{
		db: db.GetDB(),
	}
}

func (r *SessionRepository) Create(session *models.Session) error {
	return r.db.Create(session).Error
}

func (r *SessionRepository) FindByID(sessionID string) (*models.Session, error) {
	var session models.Session
	err := r.db.Where("session_id = ? AND is_valid = ?", sessionID, true).First(&session).Error
	if err != nil {
		return nil, err
	}
	return &session, nil
}

func (r *SessionRepository) FindByToken(token string) (*models.Session, error) {
	var session models.Session
	err := r.db.Where("token = ? AND is_valid = ?", token, true).First(&session).Error
	if err != nil {
		return nil, err
	}
	return &session, nil
}

func (r *SessionRepository) UpdateLastActive(sessionID string) error {
	return r.db.Model(&models.Session{}).
		Where("session_id = ?", sessionID).
		Update("last_active", time.Now()).Error
}

func (r *SessionRepository) Invalidate(sessionID string) error {
	return r.db.Model(&models.Session{}).
		Where("session_id = ?", sessionID).
		Update("is_valid", false).Error
}

func (r *SessionRepository) DeleteExpired() error {
	return r.db.Where("expires_at < ?", time.Now()).Delete(&models.Session{}).Error
}

func (r *SessionRepository) Update(session *models.Session) error {
	return r.db.Save(session).Error
}

func (r *SessionRepository) FindByUserID(userID string) ([]*models.Session, error) {
	var sessions []*models.Session
	err := r.db.Where("user_id = ? AND is_valid = ?", userID, true).Find(&sessions).Error
	if err != nil {
		return nil, err
	}
	return sessions, nil
}

func (r *SessionRepository) InvalidateAllForUser(userID string) error {
	return r.db.Model(&models.Session{}).
		Where("user_id = ?", userID).
		Update("is_valid", false).Error
} 