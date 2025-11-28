package services

import (
	"fmt"
	"LibKompen/database"
	"LibKompen/models"
)

func GetAllUsersBebasPustaka() ([]models.UsersBebasPustaka, error) {
	var users []models.UsersBebasPustaka
	result := database.DB.Find(&users)
	return users, result.Error
}

func DeleteUserBebasPustakaByID(id string) error {
	var user models.UsersBebasPustaka
	result := database.DB.First(&user, id)
	if result.Error != nil {
		return result.Error
	}

	var count int64
	database.DB.Table("visitor_summary").Where("dibuat_oleh = ?", user.IdUsers).Count(&count)
	if count > 0 {
		return fmt.Errorf("user_has_dependencies:%d", count)
	}

	return database.DB.Delete(&user).Error
}
