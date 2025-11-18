package request

import "time"

type TenggatWaktu struct {
	WaktuMulai time.Time `json:"waktu_mulai"`
	WaktuAkhir time.Time `json:"waktu_akhir"`
}
