package websocket

import (
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
}

func NewHub() *Hub {
    return &Hub{
        clients:     make(map[*Client]bool),
        broadcast:   make(chan []byte),
        register:    make(chan *Client),
        unregister:  make(chan *Client),
        userClients: make(map[string]*Client),
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