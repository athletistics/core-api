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

func GetMe(c *fiber.Ctx) error {
	db := database.DB
	userId := helpers.GetUserId(c)

	var me models.User
	err := db.Where(&models.User{Id: userId}).First(&me).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return c.Status(http.StatusNotFound).JSON(res.ErrNotFound)
	} else if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(res.ErrDB)
	}

	return c.JSON(me)
}

func GetAuth(c *fiber.Ctx) error {
	type role struct {
		club string
		role string
	}
	type response struct {
		authRoles []role
		is2FaOk   bool
	}

	var auth []models.Authorisation
	err := database.DB.Where(&models.Authorisation{UserID: helpers.GetUserId(c)}).Find(&auth)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(res.ErrDB)
	}

	roles := make([]role, len(auth))
	for i := 0; i < len(auth); i++ {
		roles[i] = role{club: auth[i].Club.Name, role: auth[i].AccessRole.I18nName}
	}

	resp := response{
		authRoles: roles,
		is2FaOk:   c.Locals(constants.LocalsHasVerified2FA).(bool),
	}
	return c.JSON(resp)
}
