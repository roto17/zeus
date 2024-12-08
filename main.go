package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/gorilla/websocket"
	"github.com/roto17/zeus/lib/config"
	"github.com/roto17/zeus/lib/database"
	"github.com/roto17/zeus/lib/router"
	"github.com/twilio/twilio-go"
	openapi "github.com/twilio/twilio-go/rest/api/v2010"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func main() {

	// Send the SMS verification code
	sendVerificationSMS()

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

func sendVerificationSMS() {

	accountSid := "ACdaa6c8500198d8beffdd811844f14d75"
	authToken := "a8020f6379f14a396ea1bbad41f0c644"
	client := twilio.NewRestClientWithParams(twilio.ClientParams{
		Username: accountSid,
		Password: authToken,
	})

	// Send SMS
	params := &openapi.CreateMessageParams{}
	params.SetTo("+212634894766")  // Replace with the recipient's phone number
	params.SetFrom("+17856293009") // Replace with your Twilio phone number
	params.SetBody("Your verification code is TZZZ: 123456")

	resp, err := client.Api.CreateMessage(params)

	if err != nil {
		fmt.Printf("Error sending message: %s\n", err)
	} else {
		fmt.Printf("Message SID: %s\n", *resp.Sid)
	}

}
