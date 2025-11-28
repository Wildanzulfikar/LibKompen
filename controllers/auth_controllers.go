package controllers

import (
	"fmt"
	"os"
	"time"

	"LibKompen/models"
	"LibKompen/services"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
)

type AuthRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type RegisterRequest struct {
	Name            string `json:"name"`
	Email           string `json:"email"`
	Username        string `json:"username"`
	Password        string `json:"password"`
	ConfirmPassword string `json:"confirm"`
	Role            string `json:"role"`
}

func Register(c *fiber.Ctx) error {
	var body RegisterRequest
	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Format request tidak valid",
			"error":   err.Error(),
		})
	}
	body.Role = services.SetDefaultRole(body.Role)

	if err := services.ValidateRegister(body.Username, body.Password, body.ConfirmPassword, body.Name, body.Email); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": err.Error(),
		})
	}
	existing, err := services.FindUserByUsername(body.Username)
	if err == nil && existing != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Username sudah terdaftar, silakan gunakan username lain",
		})
	}

	hashedPasswordBytes, err := bcrypt.GenerateFromPassword([]byte(body.Password), bcrypt.DefaultCost)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Gagal memproses password",
			"error":   err.Error(),
		})
	}
	hashedPassword := string(hashedPasswordBytes)

	user := models.UsersBebasPustaka{
		Name:      body.Name,
		Email:     body.Email,
		Username:  body.Username,
		Password:  string(hashedPassword),
		Role:      body.Role,
		Status:    true,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	err = services.CreateUser(&user)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Gagal membuat user baru",
			"error":   err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"status":  "success",
		"message": "Registrasi berhasil",
		"data": fiber.Map{
			"user_id":  user.IdUsers,
			"username": user.Username,
			"name":     user.Name,
			"email":    user.Email,
			"role":     user.Role,
		},
	})
}

func Login(c *fiber.Ctx) error {
	var body AuthRequest

	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Format request tidak valid",
			"error":   err.Error(),
		})
	}

	if len(body.Username) < 4 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Username minimal 4 karakter",
		})
	}
	if len(body.Password) < 6 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Password minimal 6 karakter",
		})
	}

	fmt.Println("USERNAME INPUT:", body.Username)
	fmt.Println("PASSWORD INPUT:", body.Password)

	user, err := services.FindUserByUsername(body.Username)
	if err != nil || user == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status":  "error",
			"message": "Username atau password salah",
		})
	}

	fmt.Println("USER DB USERNAME:", user.Username)
	fmt.Println("USER DB PASSWORD HASH:", user.Password)

	err = services.CheckPassword(user.Password, body.Password)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status":  "error",
			"message": "Username atau password salah",
		})
	}

	var now = time.Now()
	user.LastLogin = now
	services.UpdateLastLogin(user)

	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		secret = "secret"
	}

	claims := jwt.MapClaims{
		"id_users": user.IdUsers,
		"username": user.Username,
		"role":     user.Role,
		"exp":      time.Now().Add(24 * time.Hour).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString([]byte(secret))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Gagal membuat token",
			"error":   err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Login berhasil",
		"data": fiber.Map{
			"access_token": signed,
			"token_type":   "Bearer",
			"expires_in":   86400,
			"user": fiber.Map{
				"user_id":  user.IdUsers,
				"username": user.Username,
				"name":     user.Name,
				"email":    user.Email,
				"role":     user.Role,
			},
		},
	})
}

func Me(c *fiber.Ctx) error {
	id := c.Locals("user_id")
	if id == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"success": false,
			"message": "Token tidak valid atau tidak ditemukan",
		})
	}

	user, err := services.FindUserByID(id)
	if err != nil || user == nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status":  "error",
			"message": "User tidak ditemukan",
		})
	}

	user.Password = ""
	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "Data user berhasil diambil",
		"data":    user,
	})
}

func Logout(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"success": true,
		"message": "Logout berhasil",
	})
}

func isValidEmail(email string) bool {
	if len(email) < 6 || len(email) > 50 {
		return false
	}
	at := false
	dot := false
	for i, c := range email {
		if c == '@' && i > 0 {
			at = true
		}
		if c == '.' && at {
			dot = true
		}
	}
	return at && dot
}
