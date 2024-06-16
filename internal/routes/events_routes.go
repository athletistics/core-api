// Author  : Nicolas Baeriswyl <nicolasfabrice.baeriswyl@edu.hefr.ch>
// Created : 18.06.2023

package routes

import (
	"fmt"
	eventHandler "github.com/athletistics/core-api/internal/handlers/event"
	"github.com/gofiber/fiber/v2"
)

func SetupEventsRoutes(router fiber.Router) {
	node := router.Group("/events")

	node.Get("/", eventHandler.GetEvents)
}

func SetupEventsRoutesTotpProtected(router fiber.Router) {
	node := router.Group("/events")

	node.Post("/", eventHandler.CreateEvent)
	node.Put(fmt.Sprintf("/:%v", eventHandler.ParamEventId), eventHandler.UpdateEvent)
	node.Delete(fmt.Sprintf("/:%v", eventHandler.ParamEventId), eventHandler.DeleteEvent)
}
