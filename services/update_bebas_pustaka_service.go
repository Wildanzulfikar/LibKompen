package services

import (
	"LibKompen/database"
	"LibKompen/models"
	"encoding/json"
	"net/http"

	"fmt"

	"github.com/gofiber/fiber/v2"
)

func UpdateBebasPustakaService(c *fiber.Ctx) error {
	var rawUID interface{}
	if v := c.Locals("id_users"); v != nil {
		rawUID = v
	} else if v := c.Locals("user_id"); v != nil {
		rawUID = v
	} else if v := c.Locals("idUsers"); v != nil {
		rawUID = v
	} else if v := c.Locals("id"); v != nil {
		rawUID = v
	}

	var userID uint
	switch v := rawUID.(type) {
	case uint:
		userID = v
	case int:
		userID = uint(v)
	case float64:
		userID = uint(v)
	case string:
		var tmp int
		fmt.Sscanf(v, "%d", &tmp)
		userID = uint(tmp)
	default:
		return c.Status(401).JSON(fiber.Map{"status": "error", "message": "ID user tidak ditemukan di token"})
	}

	var body map[string]interface{}
	if err := c.BodyParser(&body); err != nil {
		return c.Status(400).JSON(fiber.Map{"status": "error", "message": "Body request tidak valid"})
	}

	rawKode := body["kode_user"]
	if rawKode == nil {
		return c.Status(400).JSON(fiber.Map{"status": "error", "message": "kode_user wajib diisi"})
	}

	kodeUser := fmt.Sprint(rawKode)

	keterangan, _ := body["keterangan"].(string)

	var count int64
	database.DB.Table("status_approval").Where("kode_user = ?", kodeUser).Count(&count)
	if count == 0 {
		approval := models.ApprovalBebasPustaka{
			KodeUser:   kodeUser,
			Keterangan: keterangan,
			IDUsers:    uint(userID),
		}
		database.DB.Create(&approval)
	} else {
		database.DB.Table("status_approval").Where("kode_user = ?", kodeUser).Updates(map[string]interface{}{
			"keterangan": keterangan,
			"id_users":   uint(userID),
		})
	}

	// Setelah update, cek status pinjaman real-time
	resp, err := http.Get("http://localhost:8080/loan?member_id=" + kodeUser)
	if err == nil {
		defer resp.Body.Close()
		var loanPayload map[string]interface{}
		if err := json.NewDecoder(resp.Body).Decode(&loanPayload); err == nil {
			loansData, ok := loanPayload["data"].([]interface{})
			if ok {
				adaTanggungan := false
				for _, loanItem := range loansData {
					loanMap, ok := loanItem.(map[string]interface{})
					if !ok {
						continue
					}
					isReturn, ok := loanMap["is_return"].(bool)
					if !ok || !isReturn {
						adaTanggungan = true
						break
					}
				}
				// Update status_approval dan keterangan di DB jika perlu
				if adaTanggungan {
					database.DB.Table("status_approval").Where("kode_user = ?", kodeUser).Updates(map[string]interface{}{
						"status_approval": false,
					})
				} else {
					database.DB.Table("status_approval").Where("kode_user = ?", kodeUser).Updates(map[string]interface{}{
						"status_approval": true,
						"keterangan":      "-",
					})
				}
			}
		}
	}

	return c.JSON(fiber.Map{
		"status":     "success",
		"message":    "Keterangan berhasil diperbarui. Status bebas pustaka akan otomatis mengikuti status pinjaman.",
		"kode_user":  kodeUser,
		"keterangan": keterangan,
	})
}
