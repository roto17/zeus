package config

import (
	"encoding/json"
	"io"
	"log"
	"os"

	"github.com/joho/godotenv"
)

// LoadConfig loads environment variables from a .env file

// Create a variable to hold the translation map
var TranslationMap map[string]map[string]string

func LoadTranslationFile() {

	file, err := os.Open("./lib/translation/i18n.json")
	if err != nil {
		log.Fatalf("Failed to open the file: %v", err)
	}
	defer file.Close()

	// Read the file contents
	data, err := io.ReadAll(file)
	if err != nil {
		log.Fatalf("Failed to read file: %v", err)
	}

	// Unmarshal JSON data into the translation map
	err = json.Unmarshal(data, &TranslationMap)
	if err != nil {
		log.Fatalf("Failed to unmarshal JSON: %v", err)
	}

	// ******
}

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
