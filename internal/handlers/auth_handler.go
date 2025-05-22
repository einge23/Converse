package handlers

import (
	"converse/internal/services"
	"converse/internal/types"
	"converse/pkg/errors"
	"net/http"
	"regexp"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	authService *services.AuthService
}

func NewAuthHandler() *AuthHandler {
    return &AuthHandler{
        authService: services.NewAuthService(),
    }
}


func(h *AuthHandler) Register(c *gin.Context) {
	var req types.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, &errors.AppError{
            Code:    http.StatusBadRequest,
            Message: "Invalid request body",
            Details: err.Error(),
        })
        return
    }

	if err := validatePassword(req.Password); err != nil {
        c.JSON(http.StatusBadRequest, err)
        return
    }

    if err := validateUsername(req.Username); err != nil {
        c.JSON(http.StatusBadRequest, err)
        return
    }

	response, err := h.authService.Register(req)
    if err != nil {
        switch err.(type) {
        case *errors.AppError:
            appErr := err.(*errors.AppError)
            c.JSON(appErr.Code, appErr)
        default:
            c.JSON(http.StatusInternalServerError, &errors.AppError{
                Code:    http.StatusInternalServerError,
                Message: "Internal server error",
            })
        }
        return
    }

    c.JSON(http.StatusCreated, response)
}

func validatePassword(password string) *errors.AppError {
    if len(password) < 8 {
        return errors.NewBadRequestError(
            "Password is too weak",
            "Password must be at least 8 characters long",
        )
    }

    hasNumber := regexp.MustCompile(`[0-9]`).MatchString(password)
    if !hasNumber {
        return errors.NewBadRequestError(
            "Password is too weak",
            "Password must contain at least one number",
        )
    }

    hasUpper := regexp.MustCompile(`[A-Z]`).MatchString(password)
    if !hasUpper {
        return errors.NewBadRequestError(
            "Password is too weak",
            "Password must contain at least one uppercase letter",
        )
    }

    return nil
}

func validateUsername(username string) *errors.AppError {
    validUsername := regexp.MustCompile(`^[a-zA-Z0-9_-]+$`).MatchString(username)
    if !validUsername {
        return errors.NewBadRequestError(
            "Invalid username format",
            "Username can only contain letters, numbers, underscores, and hyphens",
        )
    }

    return nil
}
