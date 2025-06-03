package websocket

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

const (
	writeWait      = 10 * time.Second
	pongWait       = 60 * time.Second
	pingPeriod     = (pongWait * 9) / 10
	maxMessageSize = 512
)

type Client struct {
	hub    *Hub
	conn   *websocket.Conn
	send   chan []byte
	UserID string
}

func (c *Client) readPump() {
	defer func() {
		c.hub.unregister <- c
		c.conn.Close()
	}()

	c.conn.SetReadLimit(maxMessageSize)
    c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error {
        c.conn.SetReadDeadline(time.Now().Add(pongWait))
        return nil
    })

	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
            if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
                log.Printf("error: %v", err)
            }
            break
        }

		// Parse the incoming message
        var incomingMsg IncomingMessage
        if err := json.Unmarshal(message, &incomingMsg); err != nil {
            log.Printf("Error parsing message from %s: %v", c.UserID, err)
            c.hub.sendErrorToClient(c, "Invalid message format")
            continue
        }

        log.Printf("Received message from %s: type=%s, content=%s", c.UserID, incomingMsg.Type, incomingMsg.Content)        // Route the message based on type
        switch incomingMsg.Type {
        case MessageTypeNewMessage:
            c.hub.ProcessIncomingMessage(c, incomingMsg)
        case MessageTypeTyping:
            c.handleTypingIndicator(incomingMsg, true)
        case MessageTypeStopTyping:
            c.handleTypingIndicator(incomingMsg, false)
        case MessageTypePing:
            c.handlePing()
        default:
            log.Printf("Unknown message type from %s: %s", c.UserID, incomingMsg.Type)
            c.hub.sendErrorToClient(c, "Unknown message type")
        }
	}
}

// handleTypingIndicator handles typing indicator messages
func (c *Client) handleTypingIndicator(msg IncomingMessage, isTyping bool) {
    messageType := MessageTypeStopTyping
    if isTyping {
        messageType = MessageTypeTyping
    }

    typingMsg := OutgoingMessage{
        Type:     messageType,
        RoomID:   msg.RoomID,
        ThreadID: msg.ThreadID,
        SenderID: c.UserID,
    }

    // Route typing indicator to appropriate recipients
    if msg.RoomID != nil {
        c.hub.SendToRoom(*msg.RoomID, typingMsg, c.UserID)
    } else if msg.ThreadID != nil {
        c.hub.SendToThread(*msg.ThreadID, typingMsg, c.UserID)
    }
}

func (c *Client) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
                c.conn.WriteMessage(websocket.CloseMessage, []byte{})
                return
            }

			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)

			// Add queued messages to the current websocket message
			n := len(c.send)
			for range n {
                w.Write([]byte{'\n'})
                w.Write(<-c.send)
            }

			if err := w.Close(); err != nil {
                return
            }

		case <-ticker.C:
            c.conn.SetWriteDeadline(time.Now().Add(writeWait))
            if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
                return
            }
        }
	}
}

// handlePing responds to ping messages with a pong
func (c *Client) handlePing() {
    pongMsg := OutgoingMessage{
        Type:     MessageTypePong,
        SenderID: c.UserID,
    }

    messageBytes, err := json.Marshal(pongMsg)
    if err != nil {
        log.Printf("Error marshaling pong message for %s: %v", c.UserID, err)
        return
    }

    select {
    case c.send <- messageBytes:
        log.Printf("Pong sent to %s", c.UserID)
    default:
        log.Printf("Failed to send pong to %s", c.UserID)
    }
}

func ServeWS(hub *Hub, c *gin.Context) {
	userID, exists := c.Get("user_id")
    if !exists {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
        return
    }

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println(err)
		return
	}

	client := &Client{
		hub: hub,
		conn: conn,
		send:   make(chan []byte, 256),
		UserID: userID.(string),
	}

	client.hub.register <- client

	go client.writePump()
    go client.readPump()
}
