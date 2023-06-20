// Brief   :
// Author  : Nicolas Baeriswyl <nicolasfabrice.baeriswyl@edu.hefr.ch>
// Created : 18.06.2023

package routes

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	clubHandler "gitlab.forge.hefr.ch/tb-athletistics/backend/internal/handlers/club"
)

func SetupClubRoutes(router fiber.Router) {
	node := router.Group(fmt.Sprintf("/club/:%v", clubHandler.ParamClubUuid))

	node.Get("/", clubHandler.GetClub)
	node.Get("/teams", clubHandler.GetTeams)
	node.Post("/teams", clubHandler.CreateTeam)
}
