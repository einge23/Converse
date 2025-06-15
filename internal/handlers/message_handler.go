package handlers

import (
	"net/http"
	"strconv"

	"converse/internal/services"

	"github.com/gin-gonic/gin"
)

// MessageHandler handles HTTP requests related to messages
type MessageHandler struct {
	messageService *services.MessageService
}

// NewMessageHandler creates a new message handler
func NewMessageHandler() *MessageHandler {
	return &MessageHandler{
		messageService: services.NewMessageService(),
	}
}

// GetMessagesByRoomID handles the request to get paginated messages for a room
func (h *MessageHandler) GetMessagesByRoomID(c *gin.Context) {
	roomID := c.Param("room_id")
	if roomID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Room ID is required"})
		return
	}

	// Parse pagination parameters with defaults
	page, pageSize := h.getPaginationParams(c)

	// Get messages from service
	paginatedMessages, err := h.messageService.GetMessagesByRoomID(roomID, page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve messages"})
		return
	}

	c.JSON(http.StatusOK, paginatedMessages)
}

// GetMessagesByThreadID handles the request to get paginated messages for a thread
func (h *MessageHandler) GetMessagesByThreadID(c *gin.Context) {
	threadID := c.Param("thread_id")
	if threadID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Thread ID is required"})
		return
	}

	// Parse pagination parameters with defaults
	page, pageSize := h.getPaginationParams(c)

	// Get messages from service
	paginatedMessages, err := h.messageService.GetMessagesByThreadID(threadID, page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve messages"})
		return
	}

	c.JSON(http.StatusOK, paginatedMessages)
}

// getPaginationParams extracts and validates pagination parameters from the request
func (h *MessageHandler) getPaginationParams(c *gin.Context) (int, int) {
	// Default values
	defaultPage := 1
	defaultPageSize := 20
	maxPageSize := 100

	// Get page parameter
	pageStr := c.DefaultQuery("page", strconv.Itoa(defaultPage))
	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = defaultPage
	}

	// Get page_size parameter
	pageSizeStr := c.DefaultQuery("page_size", strconv.Itoa(defaultPageSize))
	pageSize, err := strconv.Atoi(pageSizeStr)
	if err != nil || pageSize < 1 {
		pageSize = defaultPageSize
	}

	// Limit the maximum page size
	if pageSize > maxPageSize {
		pageSize = maxPageSize
	}

	return page, pageSize
}
