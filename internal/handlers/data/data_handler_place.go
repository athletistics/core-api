// Brief   : Database readers for countries
// Author  : Nicolas Baeriswyl <nicolasfabrice.baeriswyl@edu.hefr.ch>
// Created : 27.06.2023

package dataHandler

import (
	"gitlab.forge.hefr.ch/tb-athletistics/backend/database"
	"gitlab.forge.hefr.ch/tb-athletistics/backend/internal/models"
)

// SelectCountries reads all countries stored in the database.
func SelectCountries() ([]models.Country, error) {
	db := database.DB

	var countries []models.Country
	err := db.Find(&countries).Error

	return countries, err
}
