
package controllers

import (
	"github.com/gofiber/fiber/v2"
	"LibKompen/services"
)
func GetAllLoan(c *fiber.Ctx) error {
	formatted, err := services.GetAllLoanFormatted()
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(formatted)
}

func GetLoanDetail(c *fiber.Ctx) error {
	loanID := c.Params("loan_id")
	result, err := services.FetchLoanDetail(loanID)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"status": "error", "message": err.Error()})
	}
	return c.JSON(result)
}

