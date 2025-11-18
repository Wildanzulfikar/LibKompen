package controllers

import (
    "LibKompen/database"
	"LibKompen/models"
	"LibKompen/models/request"

	"github.com/gofiber/fiber/v2"
)

func CreateTenggat(c *fiber.Ctx) error {
	var payload []models.TenggatWaktu

	if err := c.BodyParser(&payload); err != nil {
		return c.Status(400).JSON(fiber.Map {
			"error" : err.Error(),
		})
	}

	tenggat := models.TenggatWaktu {
		WaktuMulai : payload.WaktuMulai,
		WaktuAkhir : payload.WaktuAkhir,
	}

	if err := models.DB.create(&tenggat).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error" : err.Error(),
		})
	}
}