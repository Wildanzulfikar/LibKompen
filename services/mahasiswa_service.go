package services

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"
)

// func GetMahasiswaBebasPustakaService(memberID string) ([]map[string]interface{}, error) {
// 	sikompenURL := "http://localhost:8000/api/mahasiswa?limit=0"
// 	respMahasiswa, err := http.Get(sikompenURL)
// 	if err != nil {
// 		return nil, fmt.Errorf("Gagal ambil data mahasiswa Sikompen: %w", err)
// 	}
// 	defer respMahasiswa.Body.Close()

// 	var mahasiswaList []map[string]interface{}
// 	bodyBytes, _ := io.ReadAll(respMahasiswa.Body)
// 	if err := json.Unmarshal(bodyBytes, &mahasiswaList); err != nil {
// 		return nil, fmt.Errorf("Gagal decode mahasiswa: %w", err)
// 	}

// 	var hasil []map[string]interface{}
// 	for _, m := range mahasiswaList {
// 		kodeUser, ok := m["kode_user"].(string)
// 		if !ok {
// 			continue
// 		}
// 		if memberID != "" && kodeUser != memberID {
// 			continue
// 		}
// 		namaUser, _ := m["nama_user"].(string)
// 		prodi, _ := m["prodi"].(string)
// 		kelas, _ := m["kelas"].(string)
// 		semester, _ := m["semester"].(string)
// 		idMahasiswa := m["id_mahasiswa"]

// 		response := map[string]interface{}{
// 			"id_mahasiswa":    idMahasiswa,
// 			"nim":             kodeUser,
// 			"nama":            namaUser,
// 			"prodi":           prodi,
// 			"kelas":           kelas,
// 			"semester":        semester,
// 			"status":          "Bebas Pustaka",
// 			"status_pinjaman": "Lunas",
// 			"keterangan":      "-",
// 		}

// 		loanURL := fmt.Sprintf("http://localhost:8080/loan?member_id=%s", kodeUser)
// 		loanResp, err := http.Get(loanURL)
// 		if err != nil {
// 			hasil = append(hasil, response)
// 			continue
// 		}
// 		defer loanResp.Body.Close()

// 		bodyBytes, _ := io.ReadAll(loanResp.Body)
// 		var loanPayload map[string]interface{}
// 		if err := json.Unmarshal(bodyBytes, &loanPayload); err != nil {
// 			hasil = append(hasil, response)
// 			continue
// 		}

// 		loansData, ok := loanPayload["data"].([]interface{})
// 		if !ok {
// 			hasil = append(hasil, response)
// 			continue
// 		}

// 		adaTanggungan := false
// 		for _, loanItem := range loansData {
// 			loanMap, ok := loanItem.(map[string]interface{})
// 			if !ok {
// 				continue
// 			}
// 			isReturn, ok := loanMap["is_return"].(bool)
// 			if ok && !isReturn {
// 				adaTanggungan = true
// 				break
// 			}
// 		}

// 		if adaTanggungan {
// 			response["status"] = "Tanggungan"
// 			response["status_pinjaman"] = "Belum"

// 			statusURL := fmt.Sprintf("http://localhost:8000/api/status_approval?kode_user=%s", kodeUser)
// 			statusResp, err := http.Get(statusURL)
// 			if err == nil {
// 				defer statusResp.Body.Close()
// 				bodyBytes, _ := io.ReadAll(statusResp.Body)
// 				var statApproval []map[string]interface{}
// 				if err := json.Unmarshal(bodyBytes, &statApproval); err == nil {
// 					if len(statApproval) > 0 {
// 						keterangan, ok := statApproval[0]["keterangan"].(string)
// 						if ok {
// 							response["keterangan"] = keterangan
// 						}
// 					}
// 				}
// 			}
// 		}

// 		hasil = append(hasil, response)
// 	}
// 	return hasil, nil
// }

// func GetMahasiswaBebasPustakaService(memberID string) ([]map[string]interface{}, error) {
// 	sikompenURL := "http://localhost:8000/api/mahasiswa?limit=0"
// 	respMahasiswa, err := http.Get(sikompenURL)
// 	if err != nil {
// 		return nil, fmt.Errorf("Gagal ambil data mahasiswa Sikompen: %w", err)
// 	}
// 	defer respMahasiswa.Body.Close()

// 	// ================================
// 	// FIXED BAGIAN 1: Decode API mahasiswa
// 	// ================================
// 	bodyBytes, _ := io.ReadAll(respMahasiswa.Body)

