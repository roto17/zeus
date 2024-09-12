// File: main.go
package main

import (
	"fmt"

	"github.com/roto17/zeus/lib" // Replace with your actual module path
)

func main() {
	lib.LoadConfig()
	dbUser := lib.GetEnv("dbport")
	dbPassword := lib.GetEnv("dbname")

	fmt.Println("dbport:", dbUser)
	fmt.Println("dbname:", dbPassword)
}
