// Brief   :
// Author  : Nicolas Baeriswyl <nicolasfabrice.baeriswyl@edu.hefr.ch>
// Created : 27.06.2023

package dataHandler

import "github.com/gofiber/fiber/v2"

// GetCategories returns a list of all categories registered for the authenticated user's sport.
//
//	@Summary	Get a list of all categories for a specific sport.
//	@Tags		Data
//	@Produce	json
//	@Success	200	{object}	[]models.Category
//	@Failure	401	{object}	res.Error
//	@Security	ApiKeyAuth
//	@Router		/data/categories [get]
func GetCategories(c *fiber.Ctx) error {
	return nil
}

// GetPlaces returns a list of places based on at least two characters for a search query. The characters can either be
// part of the postal code, or the name of the place.
//
//	@Summary	Searches for places.
//	@Tags		Data
//	@Produce	json
//
//	@Param		country	query		string	true	"Country ID"
//	@Param		state	query		string	false	"State ID"
//	@Param		q		query		string	true	"Search query, postal code or name of the place"
//
//	@Success	200		{object}	[]models.PlaceOut
//	@Failure	401		{object}	res.Error
//	@Security	ApiKeyAuth
//	@Router		/data/places [get]
func GetPlaces(c *fiber.Ctx) error {
	return nil
}

// GetStates returns a list of all states of a given country identified by its ID.
//
//	@Summary	Get all states of a country.
//	@Tags		Data
//	@Produce	json
//
//	@Param		country	query		string	true	"Country ID"
//
//	@Success	200		{object}	[]models.State
//	@Failure	401		{object}	res.Error
//	@Security	ApiKeyAuth
//	@Router		/data/places/states [get]
func GetStates(c *fiber.Ctx) error {
	return nil
}

// GetCountries returns a list of all available countries.
//
//	@Summary	Get all countries.
//	@Tags		Data
//	@Produce	json
//	@Success	200	{object}	[]models.Country
//	@Failure	401	{object}	res.Error
//	@Security	ApiKeyAuth
//	@Router		/data/places/countries [get]
func GetCountries(c *fiber.Ctx) error {
	return nil
}

// GetLanguages returns a list of all available languages.
//
//	@Summary	Get all languages.
//	@Tags		Data
//	@Produce	json
//	@Success	200	{object}	[]models.Language
//	@Failure	401	{object}	res.Error
//	@Security	ApiKeyAuth
//	@Router		/data/languages [get]
func GetLanguages(c *fiber.Ctx) error {
	return nil
}

// GetAccessRoles returns a list of all available access roles. This route is only available for head coaches.
//
//	@Summary	Lists all available access roles that can be given to a user.
//	@Tags		Data
//	@Produce	json
//	@Success	200	{object}	[]models.AccessRole
//	@Failure	401	{object}	res.Error
//	@Failure	403	{object}	res.Error
//	@Security	ApiKeyAuth
//	@Router		/data/accessRoles [get]
func GetAccessRoles(c *fiber.Ctx) error {
	return nil
}

// GetEventTypes lists all event types that can be attributed to an event.
//
//	@Summary
//	@Tags		Data
//	@Produce	json
//	@Success	200	{object}	[]models.EventType
//	@Failure	401	{object}	res.Error
//	@Security	ApiKeyAuth
//	@Router		/data/eventTypes [get]
func GetEventTypes(c *fiber.Ctx) error {
	return nil
}

// GetSports lists all currently registered sports.
//
//	@Summary	Get a list of all sports.
//	@Tags		Data
//	@Produce	json
//	@Success	200	{object}	[]models.Sport
//	@Failure	401	{object}	res.Error
//	@Security	ApiKeyAuth
//	@Router		/data/sports [get]
func GetSports(c *fiber.Ctx) error {
	return nil
}

// GetIdTypes lists all currently registered ID types, e.g. Swiss-Athletics license number.
//
//	@Summary	get a list of available ID types.
//	@Tags		Data
//	@Produce	json
//	@Success	200	{object}	[]models.ExternalIdType
//	@Failure	401	{object}	res.Error
//	@Security	ApiKeyAuth
//	@Router		/data/idTypes [get]
func GetIdTypes(c *fiber.Ctx) error {
	return nil
}
