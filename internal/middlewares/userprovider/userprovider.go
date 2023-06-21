// Brief   :
// Author  : Nicolas Baeriswyl <nicolasfabrice.baeriswyl@edu.hefr.ch>
// Created : 21.06.2023

package userprovider

import (
	"github.com/gofiber/fiber/v2"
	"gitlab.forge.hefr.ch/tb-athletistics/backend/constants"
	"gitlab.forge.hefr.ch/tb-athletistics/backend/database"
	"gitlab.forge.hefr.ch/tb-athletistics/backend/helpers"
	"gitlab.forge.hefr.ch/tb-athletistics/backend/internal/models"
	res "gitlab.forge.hefr.ch/tb-athletistics/backend/internal/responses"
	"net/http"
)

func New() fiber.Handler {
	return func(c *fiber.Ctx) error {
		user := models.User{ID: helpers.GetUserId(c)}
		if err := database.DB.First(&user).Error; err != nil {
			return c.Status(http.StatusNotFound).JSON(res.ErrDB.Err(err))
		}

		c.Locals(constants.LocalsUser)

		return c.Next()
	}
}
