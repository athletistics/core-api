// Brief   : Middleware providing the user object as part of the fiber context
// Author  : Nicolas Baeriswyl <nicolasfabrice.baeriswyl@edu.hefr.ch>
// Created : 21.06.2023

package userprovider

import (
	"github.com/athletistics/core-api/constants"
	"github.com/athletistics/core-api/database"
	"github.com/athletistics/core-api/helpers"
	"github.com/athletistics/core-api/internal/models"
	res "github.com/athletistics/core-api/internal/responses"
	"github.com/gofiber/fiber/v2"
	"net/http"
)

func New() fiber.Handler {
	return func(c *fiber.Ctx) error {
		user := models.User{ID: helpers.GetUserId(c)}
		if err := database.DB.First(&user).Error; err != nil {
			return c.Status(http.StatusNotFound).JSON(res.ErrDB.Err(err))
		}

		c.Locals(constants.LocalsUser, user)

		return c.Next()
	}
}
