package controllers

import (
	"LibKompen/database"
	"LibKompen/models"
	"LibKompen/models/request"
	"time"

	"github.com/gofiber/fiber/v2"
)

func CreateTenggat(c *fiber.Ctx) error {
	userId := c.Locals("id_users").(uint)
	var payload request.TenggatWaktuRequest

	if err := c.BodyParser(&payload); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}

	tenggat := models.TenggatWaktu{
		IdUsers:    userId,
		WaktuMulai: payload.WaktuMulai,
		WaktuAkhir: payload.WaktuAkhir,
	}

	if err := database.DB.Create(&tenggat).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(201).JSON(tenggat)
}

func UpdateTenggat(c *fiber.Ctx) error {
	idTenggatWaktu := c.Params("id_tenggat_waktu")

	var tenggat models.TenggatWaktu

	if err := database.DB.First(&tenggat, "id_tenggat_waktu = ?", idTenggatWaktu).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{
			"error": "tenggat not found",
		})
	}

	var payload struct {
		WaktuMulai *time.Time `json:"waktu_mulai"`
		WaktuAkhir *time.Time `json:"waktu_akhir"`
	}

	c.BodyParser(&payload)

	if payload.WaktuMulai != nil {
		tenggat.WaktuMulai = *payload.WaktuMulai
	}

	if payload.WaktuAkhir != nil {
		tenggat.WaktuAkhir = *payload.WaktuMulai
	}

	database.DB.Save(&tenggat)
	return c.JSON(tenggat)
}

func GetAllTenggat(c *fiber.Ctx) error {
	var list []models.TenggatWaktu

	if err := database.DB.Find(&list).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(list)
}

func GetActiveTenggat(c *fiber.Ctx) error {
	var activeTenggat models.TenggatWaktu

	if err := database.DB.Order("id_tenggat_waktu DESC").First(&activeTenggat).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{
			"error": "no tenggat found",
		})
	}

	return c.JSON(activeTenggat)
}

func DeleteTenggat(c *fiber.Ctx) error {
	idTenggatWaktu := c.Params("id_tenggat_waktu")

	if err := database.DB.Delete(&models.TenggatWaktu{}, "id_tenggat_waktu = ?", idTenggatWaktu).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "tenggat waktu was deleted",
	})
}
