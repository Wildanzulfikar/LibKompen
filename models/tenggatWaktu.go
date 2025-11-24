package models

import "time"

type TenggatWaktu struct {
	IdTenggatWaktu uint              `gorm:"primaryKey" json:"id_tenggat_waktu"`
	IdUsers        uint              `json:"id_users"`
	User           UsersBebasPustaka `gorm:"foreignKey:IdUsers" json:"users_bebas_pustaka"`
	WaktuMulai     time.Time         `json:"waktu_mulai"`
	WaktuAkhir     time.Time         `json:"waktu_akhir"`
	CreatedAt      time.Time         `gorm:"column:createdAt" json:"createdAt"`
	UpdatedAt      time.Time         `gorm:"column:updatedAt" json:"updatedAt"`
}

func (TenggatWaktu) TableName() string {
	return "tenggat_waktu"
}
