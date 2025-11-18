package models

type StatusApproval struct {
	IdStatus       uint   `gorm:"primaryKey" json:"id_status"`
	KodeUser       uint   `json:"kode_user"`
	IdUser         uint   `json:"id_user"`
	StatusApproval bool   `json:"status_approval"`
	Keterangan     string `json:"keterangan"`
}