// 	// API tidak mengembalikan array, tetapi object → decode ke map
// 	var payload map[string]interface{}
// 	if err := json.Unmarshal(bodyBytes, &payload); err != nil {
// 		return nil, fmt.Errorf("Gagal decode mahasiswa: %w", err)
// 	}

// 	// Ambil data: []interface{}
// 	dataArr, ok := payload["data"].([]interface{})
// 	if !ok {
// 		return nil, fmt.Errorf("Data mahasiswa tidak berbentuk array")
// 	}

// 	// Convert ke []map[string]interface{}
// 	mahasiswaList := make([]map[string]interface{}, 0)
// 	for _, item := range dataArr {
// 		if m, ok := item.(map[string]interface{}); ok {
// 			mahasiswaList = append(mahasiswaList, m)
// 		}
// 	}
// 	// ================================
// 	// END FIXED
// 	// ================================

// 	var hasil []map[string]interface{}
// 	for _, m := range mahasiswaList {

// 		kodeUser, ok := m["kode_user"].(string)
// 		if !ok {
// 			continue
// 		}
// 		if memberID != "" && kodeUser != memberID {
// 			continue
// 		}

// 		namaUser, _ := m["nama_user"].(string)
// 		prodi, _ := m["prodi"].(string)
// 		kelas, _ := m["kelas"].(string)
// 		semester, _ := m["semester"].(string)
// 		idMahasiswa := m["id_mahasiswa"]

// 		response := map[string]interface{}{
// 			"id_mahasiswa":    idMahasiswa,
// 			"nim":             kodeUser,
// 			"nama":            namaUser,
// 			"prodi":           prodi,
// 			"kelas":           kelas,
// 			"semester":        semester,
// 			"status":          "Bebas Pustaka",
// 			"status_pinjaman": "Lunas",
// 			"keterangan":      "-",
// 		}

// 		// ================================
// 		// Loan check
// 		// ================================
// 		loanURL := fmt.Sprintf("http://localhost:8080/loan?member_id=%s", kodeUser)
// 		loanResp, err := http.Get(loanURL)
// 		if err != nil {
// 			hasil = append(hasil, response)
// 			continue
// 		}
// 		defer loanResp.Body.Close()

// 		bodyBytes, _ = io.ReadAll(loanResp.Body)
// 		var loanPayload map[string]interface{}
// 		if err := json.Unmarshal(bodyBytes, &loanPayload); err != nil {
// 			hasil = append(hasil, response)
// 			continue
// 		}

// 		loansData, ok := loanPayload["data"].([]interface{})
// 		if !ok {
// 			hasil = append(hasil, response)
// 			continue
// 		}

// 		adaTanggungan := false
// 		for _, loanItem := range loansData {
// 			loanMap, ok := loanItem.(map[string]interface{})
// 			if !ok {
// 				continue
// 			}
// 			isReturn, ok := loanMap["is_return"].(bool)
// 			if ok && !isReturn {
// 				adaTanggungan = true
// 				break
// 			}
// 		}

// 		if adaTanggungan {
// 			response["status"] = "Tanggungan"
// 			response["status_pinjaman"] = "Belum"

// 			statusURL := fmt.Sprintf("http://localhost:8000/api/status_approval?kode_user=%s", kodeUser)
// 			statusResp, err := http.Get(statusURL)
// 			if err == nil {
// 				defer statusResp.Body.Close()
// 				bodyBytes, _ := io.ReadAll(statusResp.Body)
// 				var statApproval []map[string]interface{}
// 				if err := json.Unmarshal(bodyBytes, &statApproval); err == nil {
// 					if len(statApproval) > 0 {
// 						keterangan, ok := statApproval[0]["keterangan"].(string)
// 						if ok {
// 							response["keterangan"] = keterangan
// 						}
// 					}
// 				}
// 			}
// 		}

// 		hasil = append(hasil, response)
// 	}

// 	return hasil, nil
// }

// func GetMahasiswaBebasPustakaService(memberID string) ([]map[string]interface{}, error) {
// 	sikompenURL := "http://localhost:8000/api/mahasiswa?limit=0"
// 	respMahasiswa, err := http.Get(sikompenURL)
// 	if err != nil {
// 		return nil, fmt.Errorf("Gagal ambil data mahasiswa Sikompen: %w", err)
// 	}
// 	defer respMahasiswa.Body.Close()

// 	bodyBytes, _ := io.ReadAll(respMahasiswa.Body)

// 	// =========================================
// 	// FIXED: Decode langsung ke array
// 	// =========================================
// 	var mahasiswaList []map[string]interface{}
// 	if err := json.Unmarshal(bodyBytes, &mahasiswaList); err != nil {
// 		return nil, fmt.Errorf("Gagal decode mahasiswa: %w", err)
// 	}

