package utils

// Mapping kode prodi berdasarkan struktur NIM
var KodeProdiMap = map[string]string{
	// --- TEKNIK SIPIL (01)
	"0131": "Konstruksi Gedung (D3)",
	"0141": "Teknik Perancangan Jalan dan Jembatan (D4)",
	"0132": "Konstruksi Sipil (D3)",
	"0121": "Konstruksi Gedung (D4)",
	// --- TEKNIK MESIN (02)
	"0231": "Teknik Mesin (D3)",
	"0211": "Teknologi Rekayasa Manufaktur (D4)",
	"0221": "Teknologi Rekayasa Pembangkit Energi (D4)",
	"0232": "Teknologi Rekayasa Konversi Energi (D4)",
	"0233": "Teknologi Rekayasa Konversi Energi (D3/D4)",
	"0243": "Teknologi Rekayasa Konversi Energi (D4)",
	"0234": "Teknologi Rekayasa Pemeliharaan Alat Berat (D4)",
	// --- TEKNIK ELEKTRO (03)
	"0331": "Teknik Listrik (D3)",
	"0311": "Teknik Otomasi Listrik Industri (D4)",
	"0332": "Elektronika Industri (D3)",
	"0321": "Broadband Multimedia (D4)",
	"0333": "Teknik Telekomunikasi (D3)",
	"0334": "Instrumentasi (D4)",
	// --- AKUNTANSI (04)
	"0431": "Akuntansi Keuangan (D3/D4)",
	"0411": "Keuangan dan Perbankan Syariah (D4)",
	"0432": "Keuangan dan Perbankan (D3)",
	"0421": "Keuangan dan Perbankan (D4)",
	"0433": "Akuntansi Keuangan (D4)",
	"0441": "Manajemen Keuangan (D4)",
	// --- ADMINISTRASI NIAGA (05)
	"0511": "BISPRO/MICE (D4)",
	"0521": "Administrasi Bisnis Terapan (D4)",
	"0531": "Administrasi Bisnis Terapan (D3)",
	// --- TEKNIK GRAFIKA PENERBITAN (06)
	"0611": "Teknologi Industri Cetak Kemasan (D4)",
	"0621": "Desain Grafis (D4)",
	"0632": "Penerbitan (D3)",
	// --- TEKNIK INFORMATIKA DAN KOMPUTER (07)
	"0711": "Teknik Informatika (D4)",
	"0712": "Teknik Komputer dan Jaringan (D1)",
	"0721": "Teknik Multimedia Jaringan (D4)",
	"0731": "Teknik Multimedia Digital (D4)",
	// --- PRODI KHUSUS (08)
	"0831": "Manajemen Pemasaran WNBK (D3)",
	"0841": "Bahasa Inggris untuk Komunikasi Bisnis (D4)",
}

// Fungsi ekstrak nama prodi dari NIM
func GetProdiFromNIM(nim string) string {
	if len(nim) < 7 {
		return "Unknown"
	}
	// Cek prodi khusus (08xx)
	if nim[2:4] == "08" && len(nim) >= 6 {
		kode := nim[2:6] // 4 digit setelah tahun masuk
		if nama, ok := KodeProdiMap[kode]; ok {
			return nama
		}
	}
	// Normal: kode jurusan (2), jenjang (1), kode prodi (1)
	kode := nim[2:4] + string(nim[4]) + string(nim[5])
	if nama, ok := KodeProdiMap[kode]; ok {
		return nama
	}
	return "Unknown"
}
