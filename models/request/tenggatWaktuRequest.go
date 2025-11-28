package request

type TenggatWaktuRequest struct {
	IdUsers    uint   `json:"id_users"`
	WaktuMulai string `json:"waktu_mulai"`
	WaktuAkhir string `json:"waktu_akhir"`
}
