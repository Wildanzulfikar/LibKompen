package services

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type MahasiswaBebasPustakaResponse struct {
	IDMahasiswa      interface{} `json:"id_mahasiswa"`
	Nim              interface{} `json:"nim"`
	Nama             interface{} `json:"nama"`
	Prodi            interface{} `json:"prodi"`
	Kelas            interface{} `json:"kelas"`
	Semester         interface{} `json:"semester"`
	Status           interface{} `json:"status"`
	StatusPinjaman   interface{} `json:"status_pinjaman"`
	Keterangan       interface{} `json:"keterangan"`
}

func GetMahasiswaBebasPustakaService(memberID string) ([]MahasiswaBebasPustakaResponse, error) {
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

	var hasil []MahasiswaBebasPustakaResponse
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

		response := MahasiswaBebasPustakaResponse{
			IDMahasiswa:    idMahasiswa,
			Nim:            kodeUser,
			Nama:           namaUser,
			Prodi:          prodi,
			Kelas:          kelas,
			Semester:       semester,
			Status:         "Bebas Pustaka",
			StatusPinjaman: "Lunas",
			Keterangan:     "-",
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
			response.Status = "Tanggungan"
			response.StatusPinjaman = "Belum"

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
							response.Keterangan = keterangan
						}
					}
				}
			}
		}

		hasil = append(hasil, response)
	}
	return hasil, nil
}
