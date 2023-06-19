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
	type Role struct {
		Club string `json:"club"`
		Role string `json:"role"`
	}
	type Response struct {
		AuthRoles []Role `json:"authRoles"`
		Is2FaOk   bool   `json:"is2FaOk"`
		Has2Fa    bool   `json:"has2Fa"`
	}

	var auth []models.Authorisation
	err := database.DB.Where(&models.Authorisation{UserID: helpers.GetUserId(c)}).Preload("AccessRole").Preload("Club").Find(&auth).Error
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(res.ErrDB)
	}

	roles := make([]Role, len(auth))
	for i := 0; i < len(auth); i++ {
		roles[i] = Role{Club: auth[i].Club.Name, Role: auth[i].AccessRole.I18nName}
	}

	resp := Response{
		AuthRoles: roles,
		Is2FaOk:   c.Locals(constants.LocalsHasVerified2FA).(bool),
		Has2Fa:    c.Locals(constants.LocalsHas2FA).(bool),
	}
	return c.JSON(resp)
}
