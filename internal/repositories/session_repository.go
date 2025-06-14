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
    err := r.db.Select("session_id, user_id, token, refresh_token, expires_at, refresh_expires_at, created_at, is_valid, ip_address, user_agent, device_id").
        Where("session_id = ? AND is_valid = ?", sessionID, true).
        First(&session).Error
    if err != nil {
        return nil, err
    }
    return &session, nil
}

func (r *SessionRepository) FindByToken(token string) (*models.Session, error) {
    var session models.Session
    err := r.db.Select("session_id, user_id, token, refresh_token, expires_at, refresh_expires_at, created_at, is_valid, ip_address, user_agent, device_id").
        Where("token = ? AND is_valid = ?", token, true).
        First(&session).Error
    if err != nil {
        return nil, err
    }
    return &session, nil
}

func (r *SessionRepository) FindByRefreshToken(refreshToken string) (*models.Session, error) {
	var session models.Session
    err := r.db.Where("refresh_token = ? AND is_valid = ? AND refresh_expires_at > ?", 
        refreshToken, true, time.Now()).First(&session).Error
    if err != nil {
        return nil, err
    }
    return &session, nil
}

func (r *SessionRepository) Invalidate(sessionID string) error {
	return r.db.Exec("UPDATE sessions SET is_valid = false WHERE session_id = ?", sessionID).Error
}

func (r *SessionRepository) DeleteExpired() error {
	return r.db.Where("expires_at < ?", time.Now()).Delete(&models.Session{}).Error
}

func (r *SessionRepository) Update(session *models.Session) error {
	return r.db.Save(session).Error
}

func (r *SessionRepository) UpdateTokens(sessionID, newToken, newRefreshToken string, 
    tokenExpiry, refreshExpiry time.Time) error {
    return r.db.Model(&models.Session{}).
        Where("session_id = ?", sessionID).
        Updates(map[string]interface{}{
            "token":              newToken,
            "refresh_token":      newRefreshToken,
            "expires_at":         tokenExpiry,
            "refresh_expires_at": refreshExpiry,
        }).Error
}

func (r *SessionRepository) FindByUserID(userID string) ([]*models.Session, error) {
    var sessions []*models.Session
    err := r.db.Select("session_id, user_id, token, refresh_token, expires_at, refresh_expires_at, created_at, is_valid, ip_address, user_agent, device_id").
        Where("user_id = ? AND is_valid = ?", userID, true).
        Find(&sessions).Error
    if err != nil {
        return nil, err
    }
    return sessions, nil
}

func (r *SessionRepository) InvalidateAllForUser(userID string) error {
	return r.db.Exec("UPDATE sessions SET is_valid = false WHERE user_id = ?", userID).Error
} 