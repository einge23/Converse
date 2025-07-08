### I. Define Signaling Message Structures

You need Go structs to represent the different types of messages that will be exchanged over your WebSocket connections.

1.  **`Message` Base Struct:**

    -   A common wrapper for all signaling messages to allow for a `Type` field (e.g., "offer", "answer", "candidate", "join_room", etc.) and a `Payload` (which will be an interface{} to hold the specific message data).

2.  **`JoinRoomMessage`:**

    -   `RoomID` (string): The ID of the room the user wants to join.

3.  **`SDPMessage` (Offer/Answer):**

    -   `SDP` (string): The actual SDP string.
    -   `Type` (string): "offer" or "answer" (though your base `Message` might already have this).
    -   `TargetID` (string): The `socket.ID` or `PeerID` of the intended recipient.
    -   `SenderID` (string): The `socket.ID` or `PeerID` of the sender.

4.  **`ICECandidateMessage`:**

    -   `Candidate` (string): The ICE candidate string.
    -   `SDPMLineIndex` (int): The media description index.
    -   `SDPMid` (string): The media stream ID.
    -   `TargetID` (string): The `socket.ID` or `PeerID` of the intended recipient.
    -   `SenderID` (string): The `socket.ID` or `PeerID` of the sender.

5.  **`PeerConnectedMessage` / `PeerDisconnectedMessage`:**
    -   `PeerID` (string): The ID of the peer that just joined or left the room. (Sent by the server to existing peers).

**Example Structs (in `models/signaling_messages.go`):**

```go
package models

// Message is the base wrapper for all signaling messages
type Message struct {
	Type    string      `json:"type"`    // e.g., "join_room", "offer", "answer", "candidate", "peer_joined", "peer_left"
	Payload interface{} `json:"payload"` // The specific message payload
}

// JoinRoomPayload is the payload for the "join_room" message
type JoinRoomPayload struct {
	RoomID string `json:"roomId"`
}

// SDPPayload is the payload for "offer" and "answer" messages
type SDPPayload struct {
	SDP      string `json:"sdp"`
	TargetID string `json:"targetId"` // The ID of the peer to send this to
	SenderID string `json:"senderId"` // The ID of the peer who sent this
}

// ICECandidatePayload is the payload for "candidate" messages
type ICECandidatePayload struct {
	Candidate     string `json:"candidate"`
	SDPMLineIndex *int   `json:"sdpMLineIndex"` // Pointers for optional fields
	SDPMid        *string `json:"sdpMid"`      // Pointers for optional fields
	TargetID      string `json:"targetId"`
	SenderID      string `json:"senderId"`
}

// PeerEventPayload is the payload for "peer_joined" and "peer_left" messages
type PeerEventPayload struct {
	PeerID string `json:"peerId"`
}
```

### II. WebSocket Message Handling Logic

Your existing WebSocket connection manager needs to be extended to handle these new message types.

1.  **Modify Your `HandleWebSocket` Function:**

    -   When you receive a raw message from a client, deserialize it into your base `Message` struct.
    -   Use a `switch` statement on `msg.Type` to determine how to process the payload.

2.  **`JoinRoom` Handler:**

    -   **Payload:** Expects `models.JoinRoomPayload`.
    -   **Action:**
        -   Associate the `socket.ID` (or whatever ID you use for the client) with the `RoomID`.
        -   Maintain a map, e.g., `map[string][]string` for `roomID -> list of peer IDs`.
        -   **Notify Existing Peers:** If there are other peers in the room, send a `PeerConnectedMessage` (type "peer_joined") to them, containing the `PeerID` of the newly joined client.
        -   **Notify New Peer:** Send a `PeerConnectedMessage` to the newly joined client for each _existing_ peer in the room. This tells the new peer who else is already there to connect to.

