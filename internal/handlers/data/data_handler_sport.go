// Brief   : Database readers for sports and event types
// Author  : Nicolas Baeriswyl <nicolasfabrice.baeriswyl@edu.hefr.ch>
// Created : 27.06.2023

package dataHandler

import (
	"github.com/athletistics/core-api/database"
	"github.com/athletistics/core-api/internal/models"
)

func SelectSports() ([]models.Sport, error) {
	db := database.DB

	var sports []models.Sport
	err := db.Find(&sports).Error

	return sports, err
}

func SelectCategoriesBySport(sportId int) ([]models.Category, error) {
	db := database.DB

	var categories []models.Category
	err := db.Where(&models.Category{SportID: sportId}).Find(&categories).Error

	return categories, err
}

func SelectEventTypes(filter bool, requesterRoles []models.Authorisation) ([]models.EventType, error) {
	db := database.DB

	var eventTypes []models.EventType
	var err error

	// if filter is false, return all event types
	if !filter {
		err = db.Find(&eventTypes).Error
		return eventTypes, err
	}

	if len(requesterRoles) == 0 {
		return nil, nil
	}

	roles := make([]int, len(requesterRoles))
	for i := 0; i < len(requesterRoles); i++ {
		roles[i] = requesterRoles[i].AccessRoleID
	}

	err = db.
		Joins(`INNER JOIN event_type_access_roles r ON "r"."event_type_id" = "event_type"."id" AND "r"."access_role_id" IN ?`, roles).
		Distinct("id", "i18n_name").
		Find(&eventTypes).Error

	return eventTypes, err
}
