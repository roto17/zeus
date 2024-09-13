package actions

import (
	"github.com/roto17/zeus/lib/models"
	"gorm.io/gorm"
)

func CreateUser(db *gorm.DB, user *models.User) error {
	result := db.Create(&user)
	return result.Error
}

func GetUser(db *gorm.DB, id int) (models.User, error) {
	var user models.User
	result := db.First(&user, id)
	return user, result.Error
}

func UpdateUser(db *gorm.DB, user *models.User) error {
	result := db.Save(&user)
	return result.Error
}

func DeleteUser(db *gorm.DB, id int) error {
	result := db.Delete(&models.User{}, id)
	return result.Error
}
