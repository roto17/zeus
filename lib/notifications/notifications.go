package notifications

import (
	"sync"

	"github.com/gorilla/websocket"
)

var (
	clients   = make(map[*websocket.Conn]bool) // Connected clients
	clientsMu sync.Mutex                       // Mutex for concurrent access
)

// AddClient adds a new WebSocket client
func AddClient(conn *websocket.Conn) {
	clientsMu.Lock()
	clients[conn] = true
	clientsMu.Unlock()
}

// RemoveClient removes a WebSocket client
func RemoveClient(conn *websocket.Conn) {
	clientsMu.Lock()
	delete(clients, conn)
	clientsMu.Unlock()
}

// SendNotification sends a notification to all connected clients
func SendNotification(message string) {
	clientsMu.Lock()
	for client := range clients {
		err := client.WriteJSON(message)
		if err != nil {
			client.Close() // Close the connection on error
			delete(clients, client)
		}
	}
	clientsMu.Unlock()
}
