package websocket

import (
	"converse/internal/models"
	"converse/internal/repositories"
	"encoding/json"
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
        origin := r.Header.Get("Origin")
        allowedOrigins := []string{
            "http://localhost:5173",
            "https://converse-ui-development.up.railway.app",
        }
        
        for _, allowed := range allowedOrigins {
            if origin == allowed {
                return true
            }
        }
        return false
    },
}

type Hub struct {
	// Registered clients
	clients map[*Client]bool

    // Inbound messages from clients
    broadcast chan []byte

    // Register requests from clients
    register chan *Client

    // Unregister requests from clients
    unregister chan *Client

    // UserID to client mapping for DMs
    userClients map[string]*Client
    mutex sync.RWMutex

    // Repository dependencies for routing
    messageRepo *repositories.MessageRepository
}

func NewHub() *Hub {
    return &Hub{
        clients:     make(map[*Client]bool),
        broadcast:   make(chan []byte),
        register:    make(chan *Client),
        unregister:  make(chan *Client),
        userClients: make(map[string]*Client),
        messageRepo: repositories.NewMessageRepository(),
    }
}

func (h *Hub) Run() {
    for {
        select {
        case client := <-h.register:
            h.clients[client] = true
            h.mutex.Lock()
            h.userClients[client.UserID] = client
            h.mutex.Unlock()
            log.Printf("Client %s connected", client.UserID)

        case client := <-h.unregister:
            if _, ok := h.clients[client]; ok {
                delete(h.clients, client)
                h.mutex.Lock()
                delete(h.userClients, client.UserID)
                h.mutex.Unlock()
                close(client.send)
                log.Printf("Client %s disconnected", client.UserID)
            }
        case message := <-h.broadcast:
            for client := range h.clients {
                select {
                case client.send <- message:
                default:
                    close(client.send)
                    delete(h.clients, client)
                    h.mutex.Lock()
                    delete(h.userClients, client.UserID)
                    h.mutex.Unlock()
                }
            }
        }
    }
}

func (h *Hub) SendToUser(userID string, message []byte) {
    h.mutex.RLock()
    client, exists := h.userClients[userID]
    h.mutex.RUnlock()

    if exists {
        select {
        case client.send <- message:
        default:
            close(client.send)
            delete(h.clients, client)
            h.mutex.Lock()
            delete(h.userClients, userID)
            h.mutex.Unlock()
        }
    }
}

// SendToRoom sends a message to all members of a specific room
func (h *Hub) SendToRoom(roomID string, message OutgoingMessage, excludeUserID string) error {
    // Get all room members from database
    roomMembers, err := h.getRoomMembers(roomID)
    if err != nil {
        log.Printf("Error getting room members for room %s: %v", roomID, err)
        return err
    }

    // Convert message to JSON
    messageBytes, err := json.Marshal(message)
    if err != nil {
        log.Printf("Error marshaling message: %v", err)
        return err
    }

    // Send to all connected room members except the sender
    h.mutex.RLock()
    defer h.mutex.RUnlock()
    
    for _, member := range roomMembers {
        if member.UserID != excludeUserID {
            if client, exists := h.userClients[member.UserID]; exists {
                select {
                case client.send <- messageBytes:
                    log.Printf("Message sent to user %s in room %s", member.UserID, roomID)
                default:
                    // Client's send channel is full, clean up
                    h.cleanupClient(client)
                }
            }
        }
    }

    return nil
}

