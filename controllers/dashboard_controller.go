package controllers

import (
	"encoding/json"
	"io"
	"net/http"
	"sync"
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

	jurusanKodeNama := map[string]string{
		"07": "Teknik Informatika dan Komputer",
		"06": "Teknik Grafika dan Penerbitan",
		"05": "Administrasi Niaga",
		"04": "Akuntansi",
		"03": "Teknik Elektro",
		"02": "Teknik Mesin",
		"01": "Teknik Sipil",
	}
	jurusanBebasMap := make(map[string]int)

	totalMahasiswa := len(mahasiswaList)
	bebasCount := 0
	tunggakanCount := 0
	peminjamSet := make(map[string]struct{})

	var wg sync.WaitGroup
	var mu sync.Mutex

	for _, m := range mahasiswaList {
		wg.Add(1)
		go func(m map[string]interface{}) {
			defer wg.Done()
			kodeUser, ok := m["kode_user"].(string)
			if !ok {
				return
			}
			loanURL := fmt.Sprintf("http://localhost:8080/loan?member_id=%s", kodeUser)
			loanResp, err := client.Get(loanURL)
			var loansData []map[string]interface{}
			if err == nil {
				defer loanResp.Body.Close()
				bodyBytes, _ := io.ReadAll(loanResp.Body)
				// decode ke field
				var loanWrapper struct {
					Data []map[string]interface{} `json:"data"`
				}
				if err := json.Unmarshal(bodyBytes, &loanWrapper); err == nil && loanWrapper.Data != nil {
					loansData = loanWrapper.Data
				} else {
					// decode langsung ke array
					_ = json.Unmarshal(bodyBytes, &loansData)
				}
			}
			adaTunggakan := false
			totalLoan := len(loansData)
			if totalLoan > 0 {
				mu.Lock()
				peminjamSet[kodeUser] = struct{}{}
				mu.Unlock()
			}
			loanReturned := 0
			now := time.Now()
			for _, loanMap := range loansData {
				// Loan dianggap sudah dikembalikan jika ada return_date (string tidak kosong) atau status == "Lunas"
				returned := false
				if returnDate, ok := loanMap["return_date"].(string); ok && returnDate != "" {
					returned = true
				} else if status, ok := loanMap["status"].(string); ok && status == "Lunas" {
					returned = true
				}
				if !returned {
					// Cek tunggakan: due_date < hari ini
					dueDateStr, ok := loanMap["due_date"].(string)
					if ok {
						dueDate, err := time.Parse("2006-01-02", dueDateStr)
						if err == nil && dueDate.Before(now) {
							adaTunggakan = true
						}
					}
				} else {
					loanReturned++
				}
			}
			mu.Lock()
			if totalLoan == 0 || loanReturned == totalLoan {
				bebasCount++
				jurusanNama := "Unknown"
				if len(kodeUser) >= 4 {
					kodeJurusan := kodeUser[2:4]
					if nama, ok := jurusanKodeNama[kodeJurusan]; ok {
						jurusanNama = nama
					}
				}
				jurusanBebasMap[jurusanNama]++
			}
			// Hapus increment peminjamCount, diganti dengan peminjamSet
			if adaTunggakan {
				tunggakanCount++
			}
			mu.Unlock()
		}(m)
	}
	wg.Wait()

	// Pastikan semua jurusan mapping muncul walaupun 0
	for _, nama := range jurusanKodeNama {
		if _, ok := jurusanBebasMap[nama]; !ok {
			jurusanBebasMap[nama] = 0
		}
	}

	return c.JSON(fiber.Map{
		"total_mahasiswa":            totalMahasiswa,
		"total_peminjam":             len(peminjamSet),
		"total_tunggakan":            tunggakanCount,
		"total_bebas_pustaka":        bebasCount,
		"chart_bebas_kompen_jurusan": jurusanBebasMap,
	})
}
