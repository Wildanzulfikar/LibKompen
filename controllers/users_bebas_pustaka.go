package controllers

import (
	"LibKompen/services"
	"github.com/gofiber/fiber/v2"
	"strings"
)

func UsersBebasPustakaGetAll(c *fiber.Ctx) error {
	users, err := services.GetAllUsersBebasPustaka()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Gagal mengambil data users Bebas Pustaka",
			"error":   err.Error(),
		})
	}
	return c.JSON(fiber.Map{
		"status": "success",
		"data":   users,
	})
}

func DeleteUser(c *fiber.Ctx) error {
	id := c.Params("id")
	err := services.DeleteUserBebasPustakaByID(id)
	if err != nil {
		code := fiber.StatusInternalServerError
		msg := "Gagal menghapus user"
		if err.Error() == "record not found" {
			code = fiber.StatusNotFound
			msg = "User tidak ditemukan"
		} else if strings.HasPrefix(err.Error(), "user_has_dependencies") {
			code = fiber.StatusBadRequest
			msg = "User tidak dapat dihapus karena memiliki data terkait. Hapus atau pindahkan data terkait terlebih dahulu."
		}
		return c.Status(code).JSON(fiber.Map{
			"status":  "error",
			"message": msg,
			"error":   err.Error(),
		})
	}
	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "User berhasil dihapus",
	})
}
