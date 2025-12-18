package services

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
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

	// SINGLEFLIGHT
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
					if _, exists := cache[k]; !exists {
						cache[k] = m
					}
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

var allLoansCache = struct {
	sync.RWMutex
	data      []map[string]interface{}
	lastFetch time.Time
}{}

func fetchAllLoans() ([]map[string]interface{}, error) {
	allLoansCache.RLock()
	if time.Since(allLoansCache.lastFetch) < 10*time.Minute && len(allLoansCache.data) > 0 {
		defer allLoansCache.RUnlock()
		return allLoansCache.data, nil
	}
	allLoansCache.RUnlock()

	client := &http.Client{Timeout: 10 * time.Second}
	allLoans := make([]map[string]interface{}, 0)
	currentPage := 1

	for {
		urlStr := fmt.Sprintf("http://localhost:8080/loan?page=%d&per_page=1000", currentPage)
		resp, err := client.Get(urlStr)
		if err != nil {
			return nil, err
		}
		body, _ := io.ReadAll(resp.Body)
		resp.Body.Close()

		var result map[string]interface{}
		if err := json.Unmarshal(body, &result); err != nil {
			return nil, err
		}

		loansData, ok := result["data"].([]interface{})
		if !ok || len(loansData) == 0 {
			break
		}

		for _, item := range loansData {
			allLoans = append(allLoans, item.(map[string]interface{}))
		}

		meta := result["meta"].(map[string]interface{})
		total := int(meta["total"].(float64))
		if currentPage*1000 >= total {
			break
		}
		currentPage++
	}

	allLoansCache.Lock()
	allLoansCache.data = allLoans
	allLoansCache.lastFetch = time.Now()
	allLoansCache.Unlock()

	return allLoans, nil
}

func GetAllLoanFormatted(page, perPage int, search string) ([]map[string]interface{}, int, error) {
	// Ambil semua loan dari cache / fetch baru
	allLoans, err := fetchAllLoans()
	if err != nil {
		return nil, 0, err
	}

	// Ambil master mahasiswa
	mahasiswaMap := getMahasiswaCache()
	layout := "2006-01-02"

	// Filter loan sesuai mahasiswa
	filtered := make([]map[string]interface{}, 0)
	for _, loan := range allLoans {
		memberID := strings.TrimSpace(fmt.Sprintf("%v", loan["member_id"]))
		mahasiswa := mahasiswaMap[memberID]
		if mahasiswa == nil {
			continue
		}

		// hitung keterlambatan
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

		filtered = append(filtered, map[string]interface{}{
			"id_mahasiswa":  mahasiswa["id_mahasiswa"],
			"nim":           mahasiswa["kode_user"],
			"nama":          mahasiswa["nama_user"],
			"prodi":         mahasiswa["prodi"],
			"kelas":         mahasiswa["kelas"],
			"semester":      mahasiswa["semester"],
			"peminjaman":    loan["loan_date"],
			"tenggat_waktu": loan["due_date"],
			"pengembalian":  loan["return_date"],
			"keterlambatan": keterlambatan,
			"status":        loan["status"],
			"is_return":     loan["is_return"],
		})
	}

	// Pagination dari filtered
	total := len(filtered)
	start := (page - 1) * perPage
	end := start + perPage
	if start > total {
		start = total
	}
	if end > total {
		end = total
	}
	paginated := filtered[start:end]

	return paginated, total, nil
}

