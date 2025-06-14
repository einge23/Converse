package repositories

import (
	"converse/internal/db"
	"converse/internal/models"
	"errors"

	"gorm.io/gorm"
)

type MessageRepository struct {
	db *gorm.DB
}

func NewMessageRepository() *MessageRepository {
	return &MessageRepository{
		db: db.GetDB(),
	}
}

func (m *MessageRepository) StoreMessage(message *models.Message) error {
	if message == nil {
		return errors.New("message cannot be nil")
	}

	// Validate that message has either RoomID or ThreadID
	if message.RoomID == nil && message.ThreadID == nil {
		return errors.New("message must have either room_id or thread_id")
	}

	// Validate required fields
	if message.Content == "" {
		return errors.New("message content cannot be empty")
	}

	if message.ContentType == "" {
		message.ContentType = "text" // Default to text if not specified
	}

	// Create the message in database
	if err := m.db.Create(message).Error; err != nil {
		return err
	}

	return nil
}

// GetMessagesByRoomID retrieves messages for a specific room with pagination
func (m *MessageRepository) GetMessagesByRoomID(roomID string, limit int, offset int) ([]*models.Message, error) {
	var messages []*models.Message

	err := m.db.Where("room_id = ? AND deleted_at IS NULL", roomID).
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&messages).Error

	return messages, err
}

// GetMessagesByThreadID retrieves messages for a specific thread with pagination
func (m *MessageRepository) GetMessagesByThreadID(threadID string, limit int, offset int) ([]*models.Message, error) {
	var messages []*models.Message

	err := m.db.Where("thread_id = ? AND deleted_at IS NULL", threadID).
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&messages).Error

	return messages, err
}

// DB returns the database connection for use by other components
func (m *MessageRepository) DB() *gorm.DB {
	return m.db
}