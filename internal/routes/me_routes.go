// Brief   :
// Author  : Nicolas Baeriswyl <nicolasfabrice.baeriswyl@edu.hefr.ch>
// Created : 18.06.2023

package routes

import (
	"github.com/gofiber/fiber/v2"
	meHandler "gitlab.forge.hefr.ch/tb-athletistics/backend/internal/handlers/me"
)

func SetupMeRoutes(router fiber.Router) {
	node := router.Group("/me")

	node.Get("/", meHandler.GetMe)
	node.Get("/auth", meHandler.GetAuth)
	node.Get("/results", meHandler.GetResults)
	node.Post("/results", meHandler.CreateResult)
	node.Get("/results/best", meHandler.GetBestResults)
}
