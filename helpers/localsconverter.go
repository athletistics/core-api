// Brief   : Converting Fiber Context locals to their respective values
// Author  : Nicolas Baeriswyl <nicolasfabrice.baeriswyl@edu.hefr.ch>
// Created : 17.06.2023

package helpers

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"gitlab.forge.hefr.ch/tb-athletistics/backend/constants"
)

func GetUserId(ctx *fiber.Ctx) string {
	return fmt.Sprintf("%v", ctx.Locals(constants.LocalsUserId))
}
