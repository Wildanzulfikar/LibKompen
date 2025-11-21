package services

import (
	"encoding/json"
	"fmt"
	"net/http"
	"github.com/gofiber/fiber/v2"
	"LibKompen/models"
	"LibKompen/database"
)

func UpdateBebasPustakaService(c *fiber.Ctx) error {
	uid := c.Locals("user_id")
	var userID uint
	switch v := uid.(type) {
	case uint:
		userID = v
	case int:
		userID = uint(v)
	case float64:
		userID = uint(v)
	default:
		return c.Status(401).JSON(fiber.Map{"status": "error", "message": "ID user tidak ditemukan di token"})
	}

	var body map[string]interface{}
	if err := c.BodyParser(&body); err != nil {
		return c.Status(400).JSON(fiber.Map{"status": "error", "message": "Body request tidak valid"})
	}

	kodeUser, ok := body["kode_user"].(string)
	if !ok || kodeUser == "" {
		return c.Status(400).JSON(fiber.Map{"status": "error", "message": "kode_user wajib diisi"})
	}

	status, ok := body["status"].(string)
	if !ok || (status != "bebas" && status != "tanggungan") {
		return c.Status(400).JSON(fiber.Map{"status": "error", "message": "status harus 'bebas' atau 'tanggungan'"})
	}

	keterangan, _ := body["keterangan"].(string)
	statusVal := status == "bebas"

	var count int64
	database.DB.Table("status_approval").Where("kode_user = ?", kodeUser).Count(&count)
	if count == 0 {
		approval := models.ApprovalBebasPustaka{
			KodeUser:       kodeUser,
			StatusApproval: statusVal,
			Keterangan:     keterangan,
			IDUsers:        uint(userID),
		}
		database.DB.Create(&approval)
	} else {
		database.DB.Table("status_approval").Where("kode_user = ?", kodeUser).Updates(map[string]interface{}{
			"status_approval": statusVal,
			"keterangan":      keterangan,
			"id_users":        uint(userID),
		})
	}

	if statusVal {
		resp, err := http.Get("http://localhost:8080/loan?member_id=" + kodeUser)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"status": "error", "message": "Gagal ambil loan dari OPAC"})
		}
		defer resp.Body.Close()

		var loanPayload map[string]interface{}
		if err := json.NewDecoder(resp.Body).Decode(&loanPayload); err != nil {
			return c.Status(500).JSON(fiber.Map{"status": "error", "message": "Gagal decode loan response"})
		}

		loansData, ok := loanPayload["data"].([]interface{})
		if !ok {
			return c.Status(500).JSON(fiber.Map{"status": "error", "message": "Data loan tidak valid"})
		}

		for _, loanItem := range loansData {
			loanMap, ok := loanItem.(map[string]interface{})
			if !ok {
				continue
			}
			isReturn, ok := loanMap["is_return"].(bool)
			if !ok || isReturn {
				continue
			}
			loanIdFloat, ok := loanMap["loan_id"].(float64)
			if !ok {
				continue
			}
			loanId := int(loanIdFloat)

			url := fmt.Sprintf("http://localhost:8080/loan/%d/return", loanId)
			req, _ := http.NewRequest("PUT", url, nil)
			http.DefaultClient.Do(req)
		}
	}

	return c.JSON(fiber.Map{
		"status":     "success",
		"message":    "Status bebas pustaka berhasil diperbarui",
		"kode_user":  kodeUser,
		"status_bebas":     status,
		"keterangan": keterangan,
	})
}
