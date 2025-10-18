package models

import "time"

type ApprovalBebasPustaka struct {
	ID             uint      `gorm:"primaryKey" json:"id"`
	KodeUser       string    `json:"kode_user"`
	NamaUser       string    `json:"nama_user"`
	Prodi          string    `json:"prodi"`
	Kelas          string    `json:"kelas"`
	StatusOpac     bool      `json:"status_opac"`
	StatusApproval bool      `json:"status_approval"`
	Keterangan     string    `json:"keterangan"`
	ApprovedBy     uint      `json:"approved_by"`
	ApprovedAt     time.Time `json:"approved_at"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}
