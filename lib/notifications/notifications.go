package notifications

import (
	// "src/lib/utils" // Adjusted import path
	"fmt"
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/roto17/zeus/lib/utils"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		// Allow connections from localhost and your frontend's origin
		return r.Header.Get("Origin") == "http://localhost:3000" || r.Header.Get("Origin") == "http://your-frontend-origin.com"
	},
}

var (
	clients    = make(map[*websocket.Conn]string) // Use string to store role
	broadcast  = make(chan Notification)
	clientLock = sync.Mutex{}
	workerPool = make(chan Notification, 100) // Channel for worker pool
)

// Notification represents the notification message
type Notification struct {
	Content string   `json:"content"`
	From    string   `json:"from"` // Sender's role
	ToRoles []string `json:"to_roles"`
}

// AddClient adds a new WebSocket client
func AddClient(conn *websocket.Conn, role string) {
	clientLock.Lock()
	clients[conn] = role
	clientLock.Unlock()
}

// RemoveClient removes a WebSocket client
func RemoveClient(conn *websocket.Conn) {
	clientLock.Lock()
	delete(clients, conn)
	clientLock.Unlock()
}

// HandleMessages processes notifications in the worker pool
func HandleMessages() {
	for {
		notification := <-broadcast
		workerPool <- notification // Send to worker pool
	}
}

func contains(roles []string, role string) bool {
	for _, r := range roles {
		if r == role {
			return true
		}
	}
	return false
}

func worker() {
	for notification := range workerPool {
		clientLock.Lock()
		for client, role := range clients {
			if contains(notification.ToRoles, role) { // Send to any role in the list
				err := client.WriteJSON(notification)
				if err != nil {
					client.Close()
					delete(clients, client)
				}
			}
		}
		clientLock.Unlock()
	}
}

// StartWorkers starts a set number of worker goroutines
func StartWorkers(numWorkers int) {
	for i := 0; i < numWorkers; i++ {
		go worker() // Start each worker
	}
}

// RegisterUser handles user registration and sends notifications
func Notify(from_role string, to_roles []string, msg string) {
	// Create a notification to inform admins about the new registration
	notification := Notification{
		From:    from_role,
		ToRoles: to_roles, // Indicating that the message is for admins
		Content: msg,
	}
	broadcast <- notification // Send the notification
}

// WSHandler handles WebSocket connections
func WSHandler(c *gin.Context) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		fmt.Printf("Error upgrading connection: %v\n", err)
		return
	}
	defer conn.Close()

	tokenString := c.Query("token")
	role := utils.GetRoleFromToken(tokenString)

	if role == "" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Invalid token"})
		return
	}

	AddClient(conn, role)

	for {
		var msg Notification
		err := conn.ReadJSON(&msg)
		if err != nil {
			RemoveClient(conn)
			break
		}
	}
}
