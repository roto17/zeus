package translation

import (
	"github.com/roto17/zeus/lib/config"
)

// getTranslation fetches the translated message based on the tag and language
func GetTranslation(tag string, field string, lang string) string {

	// if config.TranslationMap == nil {
	// 	config.LoadTranslationFile()
	// }

	translationMap := config.TranslationMap

	if translationMap[lang] == nil {
		// Fallback to English if the language is not supported
		lang = "en"
	}

	translation := translationMap[lang][tag]
	if translation == "" {
		// Fallback if no specific message for the tag is found
		translation = translationMap["en"][tag]
	}

	return translation
}
