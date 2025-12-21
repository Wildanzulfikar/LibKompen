package services

import (
	"LibKompen/database"
	"LibKompen/models"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"
)

// =======================================
// SERVICE BEBAS PUSTAKA (FINAL VERSION)
// =======================================
// func GetMahasiswaBebasPustakaServiceFast(memberID string) ([]map[string]interface{}, error) {
// 	client := &http.Client{Timeout: 10 * time.Second}

// 	// 1. Ambil semua mahasiswa
// 	resp, err := client.Get("http://localhost:8000/api/mahasiswa?limit=0")
// 	if err != nil {
// 		return nil, fmt.Errorf("gagal ambil mahasiswa: %v", err)
// 	}
// 	defer resp.Body.Close()

// 	var mahasiswaList []map[string]interface{}
// 	if err := json.NewDecoder(resp.Body).Decode(&mahasiswaList); err != nil {
// 		return nil, fmt.Errorf("gagal decode mahasiswa: %v", err)
// 	}

// 	var (
// 		wg    sync.WaitGroup
// 		mu    sync.Mutex
// 		hasil = make([]map[string]interface{}, 0)
// 	)

// 	for _, m := range mahasiswaList {
// 		wg.Add(1)

// 		go func(m map[string]interface{}) {
// 			defer wg.Done()

// 			defer func() { // catch panic supaya goroutine tidak mati
// 				if r := recover(); r != nil {
// 					fmt.Println("panic:", r)
// 				}
// 			}()

// 			nim, ok := m["kode_user"].(string)
// 			if !ok || nim == "" {
// 				return
// 			}
// 			if memberID != "" && nim != memberID {
// 				return
// 			}

// 			// =========================
// 			// DEFAULT: Bebas Pustaka
// 			// =========================
// 			response := map[string]interface{}{
// 				"id_mahasiswa":    m["id_mahasiswa"],
// 				"nim":             nim,
// 				"nama":            m["nama_user"],
// 				"prodi":           m["prodi"],
// 				"kelas":           m["kelas"],
// 				"semester":        m["semester"],
// 				"status":          "Bebas Pustaka",
// 				"status_pinjaman": "Lunas",
// 				"keterangan":      "-",
// 			}

// 			// =========================
// 			// AMBIL DATA PINJAMAN DETAIL REALTIME
// 			// =========================
// 			loanURL := fmt.Sprintf("http://localhost:3000/api/loan/member/%s", nim)
// 			loanResp, err := client.Get(loanURL)
// 			if err == nil && loanResp.StatusCode == 200 {
// 				defer loanResp.Body.Close()

// 				var loanData struct {
// 					PinjamanAktif []interface{} `json:"pinjaman_aktif"`
// 				}

// 				if err := json.NewDecoder(loanResp.Body).Decode(&loanData); err == nil {
// 					if len(loanData.PinjamanAktif) > 0 {
// 						response["status"] = "Tanggungan"
// 						response["status_pinjaman"] = "Belum"
// 					}
// 				} else {
// 					fmt.Println("warning: gagal decode loan member untuk", nim)
// 				}
// 			} else if err != nil {
// 				fmt.Println("warning: gagal ambil loan member untuk", nim, "error:", err)
// 			}

// 			// =========================
// 			// KETERANGAN (JIKA TANGGUNGAN)
// 			// =========================
// 			if response["status"] == "Tanggungan" {
// 				var approval models.ApprovalBebasPustaka
// 				if err := database.DB.
// 					Where("kode_user = ?", nim).
// 					First(&approval).Error; err == nil {
// 					response["keterangan"] = approval.Keterangan
// 				}
// 			}

// 			mu.Lock()
// 			hasil = append(hasil, response)
// 			mu.Unlock()
// 		}(m)
// 	}

// 	wg.Wait()
// 	return hasil, nil
// }

func GetMahasiswaBebasPustakaServiceFast(memberID string) ([]map[string]interface{}, error) {
	client := &http.Client{Timeout: 10 * time.Second}

	// 1. Ambil semua mahasiswa
	resp, err := client.Get("http://localhost:8000/api/mahasiswa?limit=0")
	if err != nil {
		return nil, fmt.Errorf("gagal ambil mahasiswa: %v", err)
	}
	defer resp.Body.Close()

	var mahasiswaList []map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&mahasiswaList); err != nil {
		return nil, fmt.Errorf("gagal decode mahasiswa: %v", err)
	}

	var (
		hasil   = make([]map[string]interface{}, 0)
		mu      sync.Mutex
		wg      sync.WaitGroup
		jobs    = make(chan map[string]interface{}, len(mahasiswaList))
		workers = 10 
	)

	// Worker function
	worker := func() {
		for m := range jobs {
			nim, ok := m["kode_user"].(string)
			if !ok || nim == "" {
				wg.Done()
				continue
			}
			if memberID != "" && nim != memberID {
				wg.Done()
				continue
			}

			response := map[string]interface{}{
				"id_mahasiswa":    m["id_mahasiswa"],
				"nim":             nim,
				"nama":            m["nama_user"],
				"prodi":           m["prodi"],
				"kelas":           m["kelas"],
				"semester":        m["semester"],
				"status":          "Bebas Pustaka",
				"status_pinjaman": "Lunas",
				"keterangan":      "-",
			}

			// Ambil data pinjaman realtime
			loanURL := fmt.Sprintf("http://localhost:3000/api/loan/member/%s", nim)
			loanResp, err := client.Get(loanURL)
			if err == nil && loanResp.StatusCode == 200 {
				defer loanResp.Body.Close()

				var loanData struct {
					PinjamanAktif []interface{} `json:"pinjaman_aktif"`
				}
				if err := json.NewDecoder(loanResp.Body).Decode(&loanData); err == nil {
					if len(loanData.PinjamanAktif) > 0 {
						response["status"] = "Tanggungan"
						response["status_pinjaman"] = "Belum"
					}
				}
			}

			// Ambil keterangan approval jika Tanggungan
			if response["status"] == "Tanggungan" {
				var approval models.ApprovalBebasPustaka
				if err := database.DB.
					Where("kode_user = ?", nim).
					First(&approval).Error; err == nil {
					response["keterangan"] = approval.Keterangan
				}
			}

			mu.Lock()
			hasil = append(hasil, response)
			mu.Unlock()
			wg.Done()
		}
	}

	// Start workers
	for i := 0; i < workers; i++ {
		go worker()
	}

	// Kirim job ke channel
	for _, m := range mahasiswaList {
		wg.Add(1)
		jobs <- m
	}
	close(jobs)

	wg.Wait()
	return hasil, nil
}