3.  **`Offer` Handler:**

    -   **Payload:** Expects `models.SDPPayload`.
    -   **Action:**
        -   Find the `ClientConnection` corresponding to `data.TargetID`.
        -   Relay the entire `models.Message` (with `Type: "offer"`) to that specific client.
        -   **Important:** Ensure `SenderID` is correctly set to the _actual sender's ID_ (not necessarily provided by the client, but derived from the current WebSocket connection).

4.  **`Answer` Handler:**

    -   **Payload:** Expects `models.SDPPayload`.
    -   **Action:**
        -   Find the `ClientConnection` corresponding to `data.TargetID`.
        -   Relay the entire `models.Message` (with `Type: "answer"`) to that specific client.
        -   **Important:** Ensure `SenderID` is correctly set.

5.  **`Candidate` Handler:**

    -   **Payload:** Expects `models.ICECandidatePayload`.
    -   **Action:**
        -   Find the `ClientConnection` corresponding to `data.TargetID`.
        -   Relay the entire `models.Message` (with `Type: "candidate"`) to that specific client.
        -   **Important:** Ensure `SenderID` is correctly set.

6.  **`Disconnect` Handler:**
    -   **Action:**
        -   When a client disconnects, remove its `PeerID` from any rooms it was in.
        -   **Notify Remaining Peers:** Send a `PeerDisconnectedMessage` (type "peer_left") to all remaining peers in that room, informing them that a peer has left. This is crucial for cleaning up `RTCPeerConnection` objects on the client side.

### III. Server-Side Data Structures for Rooms and Peers

You'll need a way to manage which peers are in which rooms.

