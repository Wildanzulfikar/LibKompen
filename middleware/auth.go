package middleware

import (
	"LibKompen/database"
	"LibKompen/models"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"os"
	"strings"
)

func Protected() fiber.Handler {
	return func(c *fiber.Ctx) error {
		auth := c.Get("Authorization")
		if auth == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"success": false,
				"message": "Token tidak ditemukan, silakan login terlebih dahulu",
			})
		}

		parts := strings.Split(auth, " ")
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"success": false,
				"message": "Format token tidak valid, gunakan format: Bearer <token>",
			})
		}

		tokenString := parts[1]
		secret := os.Getenv("JWT_SECRET")
		if secret == "" {
			secret = "secret"
		}

		token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fiber.ErrUnauthorized
			}
			return []byte(secret), nil
		})

		if err != nil || !token.Valid {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"success": false,
				"message": "Token tidak valid atau telah kadaluarsa, silakan login kembali",
			})
		}

		claims := token.Claims.(jwt.MapClaims)
		var uid uint
		switch v := claims["user_id"].(type) {
		case float64:
			uid = uint(v)
		case int:
			uid = uint(v)
		case uint:
			uid = v
		case string:
			var tmp int
			fmt.Sscanf(v, "%d", &tmp)
			uid = uint(tmp)
		default:
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"success": false,
				"message": "Claim user_id tidak valid",
			})
		}

		var user models.UsersBebasPustaka
		database.DB.First(&user, uid)

		if user.IdUsers == 0 {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"success": false,
				"message": "User tidak ditemukan",
			})
		}

		c.Locals("user", user)
		c.Locals("user_id", user.IdUsers)

		return c.Next()
	}
}
