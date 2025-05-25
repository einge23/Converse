package types

type CreateFriendRequest struct {
	RecipientID string `json:"recipient_id" binding:"required,uuid"`
}