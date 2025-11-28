package services

import (
	"LibKompen/database"
	"LibKompen/models"
)

func GetAllApprovals() ([]models.ApprovalBebasPustaka, error) {
	var approvals []models.ApprovalBebasPustaka
	result := database.DB.Find(&approvals)
	return approvals, result.Error
}

func CreateApproval(approval *models.ApprovalBebasPustaka) error {
	result := database.DB.Create(approval)
	return result.Error
}

func GetLatestApprovalByKodeUser(kodeUser string) (*models.ApprovalBebasPustaka, error) {
	var approval models.ApprovalBebasPustaka
	err := database.DB.Table("status_approval").Where("kode_user = ?", kodeUser).Order("created_at desc").First(&approval).Error
	if err != nil {
		return nil, err
	}
	return &approval, nil
}

func GetUserByID(id uint) (*models.UsersBebasPustaka, error) {
	var user models.UsersBebasPustaka
	err := database.DB.Table("users").Where("id_users = ?", id).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func FindApprovalByID(id uint) (*models.ApprovalBebasPustaka, error) {
	var approval models.ApprovalBebasPustaka
	result := database.DB.First(&approval, id)
	if result.Error != nil {
		return nil, result.Error
	}
	return &approval, nil
}

func DeleteApproval(approval *models.ApprovalBebasPustaka) error {
	return database.DB.Delete(approval).Error
}
