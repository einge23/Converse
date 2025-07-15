package handlers

import (
	"converse/internal/services"
)

type UserHandler struct {
    userService *services.UserService
}

func NewUserHandler() *UserHandler {
    return &UserHandler{
        userService: services.NewUserService(),
    }
}