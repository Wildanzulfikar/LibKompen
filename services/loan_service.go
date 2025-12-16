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

	// 🔥 SINGLEFLIGHT: hanya 1 fetch jalan
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

	// 🔥 ambil cache mahasiswa SEKALI
	mahasiswaMap := getMahasiswaCache()

	layout := "2006-01-02"
	formatted := make([]map[string]interface{}, 0)

	for _, item := range loansData {
		loan := item.(map[string]interface{})

		memberID := strings.TrimSpace(fmt.Sprintf("%v", loan["member_id"]))
		mahasiswa := mahasiswaMap[memberID]

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

func FetchLoanDetail(memberID string) (map[string]interface{}, error) {
	if memberID == "" {
		return nil, fmt.Errorf("missing member_id parameter")
	}

	// Ambil semua loan milik member_id
	opacLoanURL := fmt.Sprintf("http://localhost:8080/loan?member_id=%s", url.QueryEscape(memberID))
	respLoan, err := http.Get(opacLoanURL)
	if err != nil {
		return nil, fmt.Errorf("Gagal ambil data loan dari OPAC")
	}
	defer respLoan.Body.Close()

	loanBytes, err := io.ReadAll(respLoan.Body)
	if err != nil {
		return nil, fmt.Errorf("Gagal membaca response body loan")
	}

	var loanResult map[string]interface{}
	if err := json.Unmarshal(loanBytes, &loanResult); err != nil {
		return nil, fmt.Errorf("Gagal decode data loan: %v", err)
	}

	loansData, ok := loanResult["data"].([]interface{})
	if !ok {
		return nil, fmt.Errorf("Data loan tidak ditemukan atau format salah")
	}
	if len(loansData) == 0 {
		return nil, fmt.Errorf("Tidak ada data loan untuk member_id ini")
	}

	// Data mahasiswa (ambil dari loan pertama saja, asumsikan sama)
	var mahasiswaData map[string]interface{}
	sikompenURL := fmt.Sprintf("http://localhost:8000/api/mahasiswa?nim=%s", memberID)
	respMhs, err := http.Get(sikompenURL)
	if err == nil {
		defer respMhs.Body.Close()
		mhsBytes, _ := io.ReadAll(respMhs.Body)
		var mhsArr []map[string]interface{}
		if err := json.Unmarshal(mhsBytes, &mhsArr); err == nil && len(mhsArr) > 0 {
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
		}
	}

	// Filter mahasiswaData ada (aktif)
	if mahasiswaData == nil {
		return map[string]interface{}{
			"message": "Data mahasiswa tidak ditemukan atau tidak aktif",
		}, nil
	}

	layout := "2006-01-02"
	peminjamanArr := make([]map[string]interface{}, 0)
	bukuArr := make([]map[string]interface{}, 0)
	for _, item := range loansData {
		loan := item.(map[string]interface{})
		// Ambil is_return dari loan
		var isReturn interface{}
		if v, ok := loan["is_return"]; ok {
			isReturn = v
		}
		// Ambil detail loan 
		var detailLoan map[string]interface{}
		if loanID, ok := loan["loan_id"]; ok {
			detailLoanURL := fmt.Sprintf("http://localhost:8080/loan/%v", loanID)
			respDetail, err := http.Get(detailLoanURL)
			if err == nil {
				defer respDetail.Body.Close()
				detailBytes, _ := io.ReadAll(respDetail.Body)
				_ = json.Unmarshal(detailBytes, &detailLoan)
				if isReturn == nil {
					if v, ok := detailLoan["is_return"]; ok {
						isReturn = v
					}
				}
			}
		}
		status := "Belum"
		switch v := isReturn.(type) {
		case bool:
			if v {
				status = "Lunas"
			}
		case float64:
			if v == 1 {
				status = "Lunas"
			}
		case int:
			if v == 1 {
				status = "Lunas"
			}
		case int64:
			if v == 1 {
				status = "Lunas"
			}
		case string:
			s := strings.TrimSpace(strings.ToLower(v))
			if s == "1" || s == "true" {
				status = "Lunas"
			}
		default:
			if fmt.Sprintf("%v", v) == "1" {
				status = "Lunas"
			}
		}
		keterlambatan := "-"
		dueDate, _ := loan["due_date"].(string)
		returnDate, _ := loan["return_date"].(string)
		if dueDate != "" {
			tDue, errDue := time.Parse(layout, dueDate)
			var tEnd time.Time
			var errEnd error
			if returnDate == "" || returnDate == "null" {
				tEnd = time.Now()
			} else {
				tEnd, errEnd = time.Parse(layout, returnDate)
				if errEnd != nil {
					tEnd = time.Now()
				}
			}
			if errDue == nil {
				daysLate := int(tEnd.Sub(tDue).Hours() / 24)
				if daysLate > 0 {
					keterlambatan = fmt.Sprintf("%d Hari", daysLate)
				} else {
					keterlambatan = "Tepat Waktu"
				}
			}
		}
		peminjamanArr = append(peminjamanArr, map[string]interface{}{
			"loan_id":       loan["loan_id"],
			"peminjaman":    loan["loan_date"],
			"tenggat_waktu": loan["due_date"],
			"pengembalian":  loan["return_date"],
			"status":        status,
			"keterlambatan": keterlambatan,
		})

		// Ambil detail buku dengan fetch detail loan per loan_id
		if loanID, ok := loan["loan_id"]; ok {
			detailLoanURL := fmt.Sprintf("http://localhost:8080/loan/%v", loanID)
			respDetail, err := http.Get(detailLoanURL)
			if err == nil {
				defer respDetail.Body.Close()
				detailBytes, _ := io.ReadAll(respDetail.Body)
				var detailLoan map[string]interface{}
				if err := json.Unmarshal(detailBytes, &detailLoan); err == nil {
					itemCode, _ := detailLoan["item_code"].(string)
					if itemCode != "" {
						opacItemURL := fmt.Sprintf("http://localhost:8080/item/%s", itemCode)
						respItem, err := http.Get(opacItemURL)
						if err == nil {
							defer respItem.Body.Close()
							itemBytes, _ := io.ReadAll(respItem.Body)
							var itemData map[string]interface{}
							if err := json.Unmarshal(itemBytes, &itemData); err == nil {
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
										var biblioData map[string]interface{}
										if err := json.Unmarshal(biblioBytes, &biblioData); err == nil {
											buku := map[string]interface{}{
												"loan_id":      loan["loan_id"],
												"title":        biblioData["title"],
												"edition":      biblioData["edition"],
												"isbn_issn":    biblioData["isbn_issn"],
												"publish_year": biblioData["publish_year"],
												"collation":    biblioData["collation"],
												"call_number":  biblioData["call_number"],
											}
											if buku["title"] != nil && buku["title"] != "" {
												bukuArr = append(bukuArr, buku)
											}
										}
									}
								}
							}
						}
					}
				}
			}
		}
	}

	detail := make(map[string]interface{})
	detail["mahasiswa"] = map[string]interface{}{
		"nama":     mahasiswaData["nama_user"],
		"nim":      mahasiswaData["kode_user"],
		"prodi":    mahasiswaData["prodi"],
		"kelas":    mahasiswaData["kelas"],
		"semester": mahasiswaData["semester"],
	}
	detail["peminjaman"] = peminjamanArr
	detail["buku"] = bukuArr
	return detail, nil
}
