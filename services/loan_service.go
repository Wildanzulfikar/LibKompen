package services

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

func GetAllLoanFormatted() ([]map[string]interface{}, error) {
	resp, err := http.Get("http://localhost:8080/loan")
	if err != nil {
		return nil, fmt.Errorf("Gagal ambil data loan dari OPAC")
	}
	defer resp.Body.Close()

	bodyBytes, _ := io.ReadAll(resp.Body)
	var result map[string]interface{}
	if err := json.Unmarshal(bodyBytes, &result); err != nil {
		return nil, fmt.Errorf("Gagal decode data loan")
	}

	loansData, ok := result["data"].([]interface{})
	if !ok {
		return nil, fmt.Errorf("Data loan tidak valid")
	}

	formatted := make([]map[string]interface{}, 0)
	for _, loanItem := range loansData {
		loanMap, ok := loanItem.(map[string]interface{})
		if !ok {
			continue
		}
		memberID, _ := loanMap["member_id"].(string)
		var mahasiswa map[string]interface{}
		if memberID != "" {
			sikompenURL := "http://localhost:8000/api/mahasiswa?nim=" + memberID
			respMhs, err := http.Get(sikompenURL)
			if err == nil {
				defer respMhs.Body.Close()
				mhsBytes, _ := io.ReadAll(respMhs.Body)
				var mhsArr []map[string]interface{}
				if err := json.Unmarshal(mhsBytes, &mhsArr); err == nil && len(mhsArr) > 0 {
					for _, m := range mhsArr {
						if m["kode_user"] == memberID {
							mahasiswa = m
							break
						}
					}
				}
			}
		}

		status := "Belum"
		if val, ok := loanMap["is_return"].(bool); ok && val {
			status = "Lunas"
		} else if val, ok := loanMap["is_return"].(float64); ok && val == 1 {
			status = "Lunas"
		}

		keterlambatan := "-"
		dueDate, _ := loanMap["due_date"].(string)
		returnDate, _ := loanMap["return_date"].(string)
		layout := "2006-01-02"
		var daysLate int
		if dueDate != "" {
			if returnDate == "" {
				now := fmt.Sprintf("%04d-%02d-%02d", time.Now().Year(), time.Now().Month(), time.Now().Day())
				tDue, err1 := time.Parse(layout, dueDate)
				tNow, err2 := time.Parse(layout, now)
				if err1 == nil && err2 == nil {
					daysLate = int(tNow.Sub(tDue).Hours() / 24)
				}
			} else {
				tDue, err1 := time.Parse(layout, dueDate)
				tReturn, err2 := time.Parse(layout, returnDate)
				if err1 == nil && err2 == nil {
					daysLate = int(tReturn.Sub(tDue).Hours() / 24)
				}
			}
			if daysLate > 0 {
				keterlambatan = fmt.Sprintf("%d Hari", daysLate)
			} else {
				keterlambatan = "Tepat Waktu"
			}
		}

		formatted = append(formatted, map[string]interface{}{
			"id_mahasiswa":  mahasiswa["id_mahasiswa"],
			"nim":           mahasiswa["kode_user"],
			"nama":          mahasiswa["nama_user"],
			"prodi":         mahasiswa["prodi"],
			"kelas":         mahasiswa["kelas"],
			"semester":      mahasiswa["semester"],
			"peminjaman":    loanMap["loan_date"],
			"tenggat_waktu": loanMap["due_date"],
			"pengembalian":  loanMap["return_date"],
			"keterlambatan": keterlambatan,
			"status":        status,
		})
	}

	return formatted, nil
}

