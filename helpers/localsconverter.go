// Brief   : Converting Fiber Context locals to their respective values
// Author  : Nicolas Baeriswyl <nicolasfabrice.baeriswyl@edu.hefr.ch>
// Created : 17.06.2023

package helpers

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"gitlab.forge.hefr.ch/tb-athletistics/backend/constants"
	"gitlab.forge.hefr.ch/tb-athletistics/backend/internal/models"
)

func GetUserId(c *fiber.Ctx) string {
	return fmt.Sprintf("%v", c.Locals(constants.LocalsUserId))
}

func GetUser(c *fiber.Ctx) models.User {
	return c.Locals(constants.LocalsUser).(models.User)
}
