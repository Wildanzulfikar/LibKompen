package controllers

import (
	"LibKompen/database"
	"LibKompen/models"

	"github.com/gofiber/fiber/v2"
)

func GetAllLoans(c *fiber.Ctx) error {
	var loans []models.OpacLoan

	if err := database.DB.Find(&loans).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"data": loans,
	})
}

func GetLoanByID(c *fiber.Ctx) error {
	loanID := c.Params("id")
	var loan models.OpacLoan

	if err := database.DB.Where("loan_id = ?", loanID).First(&loan).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Loan tidak ditemukan",
		})
	}

	return c.JSON(fiber.Map{
		"data": loan,
	})
}

func GetLoansByMember(c *fiber.Ctx) error {
	memberID := c.Params("member_id")
	var loans []models.OpacLoan

	if err := database.DB.Where("member_id = ?", memberID).Find(&loans).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"data": loans,
	})
}

func GetActiveLoans(c *fiber.Ctx) error {
	var loans []models.OpacLoan

	if err := database.DB.Where("is_return = ?", false).Find(&loans).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"data": loans,
	})
}

func CreateLoan(c *fiber.Ctx) error {
	var loan models.OpacLoan

	if err := c.BodyParser(&loan); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	if err := database.DB.Create(&loan).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Loan berhasil dibuat",
		"data":    loan,
	})
}

func UpdateLoan(c *fiber.Ctx) error {
	loanID := c.Params("id")
	var loan models.OpacLoan

	if err := database.DB.Where("loan_id = ?", loanID).First(&loan).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Loan tidak ditemukan",
		})
	}

	if err := c.BodyParser(&loan); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	if err := database.DB.Save(&loan).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "Loan berhasil diupdate",
		"data":    loan,
	})
}

func ReturnLoan(c *fiber.Ctx) error {
	loanID := c.Params("id")
	var loan models.OpacLoan

	if err := database.DB.Where("loan_id = ?", loanID).First(&loan).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Loan tidak ditemukan",
		})
	}

	loan.IsReturn = true

	if err := database.DB.Save(&loan).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "Buku berhasil dikembalikan",
		"data":    loan,
	})
}

func DeleteLoan(c *fiber.Ctx) error {
	loanID := c.Params("id")
	var loan models.OpacLoan

	if err := database.DB.Where("loan_id = ?", loanID).First(&loan).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Loan tidak ditemukan",
		})
	}

	if err := database.DB.Delete(&loan).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "Loan berhasil dihapus",
	})
}
