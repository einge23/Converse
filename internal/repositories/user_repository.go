package repositories

import (
	"converse/internal/db"
	"converse/internal/models"

	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository() *UserRepository {
	return &UserRepository{
		db: db.GetDB(),
	}
}

func (r *UserRepository) Create(user *models.User) error {
	return r.db.Create(user).Error
}

func (r *UserRepository) FindByUsername(username string) (*models.User, error) {
	var user models.User
	err := r.db.Select("user_id, username, email, password_hash, display_name, avatar_url, status, last_active_at, created_at, updated_at").
		Where("username = ?", username).
		First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) FindByEmail(email string) (*models.User, error) {
	var user models.User
	err := r.db.Select("user_id, username, email, password_hash, display_name, avatar_url, status, last_active_at, created_at, updated_at").
		Where("email = ?", email).
		First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) FindByUsernameOrEmail(username, email string) (*models.User, error) {
	var user models.User
	err := r.db.Select("user_id, username, email, password_hash, display_name, avatar_url, status, last_active_at, created_at, updated_at").
		Where("username = ? OR email = ?", username, email).
		First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) FindByID(userID string) (*models.User, error) {
	var user models.User
	err := r.db.Select("user_id, username, email, password_hash, display_name, avatar_url, status, last_active_at, created_at, updated_at").
		Where("user_id = ?", userID).
		First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

