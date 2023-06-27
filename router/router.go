// Brief   :
// Author  : Nicolas Baeriswyl <nicolasfabrice.baeriswyl@edu.hefr.ch>
// Created : 18.06.2023

package router

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"gitlab.forge.hefr.ch/tb-athletistics/backend/config"
	"gitlab.forge.hefr.ch/tb-athletistics/backend/internal/middlewares/totprequire"
	"gitlab.forge.hefr.ch/tb-athletistics/backend/internal/routes"
)

func SetupRoutes(app *fiber.App) {
	api := app.Group(fmt.Sprintf("/api/v%v", config.ApiVersion))

	routes.SetupMeRoutes(api)
	routes.SetupDataRoutes(api)
	routes.SetupUsersRoutes(api)
	routes.SetupEventsRoutes(api)
	routes.SetupDisciplinesRoutes(api)

	api.Use(totprequire.New())

	routes.SetupDataRoutesTotpProtected(api)
	routes.SetupUsersRoutesTotpProtected(api)
	routes.SetupEventsRoutesTotpProtected(api)
	routes.SetupDisciplinesRoutesTotpProtected(api)

	routes.SetupClubRoutes(api)
	routes.SetupTeamsRoutes(api)
	routes.SetupResultsRoutes(api)
	routes.SetupNdsRoutes(api)
}
