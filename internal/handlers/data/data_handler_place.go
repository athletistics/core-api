// Brief   : Database readers for countries
// Author  : Nicolas Baeriswyl <nicolasfabrice.baeriswyl@edu.hefr.ch>
// Created : 27.06.2023

package dataHandler

import (
	"github.com/athletistics/core-api/database"
	"github.com/athletistics/core-api/internal/models"
)

// SelectCountries reads all countries stored in the database.
func SelectCountries() ([]models.Country, error) {
	db := database.DB

	var countries []models.Country
	err := db.Find(&countries).Error

	return countries, err
}
