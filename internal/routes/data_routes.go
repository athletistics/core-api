// Author  : Nicolas Baeriswyl <nicolasfabrice.baeriswyl@edu.hefr.ch>
// Created : 18.06.2023

package routes

import (
	dataHandler "github.com/athletistics/core-api/internal/handlers/data"
	"github.com/gofiber/fiber/v2"
)

func SetupDataRoutes(router fiber.Router) {
	node := router.Group("/data")

	node.Get("/categories", dataHandler.GetCategories)
	node.Get("/countries", dataHandler.GetCountries)
	node.Get("/languages", dataHandler.GetLanguages)
	node.Get("/eventTypes", dataHandler.GetEventTypes)
	node.Get("/sports", dataHandler.GetSports)
	node.Get("/idTypes", dataHandler.GetIdTypes)
}

func SetupDataRoutesTotpProtected(router fiber.Router) {
	node := router.Group("/data")

	node.Get("/accessRoles", dataHandler.GetAccessRoles)
}
