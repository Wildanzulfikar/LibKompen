package controllers

import (
	"LibKompen/database"
	"LibKompen/models"
	"LibKompen/models/request"
	"time"

	"github.com/gofiber/fiber/v2"
)

func CreateTenggat(c *fiber.Ctx) error {
	userIdInterface := c.Locals("id_users")
	userId, ok := userIdInterface.(uint)
	if !ok || userId == 0 {
		return c.Status(401).JSON(fiber.Map{
			"success": false,
			"message": "User tidak terautentikasi",
		})
	}

	var payload request.TenggatWaktuRequest

	if err := c.BodyParser(&payload); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}

	start, err := time.Parse("2006-01-02", payload.WaktuMulai)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Format waktu_mulai tidak valid"})
	}

	end, err := time.Parse("2006-01-02", payload.WaktuAkhir)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Format waktu_akhir tidak valid"})
	}

	tenggat := models.TenggatWaktu{
		IdUsers:    userId,
		WaktuMulai: start,
		WaktuAkhir: end,
	}

	var user models.UsersBebasPustaka
	if err := database.DB.First(&user, userId).Error; err != nil {
		return c.Status(401).JSON(fiber.Map{
			"success": false,
			"message": "User tidak ditemukan di database",
		})
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
		tenggat.WaktuAkhir = *payload.WaktuAkhir
	}

	database.DB.Save(&tenggat)
	return c.JSON(tenggat)
}

func GetAllTenggat(c *fiber.Ctx) error {
	var list []models.TenggatWaktu

	if err := database.DB.Preload("User").Find(&list).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(list)
}

func GetActiveTenggat(c *fiber.Ctx) error {
	var activeTenggat models.TenggatWaktu

	if err := database.DB.Preload("User").Order("id_tenggat_waktu DESC").First(&activeTenggat).Error; err != nil {
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
