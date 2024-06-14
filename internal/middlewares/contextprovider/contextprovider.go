// Brief   : Loads the requesting user's context into fiber's "Locals" variable. Aims to reduce overall database
//			 requests.
// Author  : Nicolas Baeriswyl <nicolasfabrice.baeriswyl@edu.hefr.ch>
// Created : 13.06.2024

package contextprovider

import (
	"github.com/athletistics/core-api/constants"
	"github.com/athletistics/core-api/helpers"
	"github.com/gofiber/fiber/v2"
)

func New() fiber.Handler {
	return func(c *fiber.Ctx) error {
		c.Locals(constants.LocalsContext, helpers.LoadContext(c))

		return c.Next()
	}
}
