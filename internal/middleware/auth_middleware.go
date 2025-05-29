package middleware

import (
	"converse/internal/services"
	"converse/internal/utils"
	"converse/pkg/errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc {
	authService := services.NewAuthService()

	return func(c *gin.Context) {
		// Get and validate JWT token
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, &errors.AppError{
				Code:    http.StatusUnauthorized,
				Message: "Authorization header is required",
			})
			c.Abort()
			return
		}

		// Extract token from "Bearer <token>"
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, &errors.AppError{
				Code:    http.StatusUnauthorized,
				Message: "Invalid authorization header format",
			})
			c.Abort()
			return
		}

		tokenString := parts[1]
		claims, err := utils.ValidateToken(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, &errors.AppError{
				Code:    http.StatusUnauthorized,
				Message: "Invalid token",
			})
			c.Abort()
			return
		}

		// Get and validate session
		sessionID := c.GetHeader("X-Session-ID")
		if sessionID == "" {
			c.JSON(http.StatusUnauthorized, &errors.AppError{
				Code:    http.StatusUnauthorized,
				Message: "Session ID is required",
			})
			c.Abort()
			return
		}

		// Get client information
		ipAddress := c.ClientIP()
		userAgent := c.GetHeader("User-Agent")
		deviceID := c.GetHeader("X-Device-ID")

		user, err := authService.ValidateSession(sessionID, ipAddress, userAgent, deviceID)
		if err != nil {
			switch appErr := err.(type) {
			case *errors.AppError:
				c.JSON(appErr.Code, appErr)
			default:
				c.JSON(http.StatusUnauthorized, &errors.AppError{
					Code:    http.StatusUnauthorized,
					Message: "Invalid session",
				})
			}
			c.Abort()
			return
		}

		// Verify that the session's user ID matches the JWT claims
		if user.UserID != claims.UserID {
			c.JSON(http.StatusUnauthorized, &errors.AppError{
				Code:    http.StatusUnauthorized,
				Message: "Session user mismatch",
			})
			c.Abort()
			return
		}

		// Set user information in context
		c.Set("user_id", user.UserID)
		c.Set("email", user.Email)
		c.Set("user", user)
		c.Next()
	}
}

func OwnResourceMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, exists := c.Get("user_id")
		if !exists {
			c.JSON(http.StatusUnauthorized, &errors.AppError{
				Code:    http.StatusUnauthorized,
				Message: "User not authenticated",
			})
			c.Abort()
			return
		}

		resourceID := c.Param("user_id")
		if resourceID == "" {
			c.JSON(http.StatusBadRequest, &errors.AppError{
				Code:    http.StatusBadRequest,
				Message: "user_id parameter is required in the URL path",
				Details: "This endpoint requires a user_id parameter in the URL (e.g., /users/:user_id/resource)",
			})
			c.Abort()
			return
		}

		if userID != resourceID {
			c.JSON(http.StatusForbidden, &errors.AppError{
				Code:    http.StatusForbidden,
				Message: "Access denied - you can only access your own resources",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

