package controllers

import (
	"LibKompen/database"
	"LibKompen/models"

	"github.com/gofiber/fiber/v2"
)

func GetApprovals(c *fiber.Ctx) error {
	var approvals []models.ApprovalBebasPustaka
	database.DB.Find(&approvals)
	return c.JSON(approvals)
}

func CreateApproval(c *fiber.Ctx) error {
	var approval models.ApprovalBebasPustaka
	if err := c.BodyParser(&approval); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Gagal parse data"})
	}
	database.DB.Create(&approval)
	return c.JSON(approval)
}
