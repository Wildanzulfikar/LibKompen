package models

import "time"

type OpacLoan struct {
	LoanID      uint      `gorm:"primaryKey;column:loan_id" json:"loan_id"`
	ItemCode    string    `gorm:"column:item_code" json:"item_code"`
	MemberID    string    `gorm:"column:member_id" json:"member_id"`
	LoanDate    time.Time `gorm:"column:loan_date" json:"loan_date"`
	DueDate     time.Time `gorm:"column:due_date" json:"due_date"`
	ReturnDate  time.Time `gorm:"column:return_date" json:"return_date"`
	IsReturn    bool      `gorm:"column:is_return" json:"is_return"`
	BiblioID    uint      `gorm:"column:biblio_id" json:"biblio_id"`
}

func (OpacLoan) TableName() string {
	return "opac_loan"
}
