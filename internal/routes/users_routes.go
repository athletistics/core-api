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

	node.Get("/emergency", userHandler.GetEmergencyUsers)
}

func SetupUsersRoutesTotpProtected(router fiber.Router) {
	node := router.Group("/users")

	node.Get("/", userHandler.GetUsers)
	node.Post("/", userHandler.CreateUser)
	node.Get("/share", userHandler.GetShareCode)
	node.Post("/share", userHandler.RedeemShareCode)

	node.Post(fmt.Sprintf("/:%v/create-login", userHandler.ParamUserId), userHandler.CreateLogin)

	node.Get(fmt.Sprintf("/:%v", userHandler.ParamUserId), userHandler.GetUser)
	node.Put(fmt.Sprintf("/:%v", userHandler.ParamUserId), userHandler.UpdateUser)

	node.Get(fmt.Sprintf("/:%v/roles", userHandler.ParamUserId), userHandler.GetRoles)
	node.Post(fmt.Sprintf("/:%v/roles", userHandler.ParamUserId), userHandler.AddRole)
	node.Delete(fmt.Sprintf("/:%v/roles", userHandler.ParamUserId), userHandler.RemoveRole)

}
