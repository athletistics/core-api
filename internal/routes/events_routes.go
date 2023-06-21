// Brief   :
// Author  : Nicolas Baeriswyl <nicolasfabrice.baeriswyl@edu.hefr.ch>
// Created : 18.06.2023

package routes

import "github.com/gofiber/fiber/v2"

func SetupEventsRoutes(router fiber.Router) {
	node := router.Group("/events")

	node.Get("/")
}

func SetupEventsRoutesTotpProtected(router fiber.Router) {
	node := router.Group("/events")

	node.Post("/")
}
