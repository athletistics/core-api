// Brief   :
// Author  : Nicolas Baeriswyl <nicolasfabrice.baeriswyl@edu.hefr.ch>
// Created : 27.06.2023

package dataHandler

import (
	"fmt"
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

// SelectStatesByCountry reads a list of all states of a given country.
func SelectStatesByCountry(countryId int) ([]models.State, error) {
	db := database.DB

	var states []models.State
	err := db.Where(&models.State{CountryID: countryId}).Find(&states).Error

	return states, err
}

// SearchForPlaces reads a list of places for a given country and state.
// The query string is searched in the post code and place name columns.
func SearchForPlaces(query string, countryId int, stateId int) ([]models.PlaceOut, error) {
	db := database.DB

	req := db.Model(&models.Place{}).InnerJoins("State", db.Where(&models.State{ID: stateId, CountryID: countryId}))

	// Apply the query if one is given
	if len(query) > 0 {
		queryWithPercent := fmt.Sprintf("%v%%", query)
		req.Where(db.Where(`"post_code" LIKE ?`, queryWithPercent).Or(`"place" LIKE ?`, queryWithPercent))
	}

	var places []models.PlaceOut
	err := req.Find(&places).Error

	return places, err
}
