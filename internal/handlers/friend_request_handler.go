package handlers

import (
	"converse/internal/services"
	"converse/internal/types"
	"converse/pkg/errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type FriendRequestHandler struct {
	friendRequestService *services.FriendRequestService
}

func NewFriendRequestHandler() *FriendRequestHandler {
	return &FriendRequestHandler{
		friendRequestService: services.NewFriendRequestService(),
	}
}

func (h *FriendRequestHandler) CreateFriendRequest(c *gin.Context) {
	var req types.CreateFriendRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, &errors.AppError{
			Code:    http.StatusBadRequest,
			Message: "Invalid request body",
			Details: err.Error(),
		})
		return
	}

	// Get requester ID from auth token claims
	requesterID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, &errors.AppError{
			Code:    http.StatusUnauthorized,
			Message: "User ID not found in token",
		})
		return
	}

	requesterIDStr, ok := requesterID.(string)
	if !ok {
		c.JSON(http.StatusInternalServerError, &errors.AppError{
			Code:    http.StatusInternalServerError,
			Message: "Invalid user ID format",
		})
		return
	}

	if err := h.friendRequestService.CreateFriendRequest(req, requesterIDStr); err != nil {
		switch appErr := err.(type) {
		case *errors.AppError:
			c.JSON(appErr.Code, appErr)
		default:
			c.JSON(http.StatusInternalServerError, &errors.AppError{
				Code:    http.StatusInternalServerError,
				Message: "Internal server error",
			})
		}
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Friend request sent successfully"})
}


func (h *FriendRequestHandler) DeclineFriendRequest(c *gin.Context) {
	friendRequestID := c.Param("friend_request_id")
	if friendRequestID == "" {
		c.JSON(http.StatusBadRequest, &errors.AppError{
			Code:    http.StatusBadRequest,
			Message: "Friend request ID is required",
		})
		return
	}

	// Convert friendRequestID to uint64
	var friendRequestIDUint64 uint64
	if _, err := fmt.Sscanf(friendRequestID, "%d", &friendRequestIDUint64); err != nil {
		c.JSON(http.StatusBadRequest, &errors.AppError{
			Code:    http.StatusBadRequest,
			Message: "Invalid friend request ID format",
			Details: err.Error(),
		})
		return
	}

	if err := h.friendRequestService.DeclineFriendRequest(friendRequestIDUint64); err != nil {
		switch appErr := err.(type) {
		case *errors.AppError:
			c.JSON(appErr.Code, appErr)
		default:
			c.JSON(http.StatusInternalServerError, &errors.AppError{
				Code:    http.StatusInternalServerError,
				Message: "Internal server error",
			})
		}
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Friend request declined successfully"})
}

func (h *FriendRequestHandler) GetUserFriendRequests(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, &errors.AppError{
			Code:    http.StatusUnauthorized,
			Message: "User ID not found in token",
		})
		return
	}

	userIDStr, ok := userID.(string)
	if !ok {
		c.JSON(http.StatusInternalServerError, &errors.AppError{
			Code:    http.StatusInternalServerError,
			Message: "Invalid user ID format",
		})
		return
	}

	friendRequests, err := h.friendRequestService.GetUserFriendRequests(userIDStr)
	if err != nil {
		switch appErr := err.(type) {
		case *errors.AppError:
			c.JSON(appErr.Code, appErr)
		default:
			c.JSON(http.StatusInternalServerError, &errors.AppError{
				Code:    http.StatusInternalServerError,
				Message: "Internal server error",
			})
		}
		return
	}

	c.JSON(http.StatusOK, friendRequests)
}

func (h *FriendRequestHandler) AcceptFriendRequest(c *gin.Context) {
	friendRequestID := c.Param("friend_request_id")
	if friendRequestID == "" {
		c.JSON(http.StatusBadRequest, &errors.AppError{
			Code:    http.StatusBadRequest,
			Message: "Friend request ID is required",
		})
		return
	}

	// Convert friendRequestID to uint64
	var friendRequestIDUint64 uint64
	if _, err := fmt.Sscanf(friendRequestID, "%d", &friendRequestIDUint64); err != nil {
		c.JSON(http.StatusBadRequest, &errors.AppError{
			Code:    http.StatusBadRequest,
			Message: "Invalid friend request ID format",
			Details: err.Error(),
		})
		return
	}

	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, &errors.AppError{
			Code:    http.StatusUnauthorized,
			Message: "User ID not found in token",
		})
		return
	}

	userIDStr, ok := userID.(string)
	if !ok {
		c.JSON(http.StatusInternalServerError, &errors.AppError{
			Code:    http.StatusInternalServerError,
			Message: "Invalid user ID format",
		})
		return
	}

	if err := h.friendRequestService.AcceptFriendRequest(friendRequestIDUint64, userIDStr); err != nil {
		switch appErr := err.(type) {
		case *errors.AppError:
			c.JSON(appErr.Code, appErr)
		default:
			c.JSON(http.StatusInternalServerError, &errors.AppError{
				Code:    http.StatusInternalServerError,
				Message: "Internal server error",
			})
		}
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Friend request accepted successfully"})
}