var loanCache = struct {
	sync.RWMutex
	data map[string]map[string]interface{}
	time map[string]time.Time
}{
	data: make(map[string]map[string]interface{}),
	time: make(map[string]time.Time),
}

// // Fungsi utama
// func GetMahasiswaBebasPustakaServiceFast(memberID string) ([]map[string]interface{}, error) {
// 	client := &http.Client{Timeout: 5 * time.Second}

// 	// 1. Ambil semua mahasiswa dari Sikompen
// 	resp, err := client.Get("http://localhost:8000/api/mahasiswa?limit=0")
// 	if err != nil {
// 		return nil, fmt.Errorf("gagal ambil mahasiswa: %v", err)
// 	}
// 	defer resp.Body.Close()

// 	var mahasiswaList []map[string]interface{}
// 	if err := json.NewDecoder(resp.Body).Decode(&mahasiswaList); err != nil {
// 		return nil, fmt.Errorf("gagal decode mahasiswa: %v", err)
// 	}

// 	var (
// 		hasil   = make([]map[string]interface{}, 0)
// 		mu      sync.Mutex
// 		wg      sync.WaitGroup
// 		jobs    = make(chan map[string]interface{}, len(mahasiswaList))
// 		workers = 20
// 	)

// 	worker := func() {
// 		for m := range jobs {
// 			nim, ok := m["kode_user"].(string)
// 			if !ok || nim == "" {
// 				wg.Done()
// 				continue
// 			}
// 			if memberID != "" && nim != memberID {
// 				wg.Done()
// 				continue
// 			}

// 			// Cek cache dulu
// 			loanCache.RLock()
// 			cached, exists := loanCache.data[nim]
// 			lastFetch := loanCache.time[nim]
// 			loanCache.RUnlock()

// 			response := map[string]interface{}{
// 				"id_mahasiswa":    m["id_mahasiswa"],
// 				"nim":             nim,
// 				"nama":            m["nama_user"],
// 				"prodi":           m["prodi"],
// 				"kelas":           m["kelas"],
// 				"semester":        m["semester"],
// 				"status":          "Bebas Pustaka",
// 				"status_pinjaman": "Lunas",
// 				"keterangan":      "-",
// 			}

// 			if exists && time.Since(lastFetch) < 2*time.Minute {
// 				// pakai cache
// 				if status, ok := cached["status"].(string); ok {
// 					response["status"] = status
// 				}
// 				if sp, ok := cached["status_pinjaman"].(string); ok {
// 					response["status_pinjaman"] = sp
// 				}
// 				if ket, ok := cached["keterangan"].(string); ok {
// 					response["keterangan"] = ket
// 				}
// 			} else {
// 				// ambil realtime
// 				loanURL := fmt.Sprintf("http://localhost:3000/api/loan/member/%s", nim)
// 				loanResp, err := client.Get(loanURL)
// 				if err != nil || loanResp.StatusCode != 200 {
// 					response["status"] = "Tanggungan"
// 					response["status_pinjaman"] = "Belum"
// 				} else {
// 					defer loanResp.Body.Close()
// 					var loanData struct {
// 						PinjamanAktif []interface{} `json:"pinjaman_aktif"`
// 					}
// 					if err := json.NewDecoder(loanResp.Body).Decode(&loanData); err != nil {
// 						response["status"] = "Tanggungan"
// 						response["status_pinjaman"] = "Belum"
// 					} else if len(loanData.PinjamanAktif) > 0 {
// 						response["status"] = "Tanggungan"
// 						response["status_pinjaman"] = "Belum"
// 					}
// 				}

// 				// ambil keterangan approval
// 				if response["status"] == "Tanggungan" {
// 					var approval models.ApprovalBebasPustaka
// 					if err := database.DB.Where("kode_user = ?", nim).First(&approval).Error; err == nil {
// 						response["keterangan"] = approval.Keterangan
// 					}
// 				}

// 				// simpan ke cache
// 				loanCache.Lock()
// 				loanCache.data[nim] = map[string]interface{}{
// 					"status":          response["status"],
// 					"status_pinjaman": response["status_pinjaman"],
// 					"keterangan":      response["keterangan"],
// 				}
// 				loanCache.time[nim] = time.Now()
// 				loanCache.Unlock()
// 			}

// 			mu.Lock()
// 			hasil = append(hasil, response)
// 			mu.Unlock()
// 			wg.Done()
// 		}
// 	}

// 	// start workers
// 	for i := 0; i < workers; i++ {
// 		go worker()
// 	}

// 	// kirim job ke channel
// 	for _, m := range mahasiswaList {
// 		wg.Add(1)
// 		jobs <- m
// 	}
// 	close(jobs)
// 	wg.Wait()

// 	return hasil, nil
// }