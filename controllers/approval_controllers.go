package controllers

import (
	"LibKompen/models"
	"LibKompen/services"
	"github.com/gofiber/fiber/v2"
)

func GetApprovals(c *fiber.Ctx) error {
	       approvals, err := services.GetAllApprovals()
	       if err != nil {
		       return c.Status(500).JSON(fiber.Map{"status": "error", "message": "Gagal mengambil data approvals", "error": err.Error()})
	       }
	       return c.JSON(fiber.Map{"status": "success", "data": approvals})
}

func CreateApproval(c *fiber.Ctx) error {
	       var approval models.ApprovalBebasPustaka
	       if err := c.BodyParser(&approval); err != nil {
		       return c.Status(400).JSON(fiber.Map{"status": "error", "message": "Gagal parse data"})
	       }
	       err := services.CreateApproval(&approval)
	       if err != nil {
		       return c.Status(500).JSON(fiber.Map{"status": "error", "message": "Gagal menyimpan data approval", "error": err.Error()})
	       }
	       return c.JSON(fiber.Map{"status": "success", "data": approval})
}
