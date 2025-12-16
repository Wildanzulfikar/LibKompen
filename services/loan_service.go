package services

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"
	"golang.org/x/sync/singleflight"
)

var mhsFetchGroup singleflight.Group

var mahasiswaCache = struct {
	sync.RWMutex
	data      map[string]map[string]interface{}
	lastFetch time.Time
}{
	data: make(map[string]map[string]interface{}),
}

func getMahasiswaCache() map[string]map[string]interface{} {
	mahasiswaCache.RLock()
	if time.Since(mahasiswaCache.lastFetch) < 10*time.Minute && len(mahasiswaCache.data) > 0 {
		defer mahasiswaCache.RUnlock()
		return mahasiswaCache.data
	}
	mahasiswaCache.RUnlock()

	_, _, _ = mhsFetchGroup.Do("fetch-mahasiswa", func() (interface{}, error) {

		client := &http.Client{Timeout: 5 * time.Second}
		resp, err := client.Get("http://localhost:8000/api/mahasiswa")
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()

		body, _ := io.ReadAll(resp.Body)
		var arr []map[string]interface{}
		if err := json.Unmarshal(body, &arr); err != nil {
			return nil, err
		}

		cache := make(map[string]map[string]interface{})
		for _, m := range arr {
			if v, ok := m["kode_user"]; ok {
				k := strings.TrimSpace(fmt.Sprintf("%v", v))
				if k != "" {
					cache[k] = m
				}
			}
			if v, ok := m["id_mahasiswa"]; ok {
				k := strings.TrimSpace(fmt.Sprintf("%v", v))
				if k != "" {
					cache[k] = m
				}
			}
		}

		mahasiswaCache.Lock()
		mahasiswaCache.data = cache
		mahasiswaCache.lastFetch = time.Now()
		mahasiswaCache.Unlock()

		return nil, nil
	})

	mahasiswaCache.RLock()
	defer mahasiswaCache.RUnlock()
	return mahasiswaCache.data
}

func GetAllLoanFormatted(page, perPage int, search string) ([]map[string]interface{}, int, error) {
	url := fmt.Sprintf("http://localhost:8080/loan?page=%d&per_page=%d&search=%s", page, perPage, url.QueryEscape(search))

	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Get(url)
	if err != nil {
		return nil, 0, fmt.Errorf("gagal ambil data loan")
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, 0, fmt.Errorf("decode gagal")
	}

	loansData := result["data"].([]interface{})
	meta := result["meta"].(map[string]interface{})
	total := int(meta["total"].(float64))

	mahasiswaMap := getMahasiswaCache()

	layout := "2006-01-02"
	formatted := make([]map[string]interface{}, 0)

	for _, item := range loansData {
		loan := item.(map[string]interface{})

		memberID := strings.TrimSpace(fmt.Sprintf("%v", loan["member_id"]))
		mahasiswa := mahasiswaMap[memberID]

		if mahasiswa == nil {
			fmt.Println("TIDAK MATCH:", memberID)
		} else {
			fmt.Println("MATCH:", memberID, mahasiswa["nama_user"])
		}

		// keterlambatan
		keterlambatan := "-"
		dueDate, _ := loan["due_date"].(string)
		returnDate, _ := loan["return_date"].(string)

		if dueDate != "" {
			tDue, _ := time.Parse(layout, dueDate)
			tReturn := time.Now()
			if returnDate != "" {
				tReturn, _ = time.Parse(layout, returnDate)
			}
			days := int(tReturn.Sub(tDue).Hours() / 24)
			if days > 0 {
				keterlambatan = fmt.Sprintf("%d Hari", days)
			} else {
				keterlambatan = "Tepat Waktu"
			}
		}

		var idMhs, nim, nama, prodi, kelas, semester interface{}
		if mahasiswa != nil {
			idMhs = mahasiswa["id_mahasiswa"]
			nim = mahasiswa["kode_user"]
			nama = mahasiswa["nama_user"]
			prodi = mahasiswa["prodi"]
			kelas = mahasiswa["kelas"]
			semester = mahasiswa["semester"]
		}

		formatted = append(formatted, map[string]interface{}{
			"id_mahasiswa":  idMhs,
			"nim":           nim,
			"nama":          nama,
			"prodi":         prodi,
			"kelas":         kelas,
			"semester":      semester,
			"peminjaman":    loan["loan_date"],
			"tenggat_waktu": loan["due_date"],
			"pengembalian":  loan["return_date"],
			"keterlambatan": keterlambatan,
			"status":        loan["status"],
		})
	}

	return formatted, total, nil
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
	var memberID string
	switch v := loanResult["member_id"].(type) {
	case string:
		memberID = v
	case float64:
		memberID = fmt.Sprintf("%.0f", v)
	case int:
		memberID = fmt.Sprintf("%d", v)
	default:
		memberID = ""
	}
	if memberID != "" {
		sikompenURL := fmt.Sprintf("http://localhost:8000/api/mahasiswa?nim=%s", memberID)
		respMhs, err := http.Get(sikompenURL)
		if err == nil {
			defer respMhs.Body.Close()
			mhsBytes, _ := io.ReadAll(respMhs.Body)
			var mhsArr []map[string]interface{}
			if err := json.Unmarshal(mhsBytes, &mhsArr); err == nil && len(mhsArr) > 0 {
				// Try to find exact kode_user match, similar to above
				for _, m := range mhsArr {
					var kodeUser string
					switch kv := m["kode_user"].(type) {
					case string:
						kodeUser = kv
					case float64:
						kodeUser = fmt.Sprintf("%.0f", kv)
					}
					if kodeUser == memberID {
						mahasiswaData = m
						break
					}
				}
				if mahasiswaData == nil {
					mahasiswaData = mhsArr[0]
				}
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
