package controllers

import (
	"encoding/json"
	"io"
	"net/http"
	"time"

	"fmt"

	"github.com/gofiber/fiber/v2"
)

func GetDashboardStats(c *fiber.Ctx) error {
	client := &http.Client{Timeout: 8 * time.Second}
	mahasiswaURL := "http://localhost:8000/api/mahasiswa?limit=0"
	resp, err := client.Get(mahasiswaURL)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Gagal ambil data mahasiswa Sikompen"})
	}
	defer resp.Body.Close()

	var mahasiswaList []map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&mahasiswaList); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Gagal decode mahasiswa"})
	}

	// Build master mahasiswa map (kode_user)
	masterMhs := make(map[string]map[string]interface{})
	for _, m := range mahasiswaList {
		if kode, ok := m["kode_user"].(string); ok && kode != "" {
			masterMhs[kode] = m
		}
	}

	// Fetch all loans
	loansURL := "http://localhost:8080/loan?page=1&per_page=100000"
	respLoan, err := client.Get(loansURL)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Gagal ambil data loan"})
	}
	defer respLoan.Body.Close()
	var loanPayload map[string]interface{}
	bodyBytes, _ := io.ReadAll(respLoan.Body)
	if err := json.Unmarshal(bodyBytes, &loanPayload); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Gagal decode data loan"})
	}
	loansData, ok := loanPayload["data"].([]interface{})
	if !ok {
		loansData = []interface{}{}
	}

	// Filter loans: hanya yang member_id ada di master mahasiswa
	filteredLoans := make([]map[string]interface{}, 0)
	for _, l := range loansData {
		loanMap, ok := l.(map[string]interface{})
		if !ok {
			continue
		}
		memberID := ""
		switch v := loanMap["member_id"].(type) {
		case string:
			memberID = v
		case float64:
			memberID = fmt.Sprintf("%.0f", v)
		}
		if memberID == "" {
			continue
		}
		if _, ok := masterMhs[memberID]; ok {
			filteredLoans = append(filteredLoans, loanMap)
		}
	}
	fmt.Println("filteredLoans count:", len(filteredLoans))
	if len(filteredLoans) > 0 {
		fmt.Println("filteredLoans sample:", filteredLoans[0])
	}

	// total_mahasiswa: dari master
	totalMahasiswa := len(masterMhs)

	// Sederhana: count is_return==1 dan is_return==0 dari filteredLoans
	totalBebasPustaka := 0
	totalTunggakan := 0
	peminjamSet := make(map[string]struct{})
	for _, loan := range filteredLoans {
		memberID := ""
		switch v := loan["member_id"].(type) {
		case string:
			memberID = v
		case float64:
			memberID = fmt.Sprintf("%.0f", v)
		}
		if memberID == "" {
			continue
		}
		peminjamSet[memberID] = struct{}{}
		isReturn := -1
		switch v := loan["is_return"].(type) {
		case int:
			isReturn = v
		case int8:
			isReturn = int(v)
		case int16:
			isReturn = int(v)
		case int32:
			isReturn = int(v)
		case int64:
			isReturn = int(v)
		case float64:
			isReturn = int(v)
		case float32:
			isReturn = int(v)
		case string:
			if v == "1" {
				isReturn = 1
			} else if v == "0" {
				isReturn = 0
			}
		case bool:
			if v {
				isReturn = 1
			} else {
				isReturn = 0
			}
		default:
			continue
		}
		if isReturn == 1 {
			totalBebasPustaka++
		}
		if isReturn == 0 {
			totalTunggakan++
		}
	}
	fmt.Println("totalMahasiswa:", totalMahasiswa)
	fmt.Println("totalBebasPustaka:", totalBebasPustaka)
	fmt.Println("totalTunggakan:", totalTunggakan)
	fmt.Println("totalPeminjam:", len(peminjamSet))

	return c.JSON(fiber.Map{
		"total_mahasiswa":     totalMahasiswa,
		"total_peminjam":      len(peminjamSet),
		"total_tunggakan":     totalTunggakan,
		"total_bebas_pustaka": totalBebasPustaka,
	})
}

// Mapping kode jurusan ke nama jurusan
var kodeJurusanMap = map[string]string{
	"07": "Teknik Informatika dan Komputer",
	"06": "Teknik Grafika dan Penerbitan",
	"05": "Administrasi Niaga",
	"04": "Akuntansi",
	"03": "Teknik Elektro",
	"02": "Teknik Mesin",
	"01": "Teknik Sipil",
}

