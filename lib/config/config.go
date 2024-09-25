package config

import (
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	"github.com/roto17/zeus/lib/logs"
)

// Create a variable to hold the translation map
var TranslationMap map[string]map[string]string
var Router *gin.Engine

func LoadTranslationFile() {

	file, err := os.Open("./lib/translation/i18n.json")
	if err != nil {
		logs.AddLog("Fatal", "roto", fmt.Sprintf("Failed to open the file: %v", err))
	}
	defer file.Close()

	// Read the file contents
	data, err := io.ReadAll(file)
	if err != nil {
		logs.AddLog("Fatal", "roto", fmt.Sprintf("Failed to read file: %v", err))
	}

	// Unmarshal JSON data into the translation map
	err = json.Unmarshal(data, &TranslationMap)
	if err != nil {
		logs.AddLog("Fatal", "roto", fmt.Sprintf("Failed to unmarshal JSON: %v", err))
	} else {
		logs.AddLog("Fatal", "roto", "I18n.json file loaded successefuly")
	}

}

// LoadConfig loads environment variables from a .env file
func LoadConfig() {
	LoadTranslationFile()
	err := godotenv.Load()
	if err != nil {
		logs.AddLog("Fatal", "roto", "Error loading .env file")
	}
}

// GetEnv retrieves the value of an environment variable
func GetEnv(key string) string {
	return os.Getenv(key)
}
