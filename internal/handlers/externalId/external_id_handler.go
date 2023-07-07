// Brief   :
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

func SelectExternalUserIds(c *fiber.Ctx) error {
	// parse request
	userId := c.Query(QueryInternalId)
	user := helpers.GetUser(c)

	data, err := selectExternalUserIds(&user, userId)
	return helpers.ProcessJsonResponse(c, data, err)
}

func CreateExternalUserId(c *fiber.Ctx) error {
	// parse request
	b := new(models.ExternalIdInput)
	parseErr := c.BodyParser(b)
	user := helpers.GetUser(c)

	err := createExternalUserId(&user, b, parseErr)
	return helpers.ProcessJsonResponse(c, nil, err)
}

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
