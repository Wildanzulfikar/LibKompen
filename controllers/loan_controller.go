package controllers

import (
	"LibKompen/services"

	"github.com/gofiber/fiber/v2"
)

func GetAllLoan(c *fiber.Ctx) error {
	page := c.QueryInt("page", 1)
	perPage := c.QueryInt("per_page", 10)
	search := c.Query("search", "")

	formatted, total, err := services.GetAllLoanFormatted(page, perPage, search)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{
		"meta": fiber.Map{
			"page":     page,
			"per_page": perPage,
			"total":    total,
		},
		"data": formatted,
	})
}

func GetLoanDetail(c *fiber.Ctx) error {
	memberID := c.Params("member_id")

	if memberID == "" {
		return c.Status(400).JSON(fiber.Map{
			"status":  "error",
			"message": "member_id wajib diisi",
		})
	}

	result, err := services.FetchLoanDetailByMember(memberID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"status":  "error",
			"message": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"status": "success",
		"data":   result,
	})
}

func GetLoanByMember(c *fiber.Ctx) error {
	memberID := c.Params("member_id")

	result, err := services.FetchLoanByMemberID(memberID)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"status":  "error",
			"message": err.Error(),
		})
	}

	return c.JSON(result)
}
