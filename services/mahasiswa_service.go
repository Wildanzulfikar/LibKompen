package services

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)



func GetMahasiswaBebasPustakaService(memberID string, jurusan string, statusPustaka string, statusPinjaman string, search string, tahun string) ([]map[string]interface{}, error) {
	sikompenURL := "http://localhost:8000/api/mahasiswa?limit=0"
	respMahasiswa, err := http.Get(sikompenURL)
	if err != nil {
		return nil, fmt.Errorf("Gagal ambil data mahasiswa Sikompen: %w", err)
	}
	defer respMahasiswa.Body.Close()

	var mahasiswaList []map[string]interface{}
	bodyBytes, _ := io.ReadAll(respMahasiswa.Body)
	if err := json.Unmarshal(bodyBytes, &mahasiswaList); err != nil {
		return nil, fmt.Errorf("Gagal decode mahasiswa: %w", err)
	}

	var hasil []map[string]interface{}
	   for _, m := range mahasiswaList {
		   kodeUser, ok := m["kode_user"].(string)
		   if !ok {
			   continue
		   }
		   if memberID != "" && kodeUser != memberID {
			   continue
		   }
		   if jurusan != "" && len(kodeUser) >= 4 && kodeUser[2:4] != jurusan {
			   continue
		   }

		   namaUser, _ := m["nama_user"].(string)
		   prodi, _ := m["prodi"].(string)
		   kelas, _ := m["kelas"].(string)
		   semester, _ := m["semester"].(string)
		   idMahasiswa := m["id_mahasiswa"]

		   if search != "" {
			   if !(containsIgnoreCase(namaUser, search) || containsIgnoreCase(kodeUser, search)) {
				   continue
			   }
		   }

func containsIgnoreCase(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || (len(substr) > 0 && (containsFold(s, substr))))
}

func containsFold(s, substr string) bool {
	return len(substr) > 0 && (len(s) >= len(substr)) && (indexFold(s, substr) >= 0)
}

func indexFold(s, substr string) int {
	for i := 0; i+len(substr) <= len(s); i++ {
		if equalFold(s[i:i+len(substr)], substr) {
			return i
		}
	}
	return -1
}

func equalFold(s, t string) bool {
	if len(s) != len(t) {
		return false
	}
	for i := 0; i < len(s); i++ {
		a, b := s[i], t[i]
		if a == b {
			continue
		}
		if 'A' <= a && a <= 'Z' {
			a += 'a' - 'A'
		}
		if 'A' <= b && b <= 'Z' {
			b += 'a' - 'A'
		}
		if a != b {
			return false
		}
	}
	return true
}

		   response := map[string]interface{}{
			   "id_mahasiswa":    idMahasiswa,
			   "nim":            kodeUser,
			   "nama":           namaUser,
			   "prodi":          prodi,
			   "kelas":          kelas,
			   "semester":       semester,
			   "status":         "Bebas Pustaka",
			   "status_pinjaman": "Lunas",
			   "keterangan":     "-",
		   }

		   loanURL := fmt.Sprintf("http://localhost:8080/loan?member_id=%s", kodeUser)
		   loanResp, err := http.Get(loanURL)
		   if err != nil {
			   if tahun != "" {
				   foundYear := false
				   for _, loanItem := range loansData {
					   loanMap, ok := loanItem.(map[string]interface{})
					   if !ok {
						   continue
					   }
					   returnDate, _ := loanMap["return_date"].(string)
					   if len(returnDate) >= 4 && returnDate[:4] == tahun {
						   foundYear = true
						   break
					   }
				   }
				   if !foundYear {
					   continue
				   }
			   }
			   hasil = append(hasil, response)
			   continue
		   }
		   defer loanResp.Body.Close()

		   bodyBytes, _ := io.ReadAll(loanResp.Body)
		   var loanPayload map[string]interface{}
		   if err := json.Unmarshal(bodyBytes, &loanPayload); err != nil {
			   hasil = append(hasil, response)
			   continue
		   }

		   loansData, ok := loanPayload["data"].([]interface{})
		   if !ok {
			   hasil = append(hasil, response)
			   continue
		   }

		   adaTanggungan := false
		   for _, loanItem := range loansData {
			   loanMap, ok := loanItem.(map[string]interface{})
			   if !ok {
				   continue
			   }
			   isReturn, ok := loanMap["is_return"].(bool)
			   if ok && !isReturn {
				   adaTanggungan = true
				   break
			   }
		   }

		   if adaTanggungan {
			   response["status"] = "Tanggungan"
			   response["status_pinjaman"] = "Belum"

			   statusURL := fmt.Sprintf("http://localhost:8000/api/status_approval?kode_user=%s", kodeUser)
			   statusResp, err := http.Get(statusURL)
			   if err == nil {
				   defer statusResp.Body.Close()
				   bodyBytes, _ := io.ReadAll(statusResp.Body)
				   var statApproval []map[string]interface{}
				   if err := json.Unmarshal(bodyBytes, &statApproval); err == nil {
					   if len(statApproval) > 0 {
						   keterangan, ok := statApproval[0]["keterangan"].(string)
						   if ok {
							   response["keterangan"] = keterangan
						   }
					   }
				   }
			   }
		   }


		   if statusPustaka != "" {
			   statusVal, _ := response["status"].(string)
			   if statusVal != statusPustaka {
				   continue
			   }
		   }
		   if statusPinjaman != "" {
			   pinjamanVal, _ := response["status_pinjaman"].(string)
			   if pinjamanVal != statusPinjaman {
				   continue
			   }
		   }

		   hasil = append(hasil, response)
	   }
	   return hasil, nil
}
