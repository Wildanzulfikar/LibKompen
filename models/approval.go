package models

import "time"


type ApprovalBebasPustaka struct {
	IDStatus       int       `gorm:"primaryKey;column:id_status" json:"id_status"`
	KodeUser       string    `json:"kode_user"`
	IDUsers        uint      `json:"id_users"`
	StatusApproval bool      `json:"status_approval"`
	Keterangan     string    `json:"keterangan"`
	CreatedAt      time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt      time.Time `gorm:"column:updated_at" json:"updated_at"`
}

func (ApprovalBebasPustaka) TableName() string {
    return "status_approval"
}
