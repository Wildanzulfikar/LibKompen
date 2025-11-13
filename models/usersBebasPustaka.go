package models

import "time"

type UsersBebasPustaka struct {
	IdUsers        uint      `gorm:"primaryKey" json:"id_users"`
	IdTenggatWaktu time.Time `json:"id_tenggat_waktu"`
	Name           string    `json:"name"`
	Email          string    `json:"email"`
	Username       string    `json:"username"`
	Password       string    `json:"password"`
	Role           string    `json:"role"`
	Status         bool      `json:"status"`
	LastLogin      time.Time `json:"last_login"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

func (UsersBebasPustaka) TableName() string {
	return "users_bebas_pustaka"
}
