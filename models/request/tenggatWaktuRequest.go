package request

import "time"

type TenggatWaktuRequest struct {
	WaktuMulai time.Time `json:"waktu_mulai"`
	WaktuAkhir time.Time `json:"waktu_akhir"`
}
