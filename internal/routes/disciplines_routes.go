// Author  : Nicolas Baeriswyl <nicolasfabrice.baeriswyl@edu.hefr.ch>
// Created : 18.06.2023

package routes

import (
	"github.com/gofiber/fiber/v2"
	disciplineHandler "gitlab.forge.hefr.ch/tb-athletistics/backend/internal/handlers/discipline"
)

func SetupDisciplinesRoutes(router fiber.Router) {
	node := router.Group("/disciplines")

	node.Get("/", disciplineHandler.GetDisciplines)
}

func SetupDisciplinesRoutesTotpProtected(router fiber.Router) {
	node := router.Group("/disciplines")

	node.Post("/", disciplineHandler.CreateDisciplines)
}
