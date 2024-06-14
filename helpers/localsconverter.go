// Brief   : Converting Fiber Context locals to their respective values
// Author  : Nicolas Baeriswyl <nicolasfabrice.baeriswyl@edu.hefr.ch>
// Created : 17.06.2023

package helpers

import (
	"fmt"
	"github.com/athletistics/core-api/constants"
	"github.com/athletistics/core-api/internal/models"
	"github.com/gofiber/fiber/v2"
)

func GetUserId(c *fiber.Ctx) string {
	return fmt.Sprintf("%v", c.Locals(constants.LocalsUserId))
}

// GetUser reads the user from fiber's context
//
// Deprecated: use UserContext instead.
func GetUser(c *fiber.Ctx) models.User {
	return c.Locals(constants.LocalsUser).(models.User)
}
