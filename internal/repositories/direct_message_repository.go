package repositories

import (
	"converse/internal/db"
	"converse/internal/models"
	"errors"

	"gorm.io/gorm"
)

type DirectMessageRepository struct {
	db *gorm.DB
}

func NewDirectMessageRepository() *DirectMessageRepository {
	return &DirectMessageRepository{
		db: db.GetDB(),
	}
}

func (r *DirectMessageRepository) CreateDirectMessageThread(user1ID, user2ID string) (*models.DirectMessageThread, error) {
	// Check if a thread already exists between these two users
	var existingThread models.DirectMessageThread
	err := r.db.Where(
		"(user1_id = ? AND user2_id = ?) OR (user1_id = ? AND user2_id = ?)", 
		user1ID, user2ID, user2ID, user1ID,
	).First(&existingThread).Error

	if err == nil {
		// Thread already exists
		return &existingThread, nil
	}

	if !errors.Is(err, gorm.ErrRecordNotFound) {
		// Database error
		return nil, err
	}

	// Create new thread
	newThread := models.DirectMessageThread{
		User1ID: user1ID,
		User2ID: user2ID,
	}

	if err := r.db.Create(&newThread).Error; err != nil {
		return nil, err
	}

	return &newThread, nil
}

//tomorrow work on messages and retrieval for thread. then work on
//websocket message routing