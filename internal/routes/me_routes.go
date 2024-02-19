// Author  : Nicolas Baeriswyl <nicolasfabrice.baeriswyl@edu.hefr.ch>
// Created : 18.06.2023

package routes

import (
	meHandler "github.com/athletistics/core-api/internal/handlers/me"
	"github.com/gofiber/fiber/v2"
)

func SetupMeRoutes(router fiber.Router) {
	node := router.Group("/me")

	node.Get("/", meHandler.GetMe)
	node.Get("/auth", meHandler.GetAuth)
	node.Get("/results", meHandler.GetResults)
	node.Post("/results", meHandler.CreateResult)
	node.Get("/results/best", meHandler.GetBestResults)
}
