package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/gorilla/websocket"
	"github.com/roto17/zeus/lib/config"
	"github.com/roto17/zeus/lib/database"
	"github.com/roto17/zeus/lib/router"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func main() {
	// Load Config
	config.LoadConfig()

	// Initialize and load the database
	database.InitDB()

	// Create a cancellable context
	ctx, cancel := context.WithCancel(context.Background())

	// Set up signal handling to cancel the context on shutdown
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)

	// Run the server in a goroutine
	go func() {
		if err := router.InitRouter(ctx).Run(config.GetEnv("app_running_port")); err != nil {
			log.Fatal("Failed to run server:", err)
		}
	}()

	// Wait for a termination signal
	<-signalChan
	log.Println("Shutting down gracefully...")
	cancel() // Cancel the context
	log.Println("Shutdown complete")
}
