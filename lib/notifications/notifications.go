package notifications

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/roto17/zeus/lib/utils"
)

const (
	readDeadline      = 60 * time.Second // Constant for read deadline
	maxWorkerPoolSize = 10               // Limit for worker pool size
)

// WebSocket upgrader with custom origin check
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		origin := r.Header.Get("Origin")
		return origin == "http://localhost:3000" || origin == "http://your-frontend-origin.com"
	},
}

// Store for connected clients and channel for broadcasting notifications
var (
	clients    sync.Map
	broadcast  = make(chan Notification)
	workerPool = make(chan Notification, maxWorkerPoolSize) // Use buffered channel for limited worker pool
	jsonPool   = sync.Pool{New: func() interface{} { return new(bytes.Buffer) }}
)

// Notification represents a notification message structure
type Notification struct {
	Content string   `json:"content"`
	From    string   `json:"from"`     // Sender's role
	ToRoles []string `json:"to_roles"` // Recipient roles
}

// NewNotification is a constructor for Notification
func NewNotification(from string, to []string, content string) Notification {
	return Notification{
		From:    from,
		ToRoles: to,
		Content: content,
	}
}

// AddClient registers a new WebSocket client
func AddClient(conn *websocket.Conn, role string) {
	clients.Store(conn, role)
}

// RemoveClient unregisters a WebSocket client
func RemoveClient(conn *websocket.Conn) {
	clients.Delete(conn)
}

// HandleMessages listens for incoming notifications and dispatches them to workers
func HandleMessages() {
	for notification := range broadcast {
		workerPool <- notification // Dispatch to worker pool
	}
}

// contains checks if a given role exists in the roles slice
func contains(roles []string, role string) bool {
	roleMap := make(map[string]struct{}, len(roles))
	for _, r := range roles {
		roleMap[r] = struct{}{}
	}
	_, exists := roleMap[role]
	return exists
}

// worker processes notifications and sends them to relevant clients
func worker(ctx context.Context) {
	for {
		select {
		case notification := <-workerPool:
			clients.Range(func(key, value interface{}) bool {
				client, ok := key.(*websocket.Conn)
				if !ok {
					fmt.Printf("Warning: Failed to assert client type\n")
					return true
				}

				role, ok := value.(string)
				if !ok {
					fmt.Printf("Warning: Failed to assert role type\n")
					return true
				}

				if contains(notification.ToRoles, role) {
					buffer := jsonPool.Get().(*bytes.Buffer)
					buffer.Reset()

					if err := json.NewEncoder(buffer).Encode(notification); err != nil {
						logErrorAndRemoveClient(client, err)
						return true
					}

					if err := client.WriteMessage(websocket.TextMessage, buffer.Bytes()); err != nil {
						logErrorAndRemoveClient(client, err)
					}

					jsonPool.Put(buffer) // Return buffer to pool
				}
				return true // Continue iteration
			})
		case <-ctx.Done():
			return // Exit worker if context is done
		}
	}
}

// StartWorkers initiates a specified number of worker goroutines
func StartWorkers(numWorkers int, ctx context.Context) {
	for i := 0; i < numWorkers; i++ {
		go worker(ctx) // Start worker goroutine
	}
}

// Notify creates a notification and sends it to the broadcast channel
func Notify(fromRole string, toRoles []string, msg string) {
	notification := NewNotification(fromRole, toRoles, msg)
	broadcast <- notification // Send the notification
}

// WSHandler manages WebSocket connections and authenticates users
func WSHandler(c *gin.Context) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		fmt.Printf("Error upgrading connection: %v\n", err)
		return
	}
	defer conn.Close() // Ensure the connection is closed when done

	tokenString := c.Query("token")
	role := utils.GetRoleFromToken(tokenString)

	if role == "" {
		closeConnectionWithError(conn, websocket.ClosePolicyViolation, "Invalid token")
		return
	}

	AddClient(conn, role)
	defer RemoveClient(conn)

	// Set read deadline to prevent resource leaks
	for {
		conn.SetReadDeadline(time.Now().Add(readDeadline)) // Reset deadline
		var msg Notification
		if err := conn.ReadJSON(&msg); err != nil {
			fmt.Printf("Error reading message from client %v: %v\n", conn.RemoteAddr(), err)
			break
		}
	}
}

// Helper function to log errors and remove client
func logErrorAndRemoveClient(client *websocket.Conn, err error) {
	fmt.Printf("Error: %v\n", err)
	client.Close()
	RemoveClient(client) // Clean up client
}

// Helper function to close a connection with an error message
func closeConnectionWithError(conn *websocket.Conn, code int, message string) {
	if err := conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(code, message)); err != nil {
		fmt.Printf("Error closing connection: %v\n", err)
	}
	conn.Close()
}
