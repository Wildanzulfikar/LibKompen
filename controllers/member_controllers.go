package controllers

import (
	"LibKompen/database"
	"LibKompen/models"

	"github.com/gofiber/fiber/v2"
)

func GetAllMembers(c *fiber.Ctx) error {
	var members []models.OpacMember

	if err := database.DB.Find(&members).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"data": members,
	})
}

func GetMemberByID(c *fiber.Ctx) error {
	memberID := c.Params("id")
	var member models.OpacMember

	if err := database.DB.Where("member_id = ?", memberID).First(&member).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Member tidak ditemukan",
		})
	}

	return c.JSON(fiber.Map{
		"data": member,
	})
}

func CreateMember(c *fiber.Ctx) error {
	var member models.OpacMember

	if err := c.BodyParser(&member); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	if err := database.DB.Create(&member).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Member berhasil dibuat",
		"data":    member,
	})
}

func UpdateMember(c *fiber.Ctx) error {
	memberID := c.Params("id")
	var member models.OpacMember

	if err := database.DB.Where("member_id = ?", memberID).First(&member).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Member tidak ditemukan",
		})
	}

	if err := c.BodyParser(&member); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	if err := database.DB.Save(&member).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "Member berhasil diupdate",
		"data":    member,
	})
}

func DeleteMember(c *fiber.Ctx) error {
	memberID := c.Params("id")
	var member models.OpacMember

	if err := database.DB.Where("member_id = ?", memberID).First(&member).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Member tidak ditemukan",
		})
	}

	if err := database.DB.Delete(&member).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "Member berhasil dihapus",
	})
}
