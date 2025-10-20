package models

import "time"

type OpacMember struct {
	MemberID        string    `gorm:"primaryKey;column:member_id" json:"member_id"`
	MemberName      string    `gorm:"column:member_name" json:"member_name"`
	Gender          string    `gorm:"column:gender" json:"gender"`
	MemberTypeID    string    `gorm:"column:member_type_id" json:"member_type_id"`
	MemberAddress   string    `gorm:"column:member_address" json:"member_address"`
	PostalCode      string    `gorm:"column:postal_code" json:"postal_code"`
	InstName        string    `gorm:"column:inst_name" json:"inst_name"`
	ProgramStudi    string    `gorm:"column:program_studi" json:"program_studi"`
	MstName         string    `gorm:"column:mst_name" json:"mst_name"`
	MemberEmail     string    `gorm:"column:member_email" json:"member_email"`
	MPassword       string    `gorm:"column:mpassword" json:"mpassword"`
	LastLogin       time.Time `gorm:"column:last_login" json:"last_login"`
	InputDate       time.Time `gorm:"column:input_date" json:"input_date"`
	LastUpdate      time.Time `gorm:"column:last_update" json:"last_update"`
}

func (OpacMember) TableName() string {
	return "opac_member"
}
