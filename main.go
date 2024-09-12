// File: main.go
package main

import (
	"fmt"

	"github.com/roto17/zeus/env" // Replace with your actual module path
)

func main() {

	dbUser := env.GetEnv("dbport")
	dbPassword := env.GetEnv("dbname")
	dbsslmode := env.GetEnv("dbsslmode")

	fmt.Println("dbport:", dbUser)
	fmt.Println("dbname:", dbPassword)
	fmt.Println("dbsslmode:", dbsslmode)
}
