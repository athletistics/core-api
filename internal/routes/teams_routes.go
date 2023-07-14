// Author  : Nicolas Baeriswyl <nicolasfabrice.baeriswyl@edu.hefr.ch>
// Created : 18.06.2023

package routes

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	teamHandler "gitlab.forge.hefr.ch/tb-athletistics/backend/internal/handlers/team"
)

func SetupTeamsRoutes(router fiber.Router) {
	node := router.Group(fmt.Sprintf("/teams/:%v", teamHandler.ParamTeamId))
	node.Get("/", teamHandler.GetTeam)
	node.Get("/participants", teamHandler.GetParticipants)
	node.Post("/participants", teamHandler.AddParticipants)
	node.Delete("/participants", teamHandler.RemoveParticipants)
}
