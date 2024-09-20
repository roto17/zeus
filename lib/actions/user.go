package actions

import (
	"github.com/roto17/zeus/lib/database"
	"github.com/roto17/zeus/lib/models"
)

func CreateUser(user *models.User) error {
	result := database.DB.Create(&user)
	return result.Error
}

func GetUser(id int) (models.User, error) {
	var user models.User
	result := database.DB.First(&user, id)
	return user, result.Error
}

func UpdateUser(user *models.User) error {
	result := database.DB.Save(&user)
	return result.Error
}

func DeleteUser(id int) error {
	result := database.DB.Delete(&models.User{}, id)
	return result.Error
}
