// Brief   : External ID management handlers
// Author  : Nicolas Baeriswyl <nicolasfabrice.baeriswyl@edu.hefr.ch>
// Created : 06.07.2023

package externalId

import (
	"github.com/gofiber/fiber/v2"
	"gitlab.forge.hefr.ch/tb-athletistics/backend/database"
	"gitlab.forge.hefr.ch/tb-athletistics/backend/helpers"
	"gitlab.forge.hefr.ch/tb-athletistics/backend/internal/models"
	res "gitlab.forge.hefr.ch/tb-athletistics/backend/internal/responses"
	"net/http"
)

const (
	QueryInternalId = "internalId"
)

// SelectExternalUserIds reads the list of all external IDs of a specific user.
//
//	@Summary	List all external IDs of a specific user.
//	@Tags		External ID
//	@Produce	json
//	@Param		internalId	query		string	true	"Internal User identifier"
//	@Success	200			{object}	[]models.ExternalId
//	@Failure	401			{object}	res.Error
//	@Failure	403			{object}	res.Error	"Not enough permissions"
//	@Security	ApiKeyAuth
//	@Router		/externalId/user [get]
func SelectExternalUserIds(c *fiber.Ctx) error {
	// parse request
	userId := c.Query(QueryInternalId)
	user := helpers.GetUser(c)

	data, err := selectExternalUserIds(&user, userId)
	return helpers.ProcessJsonResponse(c, data, err)
}

// CreateExternalUserId creates a new external ID. The internal ID must belong to a user.
//
//	@Summary	Creates a new external ID. The internal ID must belong to a user.
//	@Tags		External ID
//	@Consume	json
//	@Produce	json
//	@Param		request	body		models.ExternalIdInput	true	"External ID"
//	@Success	200		{string}	string					"OK"
//	@Failure	401		{object}	res.Error
//	@Failure	403		{object}	res.Error	"Not enough permissions"
//	@Security	ApiKeyAuth
//	@Router		/externalId/user [post]
func CreateExternalUserId(c *fiber.Ctx) error {
	// parse request
	b := new(models.ExternalIdInput)
	parseErr := c.BodyParser(b)
	user := helpers.GetUser(c)

	err := createExternalUserId(&user, b, parseErr)
	return helpers.ProcessJsonResponse(c, nil, err)
}

// GetTypes lists all available ID types.
//
//	@Summary	List all available ID types.
//	@Tags		External ID
//	@Produce	json
//	@Success	200	{object}	[]models.ExternalIdType
//	@Security	ApiKeyAuth
//	@Router		/externalId/idTypes [get]
func GetTypes(c *fiber.Ctx) error {
	db := database.DB

	var types []models.ExternalIdType
	if err := db.Find(&types).Error; err != nil {
		return c.Status(http.StatusInternalServerError).JSON(res.NewErrorFromDb(err))
	}

	return c.JSON(types)
}

// selectExternalUserIds reads a list of external IDs. The internal ID must belong to a user that is
// either the requesting user themself, or to a club that the requesting user is a head coach or coach.
func createExternalUserId(user *models.User, body *models.ExternalIdInput, err error) *res.Error {
	db := database.DB
	if err != nil {
		return res.ErrBadRequest.Err(err)
	}

	// make sure user can be accessed
	targetUser := models.User{ID: body.InternalID}
	if err := db.First(&targetUser).Error; err != nil {
		return res.NewErrorFromDb(err)
	}
	if user.ID != body.InternalID && !helpers.HasRoleOnClub(user, targetUser.ClubID, helpers.RoleHeadCoach) {
		return res.ErrForbidden
	}

	// create the id
	id := models.ExternalId{
		ExternalID:       body.ExternalID,
		InternalID:       body.InternalID,
		ExternalIdTypeID: body.ExternalIdTypeID,
	}
	if err := db.Create(&id).Error; err != nil {
		return res.NewErrorFromDb(err)
	}

	return nil
}

// selectExternalUserIds reads a list of external IDs. The internal ID must belong to a user that is
// either the requesting user themself, or to a club that the requesting user is a head coach or coach,
// or the requesting user is head coach or coach on one of the clubs the target user has autorisations.
func selectExternalUserIds(user *models.User, internalId string) ([]models.ExternalId, *res.Error) {
	db := database.DB

	// make sure user can be accessed
	targetUser := models.User{ID: internalId}
	if err := db.First(&targetUser).Error; err != nil {
		return nil, res.NewErrorFromDb(err)
	}
	if user.ID != internalId && !helpers.HasRoleOnOneOfTheClubs(user, &targetUser, helpers.RoleHeadCoach, helpers.RoleCoach) {
		return nil, res.ErrForbidden
	}

	// select all external ids
	var types []models.ExternalId
	if err := db.Where(&models.ExternalId{InternalID: internalId}).Find(&types).Error; err != nil {
		return nil, res.NewErrorFromDb(err)
	}

	return types, nil
}
