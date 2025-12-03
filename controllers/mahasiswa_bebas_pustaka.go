package controllers

import (
	"LibKompen/services"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
)

func authMiddleware(c *fiber.Ctx) error {
	tokenString := c.Get("Authorization")
	if tokenString == "" {
		return c.Status(401).JSON(fiber.Map{"status": "error", "message": "Token tidak ditemukan"})
	}

	if len(tokenString) > 7 && tokenString[:7] == "Bearer " {
		tokenString = tokenString[7:]
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte("your-secret-key"), nil
	})
	if err != nil || !token.Valid {
		return c.Status(401).JSON(fiber.Map{"status": "error", "message": "Token tidak valid"})
	}
	return c.Next()
}

//	func GetMahasiswaBebasPustaka(c *fiber.Ctx) error {
//		filterMember := c.Query("member_id")
//		hasil, err := services.GetMahasiswaBebasPustakaService(filterMember)
//		if err != nil {
//			return c.Status(500).JSON(fiber.Map{"status": "error", "message": err.Error()})
//		}
//		return c.JSON(hasil)
//	}

func GetMahasiswaBebasPustaka(c *fiber.Ctx) error {
	filterMember := c.Query("member_id")
	hasil, err := services.GetMahasiswaBebasPustakaServiceFast(filterMember) // service Fast
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": err.Error()})
	}
	return c.JSON(hasil)
}

func DeleteBebasPustaka(c *fiber.Ctx) error {
	idParam := c.Params("id")
	if idParam == "" {
		return c.Status(400).JSON(fiber.Map{"status": "error", "message": "missing id parameter"})
	}

	var idInt int
	_, err := fmt.Sscanf(idParam, "%d", &idInt)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"status": "error", "message": "invalid id parameter"})
	}

	approval, err := services.FindApprovalByID(uint(idInt))
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"status": "error", "message": "Data Bebas Pustaka tidak ditemukan"})
	}
	if err := services.DeleteApproval(approval); err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": err.Error()})
	}

	return c.JSON(fiber.Map{"status": "success", "message": "Data Bebas Pustaka berhasil dihapus"})
}

func DeleteAllLoanByMember(c *fiber.Ctx) error {
	memberID := c.Params("member_id")
	if memberID == "" {
		return c.Status(400).JSON(fiber.Map{"status": "error", "message": "missing member_id parameter"})
	}

	loanURL := fmt.Sprintf("http://localhost:8080/loan?member_id=%s", memberID)
	resp, err := http.Get(loanURL)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "Gagal ambil data loan dari Opac"})
	}
	defer resp.Body.Close()

	var loanPayload map[string]interface{}
	bodyBytes, _ := io.ReadAll(resp.Body)
	if err := json.Unmarshal(bodyBytes, &loanPayload); err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "Gagal decode data loan: " + err.Error()})
	}

	loansData, ok := loanPayload["data"].([]interface{})
	if !ok {
		return c.Status(400).JSON(fiber.Map{"status": "error", "message": "Data loan tidak valid"})
	}

	deleted := 0
	for _, loanItem := range loansData {
		loanMap, ok := loanItem.(map[string]interface{})
		if !ok {
			continue
		}
		loanIdFloat, ok := loanMap["loan_id"].(float64)
		if !ok {
			continue
		}
		loanId := int(loanIdFloat)
		delURL := fmt.Sprintf("http://localhost:8080/loan/%d", loanId)
		req, _ := http.NewRequest("DELETE", delURL, nil)
		delResp, err := http.DefaultClient.Do(req)
		if err == nil && delResp.StatusCode == 200 {
			deleted++
		}
		if delResp != nil {
			delResp.Body.Close()
		}
	}

	return c.JSON(fiber.Map{
		"status":  "success",
		"message": fmt.Sprintf("%d loan berhasil dihapus untuk member_id %s", deleted, memberID),
		"deleted": deleted,
	})
}

