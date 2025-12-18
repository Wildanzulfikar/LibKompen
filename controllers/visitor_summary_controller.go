package controllers

import (
	"LibKompen/database"
	"LibKompen/models"
	"LibKompen/utils"

	"github.com/gofiber/fiber/v2"
)

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

// /api/visitor-summary/check?kode_user=...&date=...
func CheckVisitorSummary(c *fiber.Ctx) error {
	db := database.DB
	kodeUser := c.Query("kode_user")
	date := c.Query("date")
	if kodeUser == "" || date == "" {
		return c.Status(400).JSON(fiber.Map{"error": "kode_user and date required"})
	}
	var count int64
	err := db.Model(&models.VisitorSummary{}).
		Where("kode_user = ? AND DATE(created_at) = ?", kodeUser, date).
		Count(&count).Error
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "DB error"})
	}
	exists := count > 0
	return c.JSON(fiber.Map{"exists": exists})
}

// /api/visitor-summary/recap
func RecapVisitorSummary(c *fiber.Ctx) error {
	db := database.DB
	var total int64
	var today int64
	var month int64

	// Total visitor summary
	db.Model(&models.VisitorSummary{}).Count(&total)

	// Hari ini
	db.Model(&models.VisitorSummary{}).
		Where("DATE(created_at) = CURDATE()").
		Count(&today)

	// Bulan ini
	db.Model(&models.VisitorSummary{}).
		Where("YEAR(created_at) = YEAR(CURDATE()) AND MONTH(created_at) = MONTH(CURDATE())").
		Count(&month)

	return c.JSON(fiber.Map{
		"total": total,
		"today": today,
		"month": month,
	})
}

// /api/visitor-summary/recap-per-jurusan
func RecapVisitorSummaryPerJurusan(c *fiber.Ctx) error {
	db := database.DB
	type Visitor struct {
		KodeUser string
	}
	var visitors []Visitor
	if err := db.Model(&models.VisitorSummary{}).Select("kode_user").Find(&visitors).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "DB error"})
	}

	jurusanCount := make(map[string]int)
	for _, v := range visitors {
		if len(v.KodeUser) < 4 {
			continue
		}
		kodeJurusan := v.KodeUser[2:4]
		jurusanCount[kodeJurusan]++
	}

	result := make([]fiber.Map, 0)
	for kode, nama := range utils.KodeJurusanMap {
		result = append(result, fiber.Map{
			"kode":    kode,
			"jurusan": nama,
			"total":   jurusanCount[kode],
		})
	}
	return c.JSON(result)
}

// /api/visitor-summary/recap-per-bulan
func RecapVisitorSummaryPerBulan(c *fiber.Ctx) error {
	db := database.DB
	var visitors []models.VisitorSummary
	if err := db.Find(&visitors).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "DB error"})
	}

	bulanCount := make(map[int]int)
	for _, v := range visitors {
		bulan := int(v.CreatedAt.Month())
		bulanCount[bulan]++
	}

	type BulanRecap struct {
		Bulan int `json:"bulan"`
		Total int `json:"total"`
	}
	var recaps []BulanRecap
	for i := 1; i <= 12; i++ {
		recaps = append(recaps, BulanRecap{
			Bulan: i,
			Total: bulanCount[i],
		})
	}
	return c.JSON(recaps)
}

// /api/visitor-summary
func GetAllVisitorSummary(c *fiber.Ctx) error {
	db := database.DB
	jurusan := c.Query("jurusan") 
	date := c.Query("date")       

	var visitors []models.VisitorSummary
	if err := db.Find(&visitors).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "DB error"})
	}

	type Visitor struct {
		ID        int64  `json:"id"`
		KodeUser  string `json:"kode_user"`
		CreatedAt string `json:"created_at"`
	}
	var result []Visitor
	for _, v := range visitors {
		if jurusan != "" && len(v.KodeUser) >= 4 {
			if v.KodeUser[2:4] != jurusan {
				continue
			}
		}
		if date != "" {
			if v.CreatedAt.Format("2006-01-02") != date {
				continue
			}
		}
		result = append(result, Visitor{
			ID:       v.ID,
			KodeUser: v.KodeUser,
			CreatedAt: v.CreatedAt.Format("2006-01-02 15:04:05"),
		})
	}
	return c.JSON(result)
}
