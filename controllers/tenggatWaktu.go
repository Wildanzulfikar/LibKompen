package controllers

import (
	"LibKompen/models"
	"LibKompen/models/request"
	"github.com/gofiber/fiber/v2"
)

func CreateTenggat(c *fiber.Ctx) error {
	var payload []request.TenggatWaktuRequest

    if err := c.BodyParser(&payload); err != nil {
        return c.Status(400).JSON(fiber.Map{"error": err.Error()})
    }

    tenggat := models.TenggatWaktu{
        WaktuMulai: payload.WaktuMulai,
        WaktuAkhir: payload.WaktuAkhir,
    }

    if err := main.DB.Create(&tenggat).Error; err != nil {
        return c.Status(500).JSON(fiber.Map{"error": err.Error()})
    }

    return c.Status(201).JSON(tenggat)
}