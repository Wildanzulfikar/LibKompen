package models

import "time"

type VisitorSummary struct {
	ID        int64     `json:"id" gorm:"primaryKey;autoIncrement"`
	KodeUser  string    `json:"kode_user" gorm:"column:kode_user;type:varchar(50)"`
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
}

func (VisitorSummary) TableName() string {
	return "visitor_summary"
}
