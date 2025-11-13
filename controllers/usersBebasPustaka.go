package controllers

import (
	"LibKompen/database"
	"LibKompen/models"

	"github.com/gofiber/fiber/v2"
)

func UsersBebasPustakaGetAll(c *fiber.Ctx) error {
	var users []models.UsersBebasPustaka
	database.DB.Find(&users)
	return c.JSON(users)
}
