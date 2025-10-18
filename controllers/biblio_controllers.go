package controllers

import (
	"LibKompen/database"
	"LibKompen/models"

	"github.com/gofiber/fiber/v2"
)

func GetBiblio(c *fiber.Ctx) error {
	var biblios []models.OpacBiblio
	database.DB.Find(&biblios)
	return c.JSON(biblios)
}

func CreateBiblio(c *fiber.Ctx) error {
	var biblio models.OpacBiblio
	if err := c.BodyParser(&biblio); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Gagal parse data"})
	}
	database.DB.Create(&biblio)
	return c.JSON(biblio)
}