1.  **`RoomManager` or similar:**
    -   A Go `struct` (or just a `map`) that holds the state of your rooms.
    -   **`rooms` map:** `map[string]map[string]struct{}` (or `map[string]map[string]*ClientConnection`)
        -   Outer key: `RoomID` (string)
        -   Inner key: `PeerID` (string, the client's WebSocket connection ID)
        -   Value: `struct{}` (or the `*ClientConnection` if your manager holds it). Using `struct{}` is memory efficient if you just need a set of IDs.
    -   **Mutex:** Use a `sync.RWMutex` to protect access to this map, as multiple WebSocket goroutines will be reading from and writing to it concurrently.

**Example `RoomManager` (in `signaling/room_manager.go`):**

```go
package signaling

import (
	"sync"
)

// Connection represents an active WebSocket client connection.
// You likely have this in your existing WebSocket manager.
type Connection interface {
	ID() string
	WriteJSON(v interface{}) error
	// Add other methods as needed, e.g., ReadJSON for reading specific types
}

// RoomManager manages active rooms and the peers within them.
type RoomManager struct {
	mu    sync.RWMutex
	rooms map[string]map[string]Connection // roomID -> peerID -> Connection
}

func NewRoomManager() *RoomManager {
	return &RoomManager{
		rooms: make(map[string]map[string]Connection),
	}
}

// JoinRoom adds a peer to a specific room.
// Returns a list of existing peers in the room before this peer joined.
func (rm *RoomManager) JoinRoom(roomID string, conn Connection) []Connection {
	rm.mu.Lock()
	defer rm.mu.Unlock()

	if _, ok := rm.rooms[roomID]; !ok {
		rm.rooms[roomID] = make(map[string]Connection)
	}

	existingPeers := make([]Connection, 0, len(rm.rooms[roomID]))
	for _, existingConn := range rm.rooms[roomID] {
		existingPeers = append(existingPeers, existingConn)
	}

	rm.rooms[roomID][conn.ID()] = conn
	return existingPeers
}

// LeaveRoom removes a peer from a room.
// Returns a list of remaining peers in the room.
func (rm *RoomManager) LeaveRoom(roomID string, peerID string) []Connection {
	rm.mu.Lock()
	defer rm.mu.Unlock()

	if room, ok := rm.rooms[roomID]; ok {
		delete(room, peerID)
		if len(room) == 0 {
			delete(rm.rooms, roomID) // Clean up empty room
		}
		remainingPeers := make([]Connection, 0, len(room))
		for _, conn := range room {
			remainingPeers = append(remainingPeers, conn)
		}
		return remainingPeers
	}
	return nil
}

// GetPeerConnection retrieves a specific peer's connection by ID from a room.
func (rm *RoomManager) GetPeerConnection(roomID, peerID string) Connection {
	rm.mu.RLock()
	defer rm.mu.RUnlock()
	if room, ok := rm.rooms[roomID]; ok {
		return room[peerID]
	}
	return nil
}

// GetRoomPeers returns all connections in a given room.
func (rm *RoomManager) GetRoomPeers(roomID string) []Connection {
	rm.mu.RLock()
	defer rm.mu.RUnlock()
	if room, ok := rm.rooms[roomID]; ok {
		peers := make([]Connection, 0, len(room))
		for _, conn := range room {
			peers = append(peers, conn)
		}
		return peers
	}
	return nil
}

// GetRoomsForPeer returns a list of room IDs a peer is currently in.
// This might be useful for cleanup on disconnect.
func (rm *RoomManager) GetRoomsForPeer(peerID string) []string {
    rm.mu.RLock()
    defer rm.mu.RUnlock()

    var roomIDs []string
    for roomID, peers := range rm.rooms {
        if _, ok := peers[peerID]; ok {
            roomIDs = append(roomIDs, roomID)
        }
    }
    return roomIDs
}
```

### IV. Integration with Your WebSocket Manager

Your existing WebSocket manager will likely have a loop that reads messages from a client and dispatches them. This is where you'll integrate the new signaling logic.

```go
package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"

	"your_project_path/models" // Adjust import path
	"your_project_path/signaling" // Adjust import path
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // Allow all origins, adjust for production
	},
}

// WebSocketClient represents a client connection in your existing manager
// Make sure this satisfies the `signaling.Connection` interface
type WebSocketClient struct {
	conn *websocket.Conn
	id   string
}

func (c *WebSocketClient) ID() string {
	return c.id
}

func (c *WebSocketClient) WriteJSON(v interface{}) error {
	return c.conn.WriteJSON(v)
}

// WebSocketManager holds all active connections and the room manager
type WebSocketManager struct {
	connections map[string]*WebSocketClient
	roomManager *signaling.RoomManager
	// Add a mutex if `connections` can be accessed concurrently by other parts of your app
}

func NewWebSocketManager() *WebSocketManager {
	return &WebSocketManager{
		connections: make(map[string]*WebSocketClient),
		roomManager: signaling.NewRoomManager(),
	}
}

func (wm *WebSocketManager) HandleWebSocket(c *gin.Context) {
	ws, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println("upgrade:", err)
		return
	}
	defer ws.Close()

	// Assign a unique ID to the client
	// You might use a UUID generator here
	clientID := c.Query("clientId") // Or generate internally
	if clientID == "" {
		clientID = ws.RemoteAddr().String() // Simple fallback
	}

	client := &WebSocketClient{conn: ws, id: clientID}
	wm.connections[clientID] = client
	log.Printf("Client %s connected\n", clientID)

	defer func() {
		// Cleanup on disconnect
		log.Printf("Client %s disconnected\n", clientID)
		delete(wm.connections, clientID)
		disconnectedRooms := wm.roomManager.GetRoomsForPeer(clientID)
		for _, roomID := range disconnectedRooms {
			remainingPeers := wm.roomManager.LeaveRoom(roomID, clientID)
			// Notify remaining peers in the room
			peerLeftMsg := models.Message{
				Type: "peer_left",
				Payload: models.PeerEventPayload{
					PeerID: clientID,
				},
			}
			for _, peerConn := range remainingPeers {
				if err := peerConn.WriteJSON(peerLeftMsg); err != nil {
					log.Printf("Error sending peer_left to %s: %v\n", peerConn.ID(), err)
				}
			}
		}
	}()

	for {
		var msg models.Message
		err := ws.ReadJSON(&msg)
		if err != nil {
			if websocket.IsCloseError(err, websocket.CloseGoingAway, websocket.CloseNormalClosure) {
				log.Println("Websocket closed normally.")
			} else {
				log.Println("read:", err)
			}
			break
		}

		log.Printf("Received message from %s: Type=%s, Payload=%+v\n", clientID, msg.Type, msg.Payload)

		// Dispatch message based on type
		switch msg.Type {
		case "join_room":
			var payload models.JoinRoomPayload
			if err := mapstructure.Decode(msg.Payload, &payload); err != nil { // Use mapstructure for decoding interface{}
				log.Printf("Error decoding join_room payload: %v\n", err)
				continue
			}
			existingPeers := wm.roomManager.JoinRoom(payload.RoomID, client)

			// 1. Notify newly joined client about existing peers
			for _, existingPeer := range existingPeers {
				peerJoinedMsg := models.Message{
					Type: "peer_joined",
					Payload: models.PeerEventPayload{
						PeerID: existingPeer.ID(),
					},
				}
				if err := client.WriteJSON(peerJoinedMsg); err != nil {
					log.Printf("Error sending peer_joined to new client %s: %v\n", clientID, err)
				}
			}

			// 2. Notify existing peers about the new client
			newPeerJoinedMsg := models.Message{
				Type: "peer_joined",
				Payload: models.PeerEventPayload{
					PeerID: clientID, // The new client's ID
				},
			}
			for _, existingPeer := range existingPeers {
				// Don't send to self
				if existingPeer.ID() != clientID {
					if err := existingPeer.WriteJSON(newPeerJoinedMsg); err != nil {
						log.Printf("Error sending peer_joined to existing peer %s: %v\n", existingPeer.ID(), err)
					}
				}
			}

		case "offer", "answer":
			var payload models.SDPPayload
			if err := mapstructure.Decode(msg.Payload, &payload); err != nil {
				log.Printf("Error decoding SDP payload: %v\n", err)
				continue
			}
			payload.SenderID = clientID // Always set sender ID from server context

			targetConn := wm.roomManager.GetPeerConnection(payload.RoomID, payload.TargetID) // Assuming RoomID is part of SDPPayload, or you determine it from client's known room
			if targetConn != nil {
				relayMsg := models.Message{
					Type:    msg.Type, // "offer" or "answer"
					Payload: payload,
				}
				if err := targetConn.WriteJSON(relayMsg); err != nil {
					log.Printf("Error relaying %s to %s: %v\n", msg.Type, payload.TargetID, err)
				}
			} else {
				log.Printf("Target peer %s not found for %s message from %s\n", payload.TargetID, msg.Type, clientID)
			}

		case "candidate":
			var payload models.ICECandidatePayload
			if err := mapstructure.Decode(msg.Payload, &payload); err != nil {
				log.Printf("Error decoding candidate payload: %v\n", err)
				continue
			}
			payload.SenderID = clientID // Always set sender ID from server context

			targetConn := wm.roomManager.GetPeerConnection(payload.RoomID, payload.TargetID) // Assuming RoomID is part of ICECandidatePayload
			if targetConn != nil {
				relayMsg := models.Message{
					Type:    msg.Type, // "candidate"
					Payload: payload,
				}
				if err := targetConn.WriteJSON(relayMsg); err != nil {
					log.Printf("Error relaying candidate to %s: %v\n", payload.TargetID, err)
				}
			} else {
				log.Printf("Target peer %s not found for candidate message from %s\n", payload.TargetID, clientID)
			}

		default:
			log.Printf("Unknown message type: %s from client %s\n", msg.Type, clientID)
		}
	}
}

func main() {
	r := gin.Default()
	manager := NewWebSocketManager()

	r.GET("/ws", manager.HandleWebSocket)

	log.Fatal(r.Run(":8080")) // Listen on port 8080
}
```

**Note:** For decoding the `interface{}` `Payload` into specific structs, you'll need a library like `github.com/mitchellh/mapstructure`.

```bash
go get github.com/mitchellh/mapstructure
```

### V. Gin Routing

You'll have a single WebSocket endpoint (e.g., `/ws`) that your clients connect to. Gin will just route the HTTP upgrade request to your `HandleWebSocket` function.

```go
// In your main.go or router setup:
func main() {
    r := gin.Default()
    wsManager := NewWebSocketManager() // Initialize your manager

    r.GET("/ws", wsManager.HandleWebSocket)

    r.Run(":8080") // Start the Gin server
}
```

### VI. Client-Side Implementation (Briefly)

Your JavaScript client will need to:

1.  **Connect WebSocket:** `const ws = new WebSocket("ws://localhost:8080/ws?clientId=your_unique_id");`
2.  **Send `join_room`:** As soon as the WebSocket connects, send a `join_room` message.
3.  **Handle Incoming Messages:**
    -   Parse incoming JSON messages into the `Message` base struct.
    -   Use a `switch` statement on `message.type`.
    -   **`peer_joined`:** When you receive this, if you're the initiator, create an `RTCPeerConnection` for this new peer and send an `offer`. If you're not the initiator, wait for an offer.
    -   **`offer`:** Set `remoteDescription`, create `answer`, set `localDescription`, send `answer` back.
    -   **`answer`:** Set `remoteDescription`.
    -   **`candidate`:** Add `RTCIceCandidate` to the `RTCPeerConnection`.
    -   **`peer_left`:** Close the `RTCPeerConnection` associated with that `PeerID` and clean up UI elements.
4.  **Send Outgoing Messages:**
    -   When `pc.onicecandidate` fires, send a `candidate` message.
    -   When `pc.createOffer()` is done, send an `offer` message.
    -   When `pc.createAnswer()` is done, send an `answer` message.

### Summary of the Flow:

1.  **Client A & B connect to `/ws`**. Your server assigns them `clientID`s.
2.  **Client A sends `join_room`** with `RoomID: "myRoom"`.
3.  Server receives `join_room`.
    -   If "myRoom" is empty, Client A is added.
    -   If "myRoom" has Client B, Server sends `peer_joined` (with Client A's ID) to Client B. Server also sends `peer_joined` (with Client B's ID) to Client A.
4.  **Client A receives `peer_joined` (Client B's ID)**. Client A now knows about Client B.
    -   Client A creates `RTCPeerConnection`.
    -   Client A gets local media (`getUserMedia`).
    -   Client A adds tracks to `RTCPeerConnection`.
    -   Client A calls `createOffer()`.
    -   Client A calls `setLocalDescription()` with the offer.
    -   Client A sends `offer` message to server (targeting Client B).
5.  **Server relays `offer`** from Client A to Client B.
6.  **Client B receives `offer`**.
    -   Client B creates `RTCPeerConnection`.
    -   Client B gets local media.
    -   Client B adds tracks.
    -   Client B calls `setRemoteDescription()` with the received offer.
    -   Client B calls `createAnswer()`.
    -   Client B calls `setLocalDescription()` with the answer.
    -   Client B sends `answer` message to server (targeting Client A).
7.  **Server relays `answer`** from Client B to Client A.
8.  **Client A receives `answer`**.
    -   Client A calls `setRemoteDescription()` with the received answer.
9.  **Both clients generate ICE candidates** (`onicecandidate` event fires).
    -   Clients send `candidate` messages to the server (targeting the other peer).
10. **Server relays `candidate` messages** between Client A and Client B.
11. **Both clients receive ICE candidates** and call `addIceCandidate()`.
12. **WebRTC establishes direct peer-to-peer connection!**

This setup provides a robust foundation for your WebRTC signaling server in Go using Gin and WebSockets. Remember to handle errors gracefully and consider security implications for production environments.
