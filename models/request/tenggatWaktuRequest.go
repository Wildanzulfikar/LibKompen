package request

import "time"

type TenggatWaktuRequest struct {
	IdUsers    uint      `json:"id_users"`
	WaktuMulai time.Time `json:"waktu_mulai"`
	WaktuAkhir time.Time `json:"waktu_akhir"`
}
