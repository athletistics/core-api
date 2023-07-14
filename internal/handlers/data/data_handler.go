// Brief   : Data reading handlers
// Author  : Nicolas Baeriswyl <nicolasfabrice.baeriswyl@edu.hefr.ch>
// Created : 27.06.2023

package dataHandler

import (
	"github.com/gofiber/fiber/v2"
	"gitlab.forge.hefr.ch/tb-athletistics/backend/helpers"
	res "gitlab.forge.hefr.ch/tb-athletistics/backend/internal/responses"
	"net/http"
)

const (
	QuerySportId = "sport"
	QueryFilter  = "filter"
	Query        = "q"
)

// GetCategories returns a list of all categories registered for the authenticated user's sport.
//
//	@Summary	Get a list of all categories for a specific sport.
//	@Tags		Data
//	@Produce	json
//	@Param		sport	query		integer	true	"Sport ID"
//	@Success	200		{object}	[]models.Category
//	@Failure	401		{object}	res.Error
//	@Security	ApiKeyAuth
//	@Router		/data/categories [get]
func GetCategories(c *fiber.Ctx) error {
	sportId := c.QueryInt(QuerySportId)
	if sportId == 0 {
		helpers.LogInfo("bad request: missing parameter 'sportId'", helpers.GetUserId(c))
		return c.Status(http.StatusBadRequest).JSON(res.ErrBadRequest)
	}

	result, err := SelectCategoriesBySport(sportId)
	if err != nil {
		helpers.LogError(err.Error(), helpers.GetUserId(c))
		return c.Status(http.StatusInternalServerError).JSON(res.ErrDB.Err(err))
	}

	return c.JSON(result)
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
	result, err := SelectCountries()
	if err != nil {
		helpers.LogError(err.Error(), helpers.GetUserId(c))
		return c.Status(http.StatusInternalServerError).JSON(res.ErrDB.Err(err))
	}

	return c.JSON(result)
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
	result, err := SelectLanguages()
	if err != nil {
		helpers.LogError(err.Error(), helpers.GetUserId(c))
		return c.Status(http.StatusInternalServerError).JSON(res.ErrDB.Err(err))
	}

	return c.JSON(result)
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
	user := helpers.GetUser(c)
	if !helpers.HasRoleOnParentClub(&user, helpers.RoleHeadCoach) {
		helpers.LogSecurity("tried accessing access roles", helpers.GetUserId(c))
		return c.Status(http.StatusForbidden).JSON(res.ErrForbidden)
	}

	result, err := SelectAccessRoles()
	if err != nil {
		helpers.LogError(err.Error(), helpers.GetUserId(c))
		return c.Status(http.StatusInternalServerError).JSON(res.ErrDB.Err(err))
	}

	return c.JSON(result)
}

// GetEventTypes lists all event types that can be attributed to an event. The returned list can be filtered
// to only contain elements that can be used by the authenticated user to create a new event.
//
//	@Summary	Get a list of all event types.
//	@Tags		Data
//	@Produce	json
//
//	@Param		filter	query		boolean	false	"Filter, returning only elements that can be used by authenticated user"
//
//	@Success	200		{object}	[]models.EventType
//	@Failure	401		{object}	res.Error
//	@Security	ApiKeyAuth
//	@Router		/data/eventTypes [get]
func GetEventTypes(c *fiber.Ctx) error {
	user := helpers.GetUser(c)

	filter := c.QueryBool(QueryFilter, true)
	roles, err := helpers.GetAuthorisations(user.ID, user.ClubID)
	if err != nil {
		helpers.LogError(err.Error(), helpers.GetUserId(c))
		return c.Status(http.StatusInternalServerError).JSON(res.ErrDB.Err(err))
	}

	result, err := SelectEventTypes(filter, roles)
	if err != nil {
		helpers.LogError(err.Error(), helpers.GetUserId(c))
		return c.Status(http.StatusInternalServerError).JSON(res.ErrDB.Err(err))
	}

	return c.JSON(result)
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
	result, err := SelectSports()
	if err != nil {
		helpers.LogError(err.Error(), helpers.GetUserId(c))
		return c.Status(http.StatusInternalServerError).JSON(res.ErrDB.Err(err))
	}

	return c.JSON(result)
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
	result, err := SelectIdTypes()
	if err != nil {
		helpers.LogError(err.Error(), helpers.GetUserId(c))
		return c.Status(http.StatusInternalServerError).JSON(res.ErrDB.Err(err))
	}

	return c.JSON(result)
}
