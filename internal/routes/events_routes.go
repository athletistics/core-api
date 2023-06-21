// Brief   :
// Author  : Nicolas Baeriswyl <nicolasfabrice.baeriswyl@edu.hefr.ch>
// Created : 18.06.2023

package routes

import (
	"github.com/gofiber/fiber/v2"
	eventHandler "gitlab.forge.hefr.ch/tb-athletistics/backend/internal/handlers/event"
)

func SetupEventsRoutes(router fiber.Router) {
	node := router.Group("/events")

	node.Get("/", eventHandler.GetEvents)
}

func SetupEventsRoutesTotpProtected(router fiber.Router) {
	node := router.Group("/events")

	node.Post("/", eventHandler.CreateEvent)
}
