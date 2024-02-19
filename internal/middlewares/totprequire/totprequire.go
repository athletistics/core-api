// Brief   : Middleware verifying that the second factor was verified recently
// Author  : Nicolas Baeriswyl <nicolasfabrice.baeriswyl@edu.hefr.ch>
// Created : 18.06.2023

package totprequire

import (
	"github.com/athletistics/core-api/constants"
	"github.com/athletistics/core-api/helpers"
	res "github.com/athletistics/core-api/internal/responses"
	"github.com/gofiber/fiber/v2"
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
