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

var loanStatusCache = struct {
	sync.RWMutex
	data      map[string]int
	lastFetch map[string]time.Time
}{
	data:      make(map[string]int),
	lastFetch: make(map[string]time.Time),
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

func getActiveLoanCount(memberID string) int {
	loanStatusCache.RLock()
	if t, ok := loanStatusCache.lastFetch[memberID]; ok {
		if time.Since(t) < 10*time.Minute {
			defer loanStatusCache.RUnlock()
			return loanStatusCache.data[memberID]
		}
	}
	loanStatusCache.RUnlock()

	url := fmt.Sprintf("http://localhost:8080/api/loan/%s", memberID)
	resp, err := http.Get(url)
	if err != nil {
		return 0
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		return 0
	}

	aktif := 0
	if arr, ok := result["pinjaman_aktif"].([]interface{}); ok {
		aktif = len(arr)
	}

	loanStatusCache.Lock()
	loanStatusCache.data[memberID] = aktif
	loanStatusCache.lastFetch[memberID] = time.Now()
	loanStatusCache.Unlock()

	return aktif
}

func GetAllLoanFormatted(page, perPage int, search string) ([]map[string]interface{}, int, error) {
	allLoans, err := fetchAllLoans()
	if err != nil {
		return nil, 0, err
	}

	mahasiswaMap := getMahasiswaCache()

	type Agg struct {
		Mahasiswa map[string]interface{}
		Total     int
	}

	group := make(map[string]*Agg)

	for _, loan := range allLoans {
		memberID := strings.TrimSpace(fmt.Sprintf("%v", loan["member_id"]))
		if memberID == "" {
			continue
		}

		mhs := mahasiswaMap[memberID]
		if mhs == nil {
			continue
		}

		if _, ok := group[memberID]; !ok {
			group[memberID] = &Agg{
				Mahasiswa: mhs,
				Total:     0,
			}
		}

		group[memberID].Total++
	}

	result := make([]map[string]interface{}, 0)

	for nim, g := range group {
		aktif := getActiveLoanCount(nim)

		status := "Lunas"
		if aktif > 0 {
			status = "Belum Lunas"
		}

		m := g.Mahasiswa

		result = append(result, map[string]interface{}{
			"id_mahasiswa":   m["id_mahasiswa"],
			"nim":            nim,
			"nama":           m["nama_user"],
			"prodi":          m["prodi"],
			"kelas":          m["kelas"],
			"semester":       m["semester"],
			"total_pinjaman": g.Total,
			"pinjaman_aktif": aktif,
			"status":         status,
		})
	}

	total := len(result)
	start := (page - 1) * perPage
	end := start + perPage
	if start > total {
		start = total
	}
	if end > total {
		end = total
	}

	return result[start:end], total, nil
}

func FetchLoanDetailByMember(memberID string) (map[string]interface{}, error) {
	// 1. Ambil semua loan dari OPAC
	opacLoanURL := fmt.Sprintf("http://localhost:8080/api/loan/%s", memberID)
	resp, err := http.Get(opacLoanURL)
	if err != nil {
		return nil, fmt.Errorf("gagal ambil data loan")
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	var opacResult struct {
		PinjamanAktif []map[string]interface{} `json:"pinjaman_aktif"`
		Riwayat       []map[string]interface{} `json:"riwayat_pinjaman"`
		Status        string                   `json:"status"`
	}
	_ = json.Unmarshal(body, &opacResult)

	if opacResult.PinjamanAktif == nil {
		opacResult.PinjamanAktif = []map[string]interface{}{}
	}
	if opacResult.Riwayat == nil {
		opacResult.Riwayat = []map[string]interface{}{}
	}

	// 2. Hitung total pinjaman
	totalPinjaman := len(opacResult.PinjamanAktif) + len(opacResult.Riwayat)

	// 3. Hitung pinjaman aktif
	pinjamanAktif := len(opacResult.PinjamanAktif)

	// 4. Ambil data mahasiswa
	var mahasiswa map[string]interface{}
	sikompenURL := fmt.Sprintf("http://localhost:8000/api/mahasiswa?nim=%s", memberID)
	respMhs, err := http.Get(sikompenURL)
	if err == nil {
		defer respMhs.Body.Close()
		mhsBytes, _ := io.ReadAll(respMhs.Body)
		var arr []map[string]interface{}
		_ = json.Unmarshal(mhsBytes, &arr)
		for _, m := range arr {
			if fmt.Sprintf("%v", m["kode_user"]) == memberID {
				mahasiswa = m
				break
			}
		}
	}

	nama, kelas, prodi, semester := "", "", "", ""
	if mahasiswa != nil {
		if v, ok := mahasiswa["nama_user"].(string); ok {
			nama = v
		}
		if v, ok := mahasiswa["kelas"].(string); ok {
			kelas = v
		}
		if v, ok := mahasiswa["prodi"].(string); ok {
			prodi = v
		}
		if v, ok := mahasiswa["semester"].(string); ok {
			semester = v
		}
	}

	// 5. Susun response final
	status := "Lunas"
	if pinjamanAktif > 0 {
		status = "Belum Lunas"
	}

	result := map[string]interface{}{
		"nim":             memberID,
		"nama":            nama,
		"kelas":           kelas,
		"prodi":           prodi,
		"semester":        semester,
		"total_pinjaman":  totalPinjaman,
		"pinjaman_aktif":  pinjamanAktif,
		"status":          status,
		"detail_pinjaman": append(opacResult.PinjamanAktif, opacResult.Riwayat...),
	}

	return result, nil
}

func FetchLoanByMemberID(memberID string) (map[string]interface{}, error) {
	if memberID == "" {
		return nil, fmt.Errorf("member_id kosong")
	}

	// ======================
	// AMBIL DATA OPAC (SOURCE OF TRUTH)
	// ======================
	opacURL := fmt.Sprintf("http://localhost:8080/api/loan/%s", memberID)
	resp, err := http.Get(opacURL)
	if err != nil {
		return nil, fmt.Errorf("gagal ambil data loan dari OPAC")
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("gagal baca response OPAC")
	}

	var opacResult map[string]interface{}
	if err := json.Unmarshal(body, &opacResult); err != nil {
		return nil, fmt.Errorf("gagal decode response OPAC")
	}

	// ======================
	// AMBIL MAHASISWA DARI CACHE
	// ======================
	mhsMap := getMahasiswaCache()
	mhs := mhsMap[memberID]

	// ======================
	// SUSUN RESPONSE FINAL
	// ======================
	result := map[string]interface{}{
		"pinjaman_aktif":   opacResult["pinjaman_aktif"],
		"riwayat_pinjaman": opacResult["riwayat_pinjaman"],
		"status":           opacResult["status"],
	}

	if mhs != nil {
		result["mahasiswa"] = map[string]interface{}{
			"nama":     mhs["nama_user"],
			"nim":      mhs["kode_user"],
			"prodi":    mhs["prodi"],
			"kelas":    mhs["kelas"],
			"semester": mhs["semester"],
		}
	}

	return result, nil
}
