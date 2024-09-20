package actions

import (
	"github.com/roto17/zeus/lib/database"
	"github.com/roto17/zeus/lib/models"
)

func Info(username string, message string) error {
	log := models.Log{LogType: "Info", Message: message, Username: username}
	result := database.DB.Create(&log)
	return result.Error
}

func Warning(username string, message string) error {
	log := models.Log{LogType: "Warning", Message: message, Username: username}
	result := database.DB.Create(&log)
	return result.Error
}

func Debug(username string, message string) error {
	log := models.Log{LogType: "Debug", Message: message, Username: username}
	result := database.DB.Create(&log)
	return result.Error
}

func Error(username string, message string) error {
	log := models.Log{LogType: "Error", Message: message, Username: username}
	result := database.DB.Create(&log)
	return result.Error
}

func Fatal(username string, message string) error {
	log := models.Log{LogType: "Fatal", Message: message, Username: username}
	result := database.DB.Create(&log)
	return result.Error
}
