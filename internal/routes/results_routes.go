// Brief   :
// Author  : Nicolas Baeriswyl <nicolasfabrice.baeriswyl@edu.hefr.ch>
// Created : 18.06.2023

package routes

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	resultHandler "gitlab.forge.hefr.ch/tb-athletistics/backend/internal/handlers/result"
)

func SetupResultsRoutes(router fiber.Router) {
	node := router.Group("/results")

	node.Get("/static", resultHandler.GetStaticResults)
}

func SetupResultsRoutesTotpProtected(router fiber.Router) {
	node := router.Group("/results")

	node.Get("/", resultHandler.GetResults)
	node.Post("/", resultHandler.CreateResult)
	node.Get("/best", resultHandler.GetBestResults)

	node.Put(fmt.Sprintf("/:%v", resultHandler.ParamResultId), resultHandler.UpdateResult)
	node.Delete(fmt.Sprintf("/:%v", resultHandler.ParamResultId), resultHandler.DeleteResult)
}
