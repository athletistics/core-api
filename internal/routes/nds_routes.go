// Brief   : Router for all NDS routes
// Author  : Nicolas Baeriswyl <nicolasfabrice.baeriswyl@edu.hefr.ch>
// Created : 18.06.2023

package routes

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	ndsHandler "gitlab.forge.hefr.ch/tb-athletistics/backend/internal/handlers/nds"
)

func SetupNdsRoutes(router fiber.Router) {
	node := router.Group("/nds")

	node.Get("/types", ndsHandler.GetTypes)
	node.Get("/activities", ndsHandler.GetActivities)
	node.Post("/activities", ndsHandler.CreateActivity)

	activity := node.Group(fmt.Sprintf("/activities/:%v", ndsHandler.ParamActivityId))
	activity.Get("/", ndsHandler.GetActivity)
	activity.Get("/instances", ndsHandler.GetInstances)
	activity.Post("/instances", ndsHandler.CreateInstance)

	instance := activity.Group(fmt.Sprintf("/instances/:%v", ndsHandler.ParamInstanceId))
	instance.Get("/", ndsHandler.GetInstance)
	instance.Delete("/", ndsHandler.DeleteInstance)
	instance.Get("/participants", ndsHandler.GetParticipants)
	instance.Put("/participants", ndsHandler.SetParticipants)
}