// SendToThread sends a message to participants of a DM thread
func (h *Hub) SendToThread(threadID string, message OutgoingMessage, excludeUserID string) error {
    // Get thread participants from database
    participants, err := h.getThreadParticipants(threadID)
    if err != nil {
        log.Printf("Error getting thread participants for thread %s: %v", threadID, err)
        return err
    }

    // Convert message to JSON
    messageBytes, err := json.Marshal(message)
    if err != nil {
        log.Printf("Error marshaling message: %v", err)
        return err
    }

    // Send to all connected participants except the sender
    h.mutex.RLock()
    defer h.mutex.RUnlock()
    
    for _, userID := range participants {
        if userID != excludeUserID {
            if client, exists := h.userClients[userID]; exists {
                select {
                case client.send <- messageBytes:
                    log.Printf("Message sent to user %s in thread %s", userID, threadID)
                default:
                    // Client's send channel is full, clean up
                    h.cleanupClient(client)
                }
            }
        }
    }

    return nil
}

// ProcessIncomingMessage handles routing of incoming messages from clients
func (h *Hub) ProcessIncomingMessage(client *Client, incomingMsg IncomingMessage) {
    // Validate message
    if incomingMsg.Content == "" {
        h.sendErrorToClient(client, "Message content cannot be empty")
        return
    }

    // Create message model for database storage
    message := &models.Message{
        RoomID:      incomingMsg.RoomID,
        ThreadID:    incomingMsg.ThreadID,
        SenderID:    &client.UserID,
        Content:     incomingMsg.Content,
        ContentType: incomingMsg.ContentType,
    }

    if message.ContentType == "" {
        message.ContentType = "text"
    }

    // Store message in database first
    if err := h.messageRepo.StoreMessage(message); err != nil {
        log.Printf("Error storing message: %v", err)
        h.sendErrorToClient(client, "Failed to store message")
        return
    }

    // Create outgoing message for broadcasting
    outgoingMsg := OutgoingMessage{
        Type:        MessageTypeNewMessage,
        MessageID:   message.MessageID,
        RoomID:      message.RoomID,
        ThreadID:    message.ThreadID,
        SenderID:    client.UserID,
        Content:     message.Content,
        ContentType: message.ContentType,
        CreatedAt:   message.CreatedAt,
    }    // Convert message to JSON for sending back to sender
    messageBytes, err := json.Marshal(outgoingMsg)
    if err != nil {
        log.Printf("Error marshaling message: %v", err)
        h.sendErrorToClient(client, "Failed to process message")
        return
    }

    // Send message back to sender first
    select {
    case client.send <- messageBytes:
        log.Printf("Message echoed back to sender %s", client.UserID)
    default:
        log.Printf("Failed to echo message back to sender %s", client.UserID)
    }

    // Route the message based on type
    if message.RoomID != nil {
        // It's a room message
        if err := h.SendToRoom(*message.RoomID, outgoingMsg, client.UserID); err != nil {
            log.Printf("Error sending message to room: %v", err)
        }
    } else if message.ThreadID != nil {
        // It's a DM thread message
        if err := h.SendToThread(*message.ThreadID, outgoingMsg, client.UserID); err != nil {
            log.Printf("Error sending message to thread: %v", err)
        }
    }
}

// Helper methods for database queries
func (h *Hub) getRoomMembers(roomID string) ([]*models.RoomMember, error) {
    var members []*models.RoomMember
    err := h.messageRepo.DB().Where("room_id = ?", roomID).Find(&members).Error
    return members, err
}

func (h *Hub) getThreadParticipants(threadID string) ([]string, error) {
    var thread models.DirectMessageThread
    err := h.messageRepo.DB().Where("thread_id = ?", threadID).First(&thread).Error
    if err != nil {
        return nil, err
    }
    return []string{thread.User1ID, thread.User2ID}, nil
}

func (h *Hub) sendErrorToClient(client *Client, errorMsg string) {
    errorMessage := OutgoingMessage{
        Type:  MessageTypeError,
        Error: errorMsg,
    }
    
    messageBytes, _ := json.Marshal(errorMessage)
    select {
    case client.send <- messageBytes:
    default:
        h.cleanupClient(client)
    }
}

func (h *Hub) cleanupClient(client *Client) {
    close(client.send)
    delete(h.clients, client)
    delete(h.userClients, client.UserID)
}