func FetchLoanDetail(loanID string) (map[string]interface{}, error) {
	if loanID == "" {
		return nil, fmt.Errorf("missing loan_id parameter")
	}

	// =========================
	// FETCH LOAN DETAIL
	// =========================
	loanURL := fmt.Sprintf("http://localhost:8080/loan/%s", loanID)
	respLoan, err := http.Get(loanURL)
	if err != nil {
		return nil, fmt.Errorf("gagal ambil detail loan")
	}
	defer respLoan.Body.Close()

	loanBytes, err := io.ReadAll(respLoan.Body)
	if err != nil {
		return nil, fmt.Errorf("gagal membaca response loan")
	}

	var loanResult map[string]interface{}
	if err := json.Unmarshal(loanBytes, &loanResult); err != nil {
		return nil, fmt.Errorf("gagal decode loan: %v", err)
	}

	// =========================
	// NORMALIZE MEMBER ID
	// =========================
	var memberID string
	switch v := loanResult["member_id"].(type) {
	case string:
		memberID = v
	case float64:
		memberID = fmt.Sprintf("%.0f", v)
	case int:
		memberID = fmt.Sprintf("%d", v)
	}

	// =========================
	// FETCH MAHASISWA
	// =========================
	var mahasiswaData map[string]interface{}
	if memberID != "" {
		mhsURL := fmt.Sprintf("http://localhost:8000/api/mahasiswa?nim=%s", memberID)
		respMhs, err := http.Get(mhsURL)
		if err == nil {
			defer respMhs.Body.Close()
			mhsBytes, _ := io.ReadAll(respMhs.Body)

			var mhsArr []map[string]interface{}
			if json.Unmarshal(mhsBytes, &mhsArr) == nil {
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
				if mahasiswaData == nil && len(mhsArr) > 0 {
					mahasiswaData = mhsArr[0]
				}
			}
		}
	}

	// =========================
	// FETCH ITEM
	// =========================
	var itemData map[string]interface{}
	if itemCode, ok := loanResult["item_code"].(string); ok && itemCode != "" {
		itemURL := fmt.Sprintf("http://localhost:8080/item/%s", itemCode)
		respItem, err := http.Get(itemURL)
		if err == nil {
			defer respItem.Body.Close()
			itemBytes, _ := io.ReadAll(respItem.Body)
			_ = json.Unmarshal(itemBytes, &itemData)
			if d, ok := itemData["data"].(map[string]interface{}); ok {
				itemData = d
			}
		}
	}

	// =========================
	// FETCH BIBLIO
	// =========================
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
			biblioURL := fmt.Sprintf("http://localhost:8080/biblio/%s", biblioID)
			respBiblio, err := http.Get(biblioURL)
			if err == nil {
				defer respBiblio.Body.Close()
				biblioBytes, _ := io.ReadAll(respBiblio.Body)
				_ = json.Unmarshal(biblioBytes, &biblioData)
				if d, ok := biblioData["data"].(map[string]interface{}); ok {
					biblioData = d
				}
			}
		}
	}

	// =========================
	// HITUNG STATUS & TELAT
	// =========================
	status := "Belum"
	switch v := loanResult["is_return"].(type) {
	case bool:
		if v {
			status = "Lunas"
		}
	case float64:
		if v == 1 {
			status = "Lunas"
		}
	}

	keterlambatan := "-"
	layout := "2006-01-02"

	dueDate, _ := loanResult["due_date"].(string)
	returnDate, _ := loanResult["return_date"].(string)

	if dueDate != "" {
		tDue, err1 := time.Parse(layout, dueDate)
		var tEnd time.Time
		var err2 error

		if returnDate == "" || returnDate == "null" {
			tEnd = time.Now()
		} else {
			tEnd, err2 = time.Parse(layout, returnDate)
			if err2 != nil {
				tEnd = time.Now()
			}
		}

		if err1 == nil {
			daysLate := int(tEnd.Sub(tDue).Hours() / 24)
			if daysLate > 0 {
				keterlambatan = fmt.Sprintf("%d Hari", daysLate)
			} else {
				keterlambatan = "Tepat Waktu"
			}
		}
	}

	// =========================
	// FINAL RESPONSE
	// =========================
	result := map[string]interface{}{}

	if mahasiswaData != nil {
		result["mahasiswa"] = map[string]interface{}{
			"nama":     mahasiswaData["nama_user"],
			"nim":      mahasiswaData["kode_user"],
			"prodi":    mahasiswaData["prodi"],
			"kelas":    mahasiswaData["kelas"],
			"semester": mahasiswaData["semester"],
		}
	}

	result["peminjaman"] = map[string]interface{}{
		"loan_id":       loanResult["loan_id"],
		"peminjaman":    loanResult["loan_date"],
		"tenggat_waktu": loanResult["due_date"],
		"pengembalian":  loanResult["return_date"],
		"status":        status,
		"keterlambatan": keterlambatan,
	}

	if biblioData != nil {
		result["buku"] = map[string]interface{}{
			"title":        biblioData["title"],
			"edition":      biblioData["edition"],
			"isbn_issn":    biblioData["isbn_issn"],
			"publish_year": biblioData["publish_year"],
			"collation":    biblioData["collation"],
			"call_number":  biblioData["call_number"],
		}
	}

	return result, nil
}
