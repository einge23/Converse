package services

import (
	"converse/internal/models"
	"converse/internal/repositories"
)

// MessageService handles business logic for messages
type MessageService struct {
	messageRepo *repositories.MessageRepository
}

// NewMessageService creates a new message service
func NewMessageService() *MessageService {
	return &MessageService{
		messageRepo: repositories.NewMessageRepository(),
	}
}

// PaginatedMessages represents a paginated response of messages
type PaginatedMessages struct {
	Messages    []*models.Message `json:"messages"`
	CurrentPage int               `json:"current_page"`
	PageSize    int               `json:"page_size"`
	HasMore     bool              `json:"has_more"`
}

// GetMessagesByRoomID retrieves paginated messages for a specific room
func (s *MessageService) GetMessagesByRoomID(roomID string, page, pageSize int) (*PaginatedMessages, error) {
	// Calculate offset from page and pageSize
	offset := (page - 1) * pageSize
	
	// Get one more message than requested to determine if there are more pages
	limit := pageSize + 1
	
	// Fetch messages from repository
	messages, err := s.messageRepo.GetMessagesByRoomID(roomID, limit, offset)
	if err != nil {
		return nil, err
	}

	// Determine if there are more messages
	hasMore := false
	if len(messages) > pageSize {
		hasMore = true
		// Remove the extra message we fetched to check if there are more
		messages = messages[:pageSize]
	}

	return &PaginatedMessages{
		Messages:    messages,
		CurrentPage: page,
		PageSize:    pageSize,
		HasMore:     hasMore,
	}, nil
}

// GetMessagesByThreadID retrieves paginated messages for a specific thread
func (s *MessageService) GetMessagesByThreadID(threadID string, page, pageSize int) (*PaginatedMessages, error) {
	// Calculate offset from page and pageSize
	offset := (page - 1) * pageSize
	
	// Get one more message than requested to determine if there are more pages
	limit := pageSize + 1
	
	// Fetch messages from repository
	messages, err := s.messageRepo.GetMessagesByThreadID(threadID, limit, offset)
	if err != nil {
		return nil, err
	}

	// Determine if there are more messages
	hasMore := false
	if len(messages) > pageSize {
		hasMore = true
		// Remove the extra message we fetched to check if there are more
		messages = messages[:pageSize]
	}

	return &PaginatedMessages{
		Messages:    messages,
		CurrentPage: page,
		PageSize:    pageSize,
		HasMore:     hasMore,
	}, nil
}