// 	var hasil []map[string]interface{}
// 	for _, m := range mahasiswaList {

// 		kodeUser, ok := m["kode_user"].(string)
// 		if !ok {
// 			continue
// 		}
// 		if memberID != "" && kodeUser != memberID {
// 			continue
// 		}

// 		namaUser, _ := m["nama_user"].(string)
// 		prodi, _ := m["prodi"].(string)
// 		kelas, _ := m["kelas"].(string)
// 		semester, _ := m["semester"].(string)
// 		idMahasiswa := m["id_mahasiswa"] // bisa null

// 		response := map[string]interface{}{
// 			"id_mahasiswa":    idMahasiswa,
// 			"nim":             kodeUser,
// 			"nama":            namaUser,
// 			"prodi":           prodi,
// 			"kelas":           kelas,
// 			"semester":        semester,
// 			"status":          "Bebas Pustaka",
// 			"status_pinjaman": "Lunas",
// 			"keterangan":      "-",
// 		}

// 		// ================================
// 		// Loan check
// 		// ================================
// 		loanURL := fmt.Sprintf("http://localhost:8080/loan?member_id=%s", kodeUser)
// 		loanResp, err := http.Get(loanURL)
// 		if err != nil {
// 			hasil = append(hasil, response)
// 			continue
// 		}
// 		defer loanResp.Body.Close()

// 		loanBytes, _ := io.ReadAll(loanResp.Body)
// 		var loansData []map[string]interface{}
// 		if err := json.Unmarshal(loanBytes, &loansData); err != nil {
// 			hasil = append(hasil, response)
// 			continue
// 		}

// 		adaTanggungan := false
// 		for _, loanMap := range loansData {
// 			isReturn, ok := loanMap["is_return"].(bool)
// 			if ok && !isReturn {
// 				adaTanggungan = true
// 				break
// 			}
// 		}

// 		if adaTanggungan {
// 			response["status"] = "Tanggungan"
// 			response["status_pinjaman"] = "Belum"

// 			statusURL := fmt.Sprintf("http://localhost:8000/api/status_approval?kode_user=%s", kodeUser)
// 			statusResp, err := http.Get(statusURL)
// 			if err == nil {
// 				defer statusResp.Body.Close()
// 				bodyBytes, _ := io.ReadAll(statusResp.Body)
// 				var statApproval []map[string]interface{}
// 				if err := json.Unmarshal(bodyBytes, &statApproval); err == nil {
// 					if len(statApproval) > 0 {
// 						keterangan, ok := statApproval[0]["keterangan"].(string)
// 						if ok {
// 							response["keterangan"] = keterangan
// 						}
// 					}
// 				}
// 			}
// 		}

// 		hasil = append(hasil, response)
// 	}

// 	return hasil, nil
// }

// func GetMahasiswaBebasPustakaServiceFast() ([]map[string]interface{}, error) {
// 	client := &http.Client{Timeout: 8 * time.Second}

// 	// 1. Ambil semua mahasiswa
// 	mahasiswaURL := "http://localhost:8000/api/mahasiswa?limit=0"
// 	resp, err := client.Get(mahasiswaURL)
// 	if err != nil {
// 		return nil, fmt.Errorf("gagal ambil mahasiswa: %v", err)
// 	}
// 	defer resp.Body.Close()

// 	var mahasiswaList []map[string]interface{}
// 	if err := json.NewDecoder(resp.Body).Decode(&mahasiswaList); err != nil {
// 		return nil, fmt.Errorf("gagal decode mahasiswa: %v", err)
// 	}

// 	// 2. Slice hasil dan mutex
// 	hasil := make([]map[string]interface{}, 0, len(mahasiswaList))
// 	var wg sync.WaitGroup
// 	var mu sync.Mutex

// 	// 3. Loop mahasiswa, ambil loan & status paralel
// 	for _, m := range mahasiswaList {
// 		wg.Add(1)
// 		go func(m map[string]interface{}) {
// 			defer wg.Done()

// 			kodeUser, ok := m["kode_user"].(string)
// 			if !ok {
// 				return
// 			}

// 			response := map[string]interface{}{
// 				"id_mahasiswa":    m["id_mahasiswa"],
// 				"nim":             kodeUser,
// 				"nama":            m["nama_user"],
// 				"prodi":           m["prodi"],
// 				"kelas":           m["kelas"],
// 				"semester":        m["semester"],
// 				"status":          "Bebas Pustaka",
// 				"status_pinjaman": "Lunas",
// 				"keterangan":      "-",
// 			}

