// File: main.go
package main

import (
	"context"
	"fmt"
	"log"

	"github.com/roto17/zeus/lib/database" // Replace with your actual module path
)

func main() {

	conn, err := database.DBConnection()
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close(context.Background())

	var greeting string
	err = conn.QueryRow(context.Background(), "SELECT 'Hello, world!'").Scan(&greeting)
	if err != nil {
		log.Fatal("QueryRow failed:", err)
	}

	fmt.Println(greeting)

	// conn.Close(context.Background())

}
