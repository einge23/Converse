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

	// Get client information
	ipAddress := c.ClientIP()
	userAgent := c.GetHeader("User-Agent")
	deviceID := c.GetHeader("X-Device-ID")

	response, err := h.authService.Register(req, ipAddress, userAgent, deviceID)
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

	c.JSON(http.StatusCreated, response)
}

func(h *AuthHandler) Login(c *gin.Context) {
	var req types.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, &errors.AppError{
			Code:    http.StatusBadRequest,
			Message: "Invalid request body",
			Details: err.Error(),
		})
		return
	}

	// Get client information
	ipAddress := c.ClientIP()
	userAgent := c.GetHeader("User-Agent")
	deviceID := c.GetHeader("X-Device-ID")

	response, err := h.authService.Login(req, ipAddress, userAgent, deviceID)
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

	c.JSON(http.StatusOK, response)
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

func (h *AuthHandler) ValidateSession(c *gin.Context) {
	sessionID := c.GetHeader("X-Session-ID")
	if sessionID == "" {
		c.JSON(http.StatusBadRequest, &errors.AppError{
			Code:    http.StatusBadRequest,
			Message: "Session ID is required",
		})
		return
	}

	// Get client information
	ipAddress := c.ClientIP()
	userAgent := c.GetHeader("User-Agent")
	deviceID := c.GetHeader("X-Device-ID")

	user, err := h.authService.ValidateSession(sessionID, ipAddress, userAgent, deviceID)
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

	c.JSON(http.StatusOK, gin.H{
		"user": user,
	})
}

func (h *AuthHandler) Logout(c *gin.Context) {
	sessionID := c.GetHeader("X-Session-ID")
	if sessionID == "" {
		c.JSON(http.StatusBadRequest, &errors.AppError{
			Code:    http.StatusBadRequest,
			Message: "Session ID is required",
		})
		return
	}

	userID, _ := c.Get("user_id")
	if err := h.authService.Logout(sessionID, userID.(string)); err != nil {
		switch appErr := err.(type) {
		case *errors.AppError:
			c.JSON(appErr.Code, appErr)
		default:
			c.JSON(http.StatusInternalServerError, &errors.AppError{
				Code:    http.StatusInternalServerError,
				Message: "Failed to logout",
				Details: err.Error(),
			})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Successfully logged out",
	})
}

func (h *AuthHandler) RefreshToken(c *gin.Context) {
	var req types.RefreshTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, &errors.AppError{
            Code:    http.StatusBadRequest,
            Message: "Invalid request body",
            Details: err.Error(),
        })
        return
    }

	ipAddress := c.ClientIP()
	userAgent := c.GetHeader("User-Agent")
	deviceID := c.GetHeader("X-Device-ID")

    response, err := h.authService.RefreshToken(req.RefreshToken, ipAddress, userAgent, deviceID)
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
	c.JSON(http.StatusOK, response)
}

func (h *AuthHandler) Me(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, &errors.AppError{
			Code:    http.StatusUnauthorized,
			Message: "User ID not found in token",
		})
		return
	}

	user, err := h.authService.Me(userID.(string))
	if err != nil {
		switch appErr := err.(type) {
		case *errors.AppError:
			c.JSON(appErr.Code, appErr)
		default:
			c.JSON(http.StatusInternalServerError, &errors.AppError{
				Code:    http.StatusInternalServerError,
				Message: "Failed to retrieve user information",
				Details: err.Error(),
			})
		}
		return
	}

	c.JSON(http.StatusOK, user)
}
