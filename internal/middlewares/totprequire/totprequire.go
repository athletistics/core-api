// Brief   :
// Author  : Nicolas Baeriswyl <nicolasfabrice.baeriswyl@edu.hefr.ch>
// Created : 18.06.2023

package totprequire

import (
	"github.com/gofiber/fiber/v2"
	"gitlab.forge.hefr.ch/tb-athletistics/backend/constants"
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
			return c.Status(http.StatusUnauthorized).JSON(res.ErrMissingTotp)
		}

		return c.Next()
	}
}
