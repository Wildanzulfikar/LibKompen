package services

import (
	"LibKompen/database"
	"LibKompen/models"
	"golang.org/x/crypto/bcrypt"
	"errors"
	"regexp"
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

func ValidateRegister(username, password, confirm, name, email string) error {
	if len(username) < 4 {
		return errors.New("Username minimal 4 karakter")
	}
	for _, ch := range username {
		if !(ch >= 'a' && ch <= 'z') && !(ch >= 'A' && ch <= 'Z') && !(ch >= '0' && ch <= '9') {
			return errors.New("Username hanya boleh huruf dan angka")
		}
	}
	if len(password) < 6 {
		return errors.New("Password minimal 6 karakter")
	}
	if password != confirm {
		return errors.New("Konfirmasi password tidak cocok")
	}
	if name == "" {
		return errors.New("Nama wajib diisi")
	}
	if email == "" {
		return errors.New("Email wajib diisi")
	}
	if !isValidEmail(email) {
		return errors.New("Format email tidak valid")
	}
	return nil
}

func isValidEmail(email string) bool {
	if len(email) < 6 || len(email) > 50 {
		return false
	}
	re := regexp.MustCompile(`^[^@\s]+@[^@\s]+\.[^@\s]+$`)
	return re.MatchString(email)
}

func SetDefaultRole(role string) string {
	if role == "" {
		return "pustakawan"
	}
	return role
}
