package websocket

import "time"

// WebSocketMessageType represents the type of WebSocket message
type WebSocketMessageType string

const (
	MessageTypeNewMessage    WebSocketMessageType = "new_message"
	MessageTypeUserJoined    WebSocketMessageType = "user_joined"
	MessageTypeUserLeft      WebSocketMessageType = "user_left"
	MessageTypeTyping        WebSocketMessageType = "typing"
	MessageTypeStopTyping    WebSocketMessageType = "stop_typing"
	MessageTypeError         WebSocketMessageType = "error"
	MessageTypePing          WebSocketMessageType = "ping"
	MessageTypePong          WebSocketMessageType = "pong"
)

// IncomingMessage represents a message received from a client
type IncomingMessage struct {
	Type      WebSocketMessageType `json:"type"`
	RoomID    *string             `json:"room_id,omitempty"`
	ThreadID  *string             `json:"thread_id,omitempty"`
	Content   string              `json:"content"`
	ContentType string            `json:"content_type,omitempty"`
}

// OutgoingMessage represents a message sent to clients
type OutgoingMessage struct {
	Type        WebSocketMessageType `json:"type"`
	MessageID   string              `json:"message_id,omitempty"`
	RoomID      *string             `json:"room_id,omitempty"`
	ThreadID    *string             `json:"thread_id,omitempty"`
	SenderID    string              `json:"sender_id"`
	Content     string              `json:"content"`
	ContentType string              `json:"content_type"`
	CreatedAt   time.Time           `json:"created_at"`
	Error       string              `json:"error,omitempty"`
}