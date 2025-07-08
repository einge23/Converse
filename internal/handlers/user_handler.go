package handlers

import (
	"converse/internal/models"
	"converse/internal/services"
	"converse/pkg/errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
    userService *services.UserService
}

func NewUserHandler() *UserHandler {
    return &UserHandler{
        userService: services.NewUserService(),
    }
}

type UpdateStatusRequest struct {
    Status string `json:"status" binding:"required"`
}

func (h *UserHandler) UpdateStatus(c *gin.Context) {
    var req UpdateStatusRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, &errors.AppError{
            Code:    http.StatusBadRequest,
            Message: "Invalid request body",
            Details: err.Error(),
        })
        return
    }

    // Validate status value
    var status models.UserStatus
    switch req.Status {
    case "online":
        status = models.StatusOnline
    case "offline":
        status = models.StatusOffline
    case "away":
        status = models.StatusAway
    case "do_not_disturb":
        status = models.StatusDoNotDisturb
    default:
        c.JSON(http.StatusBadRequest, &errors.AppError{
            Code:    http.StatusBadRequest,
            Message: "Invalid status value",
            Details: "Status must be one of: online, offline, away, do_not_disturb",
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

    if err := h.userService.UpdateUserStatus(userIDStr, status); err != nil {
        c.JSON(http.StatusInternalServerError, &errors.AppError{
            Code:    http.StatusInternalServerError,
            Message: "Failed to update user status",
            Details: err.Error(),
        })
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "message": "Status updated successfully",
        "status":  req.Status,
    })
}