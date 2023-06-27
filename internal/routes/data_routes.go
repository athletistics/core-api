// Brief   :
// Author  : Nicolas Baeriswyl <nicolasfabrice.baeriswyl@edu.hefr.ch>
// Created : 18.06.2023

package routes

import (
	"github.com/gofiber/fiber/v2"
	dataHandler "gitlab.forge.hefr.ch/tb-athletistics/backend/internal/handlers/data"
)

func SetupDataRoutes(router fiber.Router) {
	node := router.Group("/data")

	node.Get("/categories", dataHandler.GetCategories)
	node.Get("/places", dataHandler.GetPlaces)
	node.Get("/places/states", dataHandler.GetStates)
	node.Get("/places/countries", dataHandler.GetCountries)
	node.Get("/languages", dataHandler.GetLanguages)
	node.Get("/eventTypes", dataHandler.GetEventTypes)
	node.Get("/sports", dataHandler.GetSports)
	node.Get("/idTypes", dataHandler.GetIdTypes)
}

func SetupDataRoutesTotpProtected(router fiber.Router) {
	node := router.Group("/data")

	node.Get("/accessRoles", dataHandler.GetAccessRoles)
}
