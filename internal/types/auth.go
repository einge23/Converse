package types

import "time"

type RegisterRequest struct {
    Username    string `json:"username" binding:"required,min=3,max=50"`
    Email       string `json:"email" binding:"required,email"`
    Password    string `json:"password" binding:"required,min=6"`
    DisplayName string `json:"display_name"`
}

type LoginRequest struct {
    Email    string `json:"email" binding:"omitempty,email"`
    Password string `json:"password" binding:"required,min=6"`
    Username string `json:"username" binding:"omitempty,min=3,max=50"`
}

type AuthResponse struct {
    Token        string    `json:"token"`
    RefreshToken string    `json:"refresh_token"`
    ExpiresAt    time.Time `json:"expires_at"`
    SessionID    string    `json:"session_id"`
}

type RefreshTokenRequest struct {
    RefreshToken string `json:"refresh_token" binding:"required"`
}