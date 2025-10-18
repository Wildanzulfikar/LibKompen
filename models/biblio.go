package models

import "time"

type OpacBiblio struct {
	BiblioID    int       `gorm:"primaryKey" json:"biblio_id"`
	Title       string    `json:"title"`
	Edition     string    `json:"edition"`
	ISBN        string    `json:"isbn_issn"`
	PublisherID int       `json:"publisher_id"`
	PublishYear string    `json:"publish_year"`
	Collation   string    `json:"collation"`
	LanguageID  string    `json:"language_id"`
	CallNumber  string    `json:"call_number"`
	Notes       string    `json:"notes"`
	Image       string    `json:"image"`
	InputDate   time.Time `json:"input_date"`
	LastUpdate  time.Time `json:"last_update"`
}
