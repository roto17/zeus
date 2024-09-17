package translation

import (
	"encoding/json"
	"io"

	// "fmt"

	"log"
	"os"
	"strings"
)

// getTranslation fetches the translated message based on the tag and language
func GetTranslation(tag string, field string, lang string) string {

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

	// Create a variable to hold the translation map
	var translationMap map[string]map[string]string

	// Unmarshal JSON data into the translation map
	err = json.Unmarshal(data, &translationMap)
	if err != nil {
		log.Fatalf("Failed to unmarshal JSON: %v", err)
	}

	// ******
	if translationMap[lang] == nil {
		// Fallback to English if the language is not supported
		lang = "en"
	}

	translation := translationMap[lang][tag]
	if translation == "" {
		// Fallback if no specific message for the tag is found
		translation = translationMap["en"][tag]
	}

	// Replace {Field} with the actual field name
	return strings.Replace(translation, "{Field}", field, 1)
}
