
package services

import (
	"LibKompen/database"
	"LibKompen/models"
	"golang.org/x/crypto/bcrypt"
)

func UpdateLastLogin(user *models.UsersBebasPustaka) error {
	return database.DB.Model(user).Update("last_login", user.LastLogin).Error
}

func FindUserByUsername(username string) (*models.UsersBebasPustaka, error) {
	var user models.UsersBebasPustaka
	result := database.DB.Where("username = ?", username).First(&user)
	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}

func CreateUser(user *models.UsersBebasPustaka) error {
	return database.DB.Create(user).Error
}

func CheckPassword(hash, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}

func FindUserByID(id interface{}) (*models.UsersBebasPustaka, error) {
	var user models.UsersBebasPustaka
	result := database.DB.First(&user, id)
	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}
