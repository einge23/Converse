package types

type CreateFriendRequest struct {
	Username string `json:"username" binding:"required"`
}