// Endpoint: /dashboard/bebas-pustaka-jurusan
func GetBebasPustakaJurusan(c *fiber.Ctx) error {
	client := &http.Client{Timeout: 8 * time.Second}
	mahasiswaURL := "http://localhost:8000/api/mahasiswa?limit=0"
	resp, err := client.Get(mahasiswaURL)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Gagal ambil data mahasiswa Sikompen"})
	}
	defer resp.Body.Close()

	var mahasiswaList []map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&mahasiswaList); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Gagal decode mahasiswa"})
	}

	// Build master mahasiswa map (kode_user)
	masterMhs := make(map[string]map[string]interface{})
	for _, m := range mahasiswaList {
		if kode, ok := m["kode_user"].(string); ok && kode != "" {
			masterMhs[kode] = m
		}
	}

	// Fetch all loans
	loansURL := "http://localhost:8080/loan?page=1&per_page=100000"
	respLoan, err := client.Get(loansURL)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Gagal ambil data loan"})
	}
	defer respLoan.Body.Close()
	var loanPayload map[string]interface{}
	bodyBytes, _ := io.ReadAll(respLoan.Body)
	if err := json.Unmarshal(bodyBytes, &loanPayload); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Gagal decode data loan"})
	}
	loansData, ok := loanPayload["data"].([]interface{})
	if !ok {
		loansData = []interface{}{}
	}

	// Filter loans: hanya yang member_id ada di master mahasiswa dan is_return==1
	jurusanCount := make(map[string]int)
	for _, l := range loansData {
		loanMap, ok := l.(map[string]interface{})
		if !ok {
			continue
		}
		memberID := ""
		switch v := loanMap["member_id"].(type) {
		case string:
			memberID = v
		case float64:
			memberID = fmt.Sprintf("%.0f", v)
		}
		if memberID == "" {
			continue
		}
		mhs, ok := masterMhs[memberID]
		if !ok {
			continue
		}
		// Cek is_return==1
		isReturn := -1
		switch v := loanMap["is_return"].(type) {
		case int:
			isReturn = v
		case int8:
			isReturn = int(v)
		case int16:
			isReturn = int(v)
		case int32:
			isReturn = int(v)
		case int64:
			isReturn = int(v)
		case float64:
			isReturn = int(v)
		case float32:
			isReturn = int(v)
		case string:
			if v == "1" {
				isReturn = 1
			} else if v == "0" {
				isReturn = 0
			}
		case bool:
			if v {
				isReturn = 1
			} else {
				isReturn = 0
			}
		default:
			continue
		}
		if isReturn != 1 {
			continue
		}
		// Ambil kode jurusan dari kode_user (2 digit ke-3 dan ke-4)
		kodeUser, _ := mhs["kode_user"].(string)
		if len(kodeUser) < 4 {
			continue
		}
		kodeJurusan := kodeUser[2:4]
		jurusanCount[kodeJurusan]++
	}

	// Build hasil
	result := make([]fiber.Map, 0)
	for kode, nama := range kodeJurusanMap {
		result = append(result, fiber.Map{
			"kode":    kode,
			"jurusan": nama,
			"total":   jurusanCount[kode],
		})
	}

	return c.JSON(result)
}

// Endpoint: /api/peminjam-per-jurusan
func GetPeminjamPerJurusan(c *fiber.Ctx) error {
	client := &http.Client{Timeout: 8 * time.Second}
	mahasiswaURL := "http://localhost:8000/api/mahasiswa?limit=0"
	resp, err := client.Get(mahasiswaURL)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Gagal ambil data mahasiswa Sikompen"})
	}
	defer resp.Body.Close()

	var mahasiswaList []map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&mahasiswaList); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Gagal decode mahasiswa"})
	}

	// Build master mahasiswa map (kode_user)
	masterMhs := make(map[string]map[string]interface{})
	for _, m := range mahasiswaList {
		if kode, ok := m["kode_user"].(string); ok && kode != "" {
			masterMhs[kode] = m
		}
	}

	// Fetch all loans
	loansURL := "http://localhost:8080/loan?page=1&per_page=100000"
	respLoan, err := client.Get(loansURL)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Gagal ambil data loan"})
	}
	defer respLoan.Body.Close()
	var loanPayload map[string]interface{}
	bodyBytes, _ := io.ReadAll(respLoan.Body)
	if err := json.Unmarshal(bodyBytes, &loanPayload); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Gagal decode data loan"})
	}
	loansData, ok := loanPayload["data"].([]interface{})
	if !ok {
		loansData = []interface{}{}
	}

	// Set member_id unik per jurusan
	jurusanPeminjam := make(map[string]map[string]struct{})
	for _, l := range loansData {
		loanMap, ok := l.(map[string]interface{})
		if !ok {
			continue
		}
		memberID := ""
		switch v := loanMap["member_id"].(type) {
		case string:
			memberID = v
		case float64:
			memberID = fmt.Sprintf("%.0f", v)
		}
		if memberID == "" {
			continue
		}
		mhs, ok := masterMhs[memberID]
		if !ok {
			continue
		}
		kodeUser, _ := mhs["kode_user"].(string)
		if len(kodeUser) < 4 {
			continue
		}
		kodeJurusan := kodeUser[2:4]
		if jurusanPeminjam[kodeJurusan] == nil {
			jurusanPeminjam[kodeJurusan] = make(map[string]struct{})
		}
		jurusanPeminjam[kodeJurusan][memberID] = struct{}{}
	}

	// Build hasil
	result := make([]fiber.Map, 0)
	for kode, nama := range kodeJurusanMap {
		total := 0
		if set, ok := jurusanPeminjam[kode]; ok {
			total = len(set)
		}
		result = append(result, fiber.Map{
			"kode":    kode,
			"jurusan": nama,
			"total":   total,
		})
	}

	return c.JSON(result)
}
