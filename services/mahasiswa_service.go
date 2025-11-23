package services

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)



func GetMahasiswaBebasPustakaService(memberID string) ([]map[string]interface{}, error) {
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
		   namaUser, _ := m["nama_user"].(string)
		   prodi, _ := m["prodi"].(string)
		   kelas, _ := m["kelas"].(string)
		   semester, _ := m["semester"].(string)
		   idMahasiswa := m["id_mahasiswa"]

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

		   hasil = append(hasil, response)
	   }
	   return hasil, nil
}
