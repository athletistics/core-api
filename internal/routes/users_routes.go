// Brief   :
// Author  : Nicolas Baeriswyl <nicolasfabrice.baeriswyl@edu.hefr.ch>
// Created : 18.06.2023

package routes

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	userHandler "gitlab.forge.hefr.ch/tb-athletistics/backend/internal/handlers/user"
)

func SetupUsersRoutes(router fiber.Router) {
	node := router.Group("/users")

	node.Get("/", userHandler.GetUsers)
	node.Post("/", userHandler.CreateUser)
	node.Get(fmt.Sprintf("/:%v", userHandler.ParamUserId), userHandler.GetUser)
}
