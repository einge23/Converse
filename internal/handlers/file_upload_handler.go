package handlers

import (
	"net/http"

	"converse/internal/services"
	"converse/pkg/errors"

	"github.com/gin-gonic/gin"
)

type UploadHandler struct{}

func NewUploadHandler() *UploadHandler {
	return &UploadHandler{}
}

func (h *UploadHandler) UploadFile(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, &errors.AppError{
			Code:    http.StatusUnauthorized,
			Message: "User ID not found in token",
		})
		return
	}

	userIDstr, ok := userID.(string)
	if !ok {
		c.JSON(http.StatusInternalServerError, &errors.AppError{
			Code:    http.StatusInternalServerError,
			Message: "Invalid user ID format",
		})
		return
	}

	file, header, err := c.Request.FormFile("file")
	if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "No file provided"})
        return
    }
    defer file.Close()

	fileType := c.DefaultPostForm("type", "general")

	result, err := services.Services.FileUpload.UploadFile(c.Request.Context(), file, header, userIDstr, fileType)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Upload failed"})
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "url": result.URL,
        "key": result.Key,
        "filename": header.Filename,
    })
}

func (h *UploadHandler) UploadAvatar(c *gin.Context) {
    userID := c.GetString("userID")
    if userID == "" {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
        return
    }

    file, header, err := c.Request.FormFile("avatar")
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "No avatar file provided"})
        return
    }
    defer file.Close()

    result, err := services.Services.FileUpload.UploadFile(c.Request.Context(), file, header, userID, "avatars")
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Avatar upload failed"})
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "avatar_url": result.URL,
        "key": result.Key,
    })
}
	