func GetBebasPustakaHistory(c *fiber.Ctx) error {
	kodeUser := c.Params("kode_user")
	if kodeUser == "" {
		return c.Status(400).JSON(fiber.Map{"status": "error", "message": "missing kode_user parameter"})
	}

	// Ambil data master mahasiswa dari Sikompen
	sikompenURL := fmt.Sprintf("http://localhost:8000/api/mahasiswa?kode_user=%s", kodeUser)
	respMahasiswa, err := http.Get(sikompenURL)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "Gagal ambil data mahasiswa Sikompen"})
	}
	defer respMahasiswa.Body.Close()

	var mahasiswaList []map[string]interface{}
	bodyBytes, _ := io.ReadAll(respMahasiswa.Body)
	if err := json.Unmarshal(bodyBytes, &mahasiswaList); err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "Gagal decode mahasiswa: " + err.Error()})
	}

	var m map[string]interface{}
	found := false
	for _, mm := range mahasiswaList {
		if ku, ok := mm["kode_user"].(string); ok && ku == kodeUser {
			m = mm
			found = true
			break
		}
	}
	if !found {
		return c.Status(404).JSON(fiber.Map{"status": "error", "message": "Mahasiswa tidak ditemukan"})
	}

	loanURL := fmt.Sprintf("http://localhost:8080/loan?member_id=%s", kodeUser)
	respLoan, err := http.Get(loanURL)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "Gagal ambil data loan dari Opac"})
	}
	defer respLoan.Body.Close()

	var loanPayload map[string]interface{}
	bodyBytes, _ = io.ReadAll(respLoan.Body)
	if err := json.Unmarshal(bodyBytes, &loanPayload); err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "Gagal decode loan: " + err.Error()})
	}

	loansData, ok := loanPayload["data"].([]interface{})
	if ok && len(loansData) > 0 {
		var maxTime string
		for _, loanItem := range loansData {
			loanMap, ok := loanItem.(map[string]interface{})
			if !ok {
				continue
			}
			lu, ok := loanMap["last_update"].(string)
			if ok && (maxTime == "" || lu > maxTime) {
				maxTime = lu
			}
		}
	}

	approval, errApproval := services.GetLatestApprovalByKodeUser(kodeUser)
	var updatedBy interface{} = nil
	var waktuStatus interface{} = nil
	var namaUserUpdate interface{} = nil
	if errApproval == nil && approval.IDUsers != 0 {
		updatedBy = approval.IDUsers
		waktuStatus = approval.CreatedAt
		user, errUser := services.GetUserByID(approval.IDUsers)
		if errUser == nil {
			namaUserUpdate = user.Name
		}
	}

	return c.JSON(fiber.Map{
		"status":            "success",
		"kode_user":         m["kode_user"],
		"nama_user":         m["nama_user"],
		"kelas":             m["kelas"],
		"updated_by":        updatedBy,
		"updated_by_name":   namaUserUpdate,
		"created_at_status": waktuStatus,
	})
}

func DeleteLoanById(c *fiber.Ctx) error {
	loanID := c.Params("loan_id")
	if loanID == "" {
		return c.Status(400).JSON(fiber.Map{"status": "error", "message": "missing loan_id parameter"})
	}

	delURL := fmt.Sprintf("http://localhost:8080/loan/%s", loanID)
	reqDel, _ := http.NewRequest("DELETE", delURL, nil)
	resp, err := http.DefaultClient.Do(reqDel)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "Gagal request ke Opac"})
	}
	defer resp.Body.Close()

	if resp.StatusCode == 200 {
		return c.JSON(fiber.Map{"status": "success", "message": "Loan berhasil dihapus", "loan_id": loanID})
	}
	return c.Status(resp.StatusCode).JSON(fiber.Map{"status": "error", "message": "Gagal hapus loan di Opac", "code": resp.StatusCode})
}
