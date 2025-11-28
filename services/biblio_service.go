package services

import (
	"LibKompen/database"
	"LibKompen/models"
)

func GetAllBiblio() ([]models.OpacBiblio, error) {
	var biblios []models.OpacBiblio
	result := database.DB.Find(&biblios)
	return biblios, result.Error
}

func CreateBiblio(biblio *models.OpacBiblio) error {
	result := database.DB.Create(biblio)
	return result.Error
}
