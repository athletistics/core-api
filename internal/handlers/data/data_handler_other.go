// Brief   : Database readers for remaining data
// Author  : Nicolas Baeriswyl <nicolasfabrice.baeriswyl@edu.hefr.ch>
// Created : 27.06.2023

package dataHandler

import (
	"gitlab.forge.hefr.ch/tb-athletistics/backend/database"
	"gitlab.forge.hefr.ch/tb-athletistics/backend/internal/models"
)

func SelectIdTypes() ([]models.ExternalIdType, error) {
	db := database.DB

	var types []models.ExternalIdType
	err := db.Find(&types).Error

	return types, err
}

func SelectAccessRoles() ([]models.AccessRole, error) {
	db := database.DB

	var accessRoles []models.AccessRole
	err := db.Find(&accessRoles).Error

	return accessRoles, err
}

func SelectLanguages() ([]models.Language, error) {
	db := database.DB

	var languages []models.Language
	err := db.Find(&languages).Error

	return languages, err
}
