package models

import "time"

type Visitor struct {
	Id         uint      `gorm:"primaryKey" json:"id"`
	KodeUser   uint      `json:"kode_user"`
	Created_At time.Time `json:"created_at"`
}
