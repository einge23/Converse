package handlers

import (
	"converse/internal/services"
	"converse/pkg/errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

type FriendshipHandler struct {
	friendshipService *services.FriendshipService
}

func NewFriendshipHandler() *FriendshipHandler {
	return &FriendshipHandler{
		friendshipService: services.NewFriendshipService(),
	}
}

func (h *FriendshipHandler) GetFriends(c *gin.Context) {
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

	friends, err := h.friendshipService.GetFriends(userIDStr)
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

	c.JSON(http.StatusOK, friends)
}
