package controllers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/gofiber/fiber/v2"
)

// /api/bebas-pustaka-range/:kode_user
func GetSikompenStatusBebasRange(c *fiber.Ctx) error {
	kodeUser := c.Params("kode_user")
	startDate := c.Query("start_date")
	endDate := c.Query("end_date")
	if kodeUser == "" || startDate == "" || endDate == "" {
		return c.Status(400).JSON(fiber.Map{"status": "error", "message": "missing parameter(s)"})
	}

	loanURL := fmt.Sprintf("http://localhost:8080/loan?member_id=%s", kodeUser)
	respLoan, err := http.Get(loanURL)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "Gagal ambil data loan dari Opac"})
	}
	defer respLoan.Body.Close()

	var loanPayload map[string]interface{}
	bodyBytes, _ := io.ReadAll(respLoan.Body)
	if err := json.Unmarshal(bodyBytes, &loanPayload); err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "Gagal decode loan: " + err.Error()})
	}

	loansData, ok := loanPayload["data"].([]interface{})
	statusApproval := "Lunas"
	if ok && len(loansData) > 0 {
		for _, loanItem := range loansData {
			loanMap, ok := loanItem.(map[string]interface{})
			if !ok {
				continue
			}
			loanDate, _ := loanMap["loan_date"].(string)
			if loanDate >= startDate && loanDate <= endDate {
				isReturn := loanMap["is_return"]
				if isReturn == false || isReturn == 0 || isReturn == "0" || isReturn == "false" {
					statusApproval = "Belum"
					break
				}
			}
		}
	}

	return c.JSON(fiber.Map{
		"status":          "success",
		"kode_user":       kodeUser,
		"status_approval": statusApproval,
		"start_date":      startDate,
		"end_date":        endDate,
	})
}
