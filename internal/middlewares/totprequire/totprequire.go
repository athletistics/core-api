// Brief   : Middleware verifying that the second factor was verified recently
// Author  : Nicolas Baeriswyl <nicolasfabrice.baeriswyl@edu.hefr.ch>
// Created : 18.06.2023

package totprequire

import (
	"github.com/gofiber/fiber/v2"
	"gitlab.forge.hefr.ch/tb-athletistics/backend/constants"
	"gitlab.forge.hefr.ch/tb-athletistics/backend/helpers"
	res "gitlab.forge.hefr.ch/tb-athletistics/backend/internal/responses"
	"net/http"
)

const (
	AuthHeaderName = "Authorization"
)

func New() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// get the auth header
		if !c.Locals(constants.LocalsHasVerified2FA).(bool) {
			helpers.LogSecurity("2FA not verified", helpers.GetUserId(c))
			return c.Status(http.StatusUnauthorized).JSON(res.ErrMissingTotp)
		}

		return c.Next()
	}
}