func FetchLoanDetail(loanID string) (map[string]interface{}, error) {
	if loanID == "" {
		return nil, fmt.Errorf("missing loan_id parameter")
	}

	// Data loan
	opacLoanURL := fmt.Sprintf("http://localhost:8080/loan/%s", loanID)
	respLoan, err := http.Get(opacLoanURL)
	if err != nil {
		return nil, fmt.Errorf("Gagal ambil detail loan dari OPAC")
	}
	defer respLoan.Body.Close()

	loanBytes, err := io.ReadAll(respLoan.Body)
	if err != nil {
		return nil, fmt.Errorf("Gagal membaca response body loan")
	}

	var loanResult map[string]interface{}
	if err := json.Unmarshal(loanBytes, &loanResult); err != nil {
		return nil, fmt.Errorf("Gagal decode detail loan: %v", err)
	}

	// Data mahasiswa
	var mahasiswaData map[string]interface{}
	memberID, _ := loanResult["member_id"].(string)
	if memberID != "" {
		sikompenURL := fmt.Sprintf("http://localhost:8000/api/mahasiswa?nim=%s", memberID)
		respMhs, err := http.Get(sikompenURL)
		if err == nil {
			defer respMhs.Body.Close()
			mhsBytes, _ := io.ReadAll(respMhs.Body)
			var mhsArr []map[string]interface{}
			if err := json.Unmarshal(mhsBytes, &mhsArr); err == nil && len(mhsArr) > 0 {
				mahasiswaData = mhsArr[0]
			}
		}
	}

	// Data item
	var itemData map[string]interface{}
	itemCode, _ := loanResult["item_code"].(string)
	if itemCode != "" {
		opacItemURL := fmt.Sprintf("http://localhost:8080/item/%s", itemCode)
		respItem, err := http.Get(opacItemURL)
		if err == nil {
			defer respItem.Body.Close()
			itemBytes, _ := io.ReadAll(respItem.Body)
			_ = json.Unmarshal(itemBytes, &itemData)
			if data, ok := itemData["data"].(map[string]interface{}); ok {
				itemData = data
			}
		}
	}

	var biblioData map[string]interface{}
	if itemData != nil {
		var biblioID string
		switch v := itemData["biblio_id"].(type) {
		case string:
			biblioID = v
		case float64:
			biblioID = fmt.Sprintf("%.0f", v)
		}
		if biblioID != "" {
			opacBiblioURL := fmt.Sprintf("http://localhost:8080/biblio/%s", biblioID)
			respBiblio, err := http.Get(opacBiblioURL)
			if err == nil {
				defer respBiblio.Body.Close()
				biblioBytes, _ := io.ReadAll(respBiblio.Body)
				_ = json.Unmarshal(biblioBytes, &biblioData)
				if data, ok := biblioData["data"].(map[string]interface{}); ok {
					biblioData = data
				}
			}
		}
	}

	detail := make(map[string]interface{})
	if mahasiswaData != nil {
		detail["mahasiswa"] = map[string]interface{}{
			"nama":     mahasiswaData["nama_user"],
			"nim":      mahasiswaData["kode_user"],
			"prodi":    mahasiswaData["prodi"],
			"kelas":    mahasiswaData["kelas"],
			"semester": mahasiswaData["semester"],
		}
	}
	if loanResult != nil {
		status := "Belum"
		if val, ok := loanResult["is_return"].(bool); ok && val {
			status = "Lunas"
		} else if val, ok := loanResult["is_return"].(float64); ok && val == 1 {
			status = "Lunas"
		}
		keterlambatan := "-"
		dueDate, _ := loanResult["due_date"].(string)
		returnDate, _ := loanResult["return_date"].(string)
		layout := "2006-01-02"
		var daysLate int
		if dueDate != "" {
			if returnDate == "" {
				now := fmt.Sprintf("%04d-%02d-%02d", time.Now().Year(), time.Now().Month(), time.Now().Day())
				tDue, err1 := time.Parse(layout, dueDate)
				tNow, err2 := time.Parse(layout, now)
				if err1 == nil && err2 == nil {
					daysLate = int(tNow.Sub(tDue).Hours() / 24)
				}
			} else {
				tDue, err1 := time.Parse(layout, dueDate)
				tReturn, err2 := time.Parse(layout, returnDate)
				if err1 == nil && err2 == nil {
					daysLate = int(tReturn.Sub(tDue).Hours() / 24)
				}
			}
			if daysLate > 0 {
				keterlambatan = fmt.Sprintf("%d Hari", daysLate)
			} else {
				keterlambatan = "Tepat Waktu"
			}
		}
		detail["peminjaman"] = map[string]interface{}{
			"peminjaman":    loanResult["loan_date"],
			"tenggat_waktu": loanResult["due_date"],
			"pengembalian":  loanResult["return_date"],
			"status":        status,
			"keterlambatan": keterlambatan,
		}
	}
	if biblioData != nil {
		detail["buku"] = map[string]interface{}{
			"title":        biblioData["title"],
			"edition":      biblioData["edition"],
			"isbn_issn":    biblioData["isbn_issn"],
			"publish_year": biblioData["publish_year"],
			"collation":    biblioData["collation"],
			"call_number":  biblioData["call_number"],
		}
	}
	return detail, nil
}
