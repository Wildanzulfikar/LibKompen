package controllers

import (
	"LibKompen/database"
	"LibKompen/models"

	"github.com/gofiber/fiber/v2"
)

// POST /api/visitor-summary

func CreateVisitorSummary(c *fiber.Ctx) error {
	db := database.DB

	var req struct {
		KodeUser string `json:"kode_user"`
	}
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request"})
	}
	if req.KodeUser == "" {
		return c.Status(400).JSON(fiber.Map{"error": "kode_user required"})
	}
	vs := models.VisitorSummary{KodeUser: req.KodeUser}
	if err := db.Create(&vs).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "DB error"})
	}
	return c.JSON(vs)
}

// GET /api/visitor-summary/check?kode_user=...&date=...
func CheckVisitorSummary(c *fiber.Ctx) error {
	db := database.DB
	kodeUser := c.Query("kode_user")
	date := c.Query("date")
	if kodeUser == "" || date == "" {
		return c.Status(400).JSON(fiber.Map{"error": "kode_user and date required"})
	}
	var count int64
	// Asumsikan ada kolom created_at di visitor_summary (jika belum ada, perlu migrasi DB)
	err := db.Model(&models.VisitorSummary{}).
		Where("kode_user = ? AND DATE(created_at) = ?", kodeUser, date).
		Count(&count).Error
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "DB error"})
	}
	exists := count > 0
	return c.JSON(fiber.Map{"exists": exists})
}
