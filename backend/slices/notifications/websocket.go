package notifications

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

const (
	// Time allowed to write a message to the peer
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer
	pongWait = 60 * time.Second

	// Send pings to peer with this period (must be less than pongWait)
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer
	maxMessageSize = 512
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		// Allow all origins for development
		// In production, you should restrict this to your frontend domain
		return true
	},
}

// Client represents a WebSocket client connection
type Client struct {
	ID     string
	UserID string
	conn   *websocket.Conn
	send   chan []byte
	hub    *WebSocketManager
}

// WebSocketManager manages all active WebSocket connections
type WebSocketManager struct {
	// Registered clients mapped by user ID
	clients map[string]map[*Client]bool

	// Register requests from clients
	register chan *Client

	// Unregister requests from clients
	unregister chan *Client

	// Broadcast messages to specific user
	broadcast chan *BroadcastMessage

	// Mutex for thread-safe operations
	mu sync.RWMutex
}

// BroadcastMessage contains message to broadcast to a specific user
type BroadcastMessage struct {
	UserID  string
	Message []byte
}

// NewWebSocketManager creates a new WebSocket manager
func NewWebSocketManager() *WebSocketManager {
	return &WebSocketManager{
		clients:    make(map[string]map[*Client]bool),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		broadcast:  make(chan *BroadcastMessage, 256),
	}
}

// Run starts the WebSocket manager's event loop
func (m *WebSocketManager) Run() {
	for {
		select {
		case client := <-m.register:
			m.mu.Lock()
			if _, ok := m.clients[client.UserID]; !ok {
				m.clients[client.UserID] = make(map[*Client]bool)
			}
			m.clients[client.UserID][client] = true
			m.mu.Unlock()
			log.Printf("WebSocket client registered for user: %s", client.UserID)

		case client := <-m.unregister:
			m.mu.Lock()
			if clients, ok := m.clients[client.UserID]; ok {
				if _, ok := clients[client]; ok {
					delete(clients, client)
					close(client.send)
					if len(clients) == 0 {
						delete(m.clients, client.UserID)
					}
				}
			}
			m.mu.Unlock()
			log.Printf("WebSocket client unregistered for user: %s", client.UserID)

		case message := <-m.broadcast:
			m.mu.RLock()
			if clients, ok := m.clients[message.UserID]; ok {
				for client := range clients {
					select {
					case client.send <- message.Message:
					default:
						// Client's send channel is full, close it
						close(client.send)
						delete(clients, client)
						if len(clients) == 0 {
							delete(m.clients, message.UserID)
						}
					}
				}
			}
			m.mu.RUnlock()
		}
	}
}

// BroadcastToUser sends a message to all connections for a specific user
func (m *WebSocketManager) BroadcastToUser(userID string, data interface{}) {
	message, err := json.Marshal(data)
	if err != nil {
		log.Printf("Error marshaling broadcast message: %v", err)
		return
	}

	m.broadcast <- &BroadcastMessage{
		UserID:  userID,
		Message: message,
	}
}

// GetConnectedUserCount returns the number of users currently connected
func (m *WebSocketManager) GetConnectedUserCount() int {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return len(m.clients)
}

// GetUserConnectionCount returns the number of connections for a specific user
func (m *WebSocketManager) GetUserConnectionCount(userID string) int {
	m.mu.RLock()
	defer m.mu.RUnlock()
	if clients, ok := m.clients[userID]; ok {
		return len(clients)
	}
	return 0
}

// readPump pumps messages from the WebSocket connection to the hub
func (c *Client) readPump() {
	defer func() {
		c.hub.unregister <- c
		c.conn.Close()
	}()

	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetReadLimit(maxMessageSize)
	c.conn.SetPongHandler(func(string) error {
		c.conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("WebSocket error: %v", err)
			}
			break
		}

		// Handle incoming messages if needed (e.g., acknowledgments)
		log.Printf("Received message from user %s: %s", c.UserID, message)
	}
}

// writePump pumps messages from the hub to the WebSocket connection
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
				// The hub closed the channel
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)

			// Add queued messages to the current WebSocket message
			n := len(c.send)
			for i := 0; i < n; i++ {
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

// ServeWS handles WebSocket requests from clients
func (m *WebSocketManager) ServeWS(w http.ResponseWriter, r *http.Request, userID string) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("WebSocket upgrade error: %v", err)
		return
	}

	client := &Client{
		ID:     generateClientID(),
		UserID: userID,
		conn:   conn,
		send:   make(chan []byte, 256),
		hub:    m,
	}

	m.register <- client

	// Start client goroutines
	go client.writePump()
	go client.readPump()
}

// generateClientID generates a unique client ID
func generateClientID() string {
	return time.Now().Format("20060102150405") + "-" + randomString(8)
}

// randomString generates a random string of specified length
func randomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[time.Now().UnixNano()%int64(len(charset))]
	}
	return string(b)
}
