// File: lib/conf.go
package env

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

// LoadConfig loads environment variables from a .env file
func LoadConfig() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

// GetEnv retrieves the value of an environment variable
func GetEnv(key string) string {
	LoadConfig()
	return os.Getenv(key)
}
