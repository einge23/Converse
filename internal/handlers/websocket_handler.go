package handlers

import (
	"converse/internal/utils"
	"converse/internal/websocket"
	"net/http"

	"github.com/gin-gonic/gin"
)

type WebSocketHandler struct {
	hub *websocket.Hub
}

func NewWebSocketHandler(hub *websocket.Hub) *WebSocketHandler {
	return &WebSocketHandler{
		hub: hub,
	}
}

func (h *WebSocketHandler) HandleConnection(c *gin.Context) {
    // Get token from query parameter
    token := c.Query("token")
    if token == "" {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Token is required as query parameter"})
        return
    }

    // Validate the JWT token
    claims, err := utils.ValidateToken(token)
    if err != nil {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
        return
    }

    // Set user information in context for WebSocket connection
    c.Set("user_id", claims.UserID)
    c.Set("email", claims.Email)

    websocket.ServeWS(h.hub, c)
}