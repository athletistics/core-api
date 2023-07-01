// Brief   :
// Author  : Nicolas Baeriswyl <nicolasfabrice.baeriswyl@edu.hefr.ch>
// Created : 18.06.2023

package meHandler

import (
	"errors"
	"github.com/gofiber/fiber/v2"
	"gitlab.forge.hefr.ch/tb-athletistics/backend/constants"
	"gitlab.forge.hefr.ch/tb-athletistics/backend/database"
	"gitlab.forge.hefr.ch/tb-athletistics/backend/helpers"
	"gitlab.forge.hefr.ch/tb-athletistics/backend/internal/models"
	res "gitlab.forge.hefr.ch/tb-athletistics/backend/internal/responses"
	"gorm.io/gorm"
	"net/http"
)

// GetMe gets details of the authenticated user
//
//	@Summary	Detailed information about the connected user
//	@Tags		Personal data
//	@Produce	json
//	@Success	200	{object}	models.User
//	@Failure	401	{object}	res.Error
//	@Security	ApiKeyAuth
//	@Router		/me [get]
func GetMe(c *fiber.Ctx) error {
	db := database.DB
	userId := helpers.GetUserId(c)

	me := models.User{ID: userId}
	err := db.
		Preload("Nationality").Preload("Country").Preload("Club").
		First(&me).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return c.Status(http.StatusNotFound).JSON(res.ErrNotFound)
	} else if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(res.ErrDB)
	}

	return c.JSON(me)
}

// GetAuth lists all roles of the authenticated user, as well as the state of the second factor.
//
//	@Summary	Get authentication status
//	@Tags		Personal data
//	@Produce	json
//	@Success	200	{object}	models.UserResponse
//	@Failure	401	{object}	res.Error
//	@Security	ApiKeyAuth
//	@Router		/me/auth [get]
func GetAuth(c *fiber.Ctx) error {
	var auth []models.Authorisation
	err := database.DB.Where(&models.Authorisation{UserID: helpers.GetUserId(c)}).Preload("AccessRole").Preload("Club").Find(&auth).Error
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(res.ErrDB)
	}

	roles := make([]models.UserRole, len(auth))
	for i := 0; i < len(auth); i++ {
		roles[i] = models.UserRole{Club: auth[i].Club.Name, Role: auth[i].AccessRole.I18nName}
	}

	resp := models.UserResponse{
		AuthRoles: roles,
		Is2FaOk:   c.Locals(constants.LocalsHasVerified2FA).(bool),
		Has2Fa:    c.Locals(constants.LocalsHas2FA).(bool),
	}
	return c.JSON(resp)
}
