package models

import "time"

type TenggatWaktu struct {
	IdTenggatWaktu uint      `gorm:"primaryKey" json:"id_tenggat_waktu"`
	IdUsers        uint      `json:"id_users"`
	WaktuMulai     time.Time `json:"waktu_mulai"`
	WaktuAkhir     time.Time `json:"waktu_akhir"`
	CreatedAt      time.Time `json:"createdAt"`
	UpdatedAt      time.Time `json:"updatedAt"`
}