// 			// Ambil loan
// 			loanURL := fmt.Sprintf("http://localhost:8080/loan?member_id=%s", kodeUser)
// 			loanResp, err := client.Get(loanURL)
// 			var loansData []map[string]interface{}
// 			if err == nil {
// 				defer loanResp.Body.Close()
// 				var loanWrapper struct {
// 					Data []map[string]interface{} `json:"data"`
// 				}
// 				if err := json.NewDecoder(loanResp.Body).Decode(&loanWrapper); err == nil {
// 					loansData = loanWrapper.Data
// 				}
// 			}

// 			// Cek tanggungan
// 			adaTanggungan := false
// 			for _, loanMap := range loansData {
// 				if isReturn, ok := loanMap["is_return"].(bool); ok && !isReturn {
// 					adaTanggungan = true
// 					break
// 				}
// 			}

// 			if adaTanggungan {
// 				response["status"] = "Tanggungan"
// 				response["status_pinjaman"] = "Belum"

// 				// Ambil status approval
// 				statusURL := fmt.Sprintf("http://localhost:8000/api/status_approval?kode_user=%s", kodeUser)
// 				statusResp, err := client.Get(statusURL)
// 				if err == nil {
// 					defer statusResp.Body.Close()
// 					var statApproval []map[string]interface{}
// 					if err := json.NewDecoder(statusResp.Body).Decode(&statApproval); err == nil {
// 						if len(statApproval) > 0 {
// 							if keterangan, ok := statApproval[0]["keterangan"].(string); ok {
// 								response["keterangan"] = keterangan
// 							}
// 						}
// 					}
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
	client := &http.Client{Timeout: 8 * time.Second}

	// Ambil semua mahasiswa
	mahasiswaURL := "http://localhost:8000/api/mahasiswa?limit=0"
	resp, err := client.Get(mahasiswaURL)
	if err != nil {
		return nil, fmt.Errorf("Gagal ambil data mahasiswa Sikompen: %v", err)
	}
	defer resp.Body.Close()

	var mahasiswaList []map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&mahasiswaList); err != nil {
		return nil, fmt.Errorf("Gagal decode mahasiswa: %v", err)
	}

	hasil := make([]map[string]interface{}, 0, len(mahasiswaList))
	var wg sync.WaitGroup
	var mu sync.Mutex

	for _, m := range mahasiswaList {
		wg.Add(1)
		go func(m map[string]interface{}) {
			defer wg.Done()

			kodeUser, ok := m["kode_user"].(string)
			if !ok || (memberID != "" && kodeUser != memberID) {
				return
			}

			response := map[string]interface{}{
				"id_mahasiswa":    m["id_mahasiswa"],
				"nim":             kodeUser,
				"nama":            m["nama_user"],
				"prodi":           m["prodi"],
				"kelas":           m["kelas"],
				"semester":        m["semester"],
				"status":          "Bebas Pustaka",
				"status_pinjaman": "Lunas",
				"keterangan":      "-",
			}

			// Ambil loan
			loanURL := fmt.Sprintf("http://localhost:8080/loan?member_id=%s", kodeUser)
			loanResp, err := client.Get(loanURL)
			var loansData []map[string]interface{}
			if err == nil {
				defer loanResp.Body.Close()
				var loanWrapper struct {
					Data []map[string]interface{} `json:"data"`
				}
				if err := json.NewDecoder(loanResp.Body).Decode(&loanWrapper); err == nil {
					loansData = loanWrapper.Data
				}
			}

			// Cek tanggungan
			adaTanggungan := false
			for _, loanMap := range loansData {
				if isReturn, ok := loanMap["is_return"].(bool); ok && !isReturn {
					adaTanggungan = true
					break
				}
			}

			if adaTanggungan {
				response["status"] = "Tanggungan"
				response["status_pinjaman"] = "Belum"

				// Ambil status approval
				statusURL := fmt.Sprintf("http://localhost:8000/api/status_approval?kode_user=%s", kodeUser)
				statusResp, err := client.Get(statusURL)
				if err == nil {
					defer statusResp.Body.Close()
					var statApproval []map[string]interface{}
					if err := json.NewDecoder(statusResp.Body).Decode(&statApproval); err == nil {
						if len(statApproval) > 0 {
							if keterangan, ok := statApproval[0]["keterangan"].(string); ok {
								response["keterangan"] = keterangan
							}
						}
					}
				}
			}

			mu.Lock()
			hasil = append(hasil, response)
			mu.Unlock()
		}(m)
	}

	wg.Wait()
	return hasil, nil
}
