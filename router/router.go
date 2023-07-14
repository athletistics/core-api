// Brief   : Setting up all routes
// Author  : Nicolas Baeriswyl <nicolasfabrice.baeriswyl@edu.hefr.ch>
// Created : 18.06.2023

package router

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/basicauth"
	"github.com/gofiber/fiber/v2/middleware/monitor"
	"gitlab.forge.hefr.ch/tb-athletistics/backend/config"
	"gitlab.forge.hefr.ch/tb-athletistics/backend/internal/middlewares/totprequire"
	"gitlab.forge.hefr.ch/tb-athletistics/backend/internal/routes"
	"os"
)

// SetupRoutes connects all routers with fiber.
func SetupRoutes(app *fiber.App) {
	api := app.Group(fmt.Sprintf("/api/v%v", config.ApiVersion))

	routes.SetupMeRoutes(api)
	routes.SetupDataRoutes(api)
	routes.SetupUsersRoutes(api)
	routes.SetupEventsRoutes(api)
	routes.SetupResultsRoutes(api)
	routes.SetupDisciplinesRoutes(api)

	api.Use(totprequire.New())

	routes.SetupDataRoutesTotpProtected(api)
	routes.SetupUsersRoutesTotpProtected(api)
	routes.SetupEventsRoutesTotpProtected(api)
	routes.SetupResultsRoutesTotpProtected(api)
	routes.SetupExternalIdRoutesTotpProtected(api)
	routes.SetupDisciplinesRoutesTotpProtected(api)

	routes.SetupClubRoutes(api)
	routes.SetupTeamsRoutes(api)
	routes.SetupNdsRoutes(api)
}

// SetupMonitor connects the monitor middleware and adds a basic authentication.
func SetupMonitor(app *fiber.App) {
	api := app.Group(fmt.Sprintf("/api/v%v/metrics", config.ApiVersion))

	api.Use(basicauth.New(basicauth.Config{
		Users: map[string]string{
			os.Getenv("METRICS_USR"): os.Getenv("METRICS_PWD"),
		},
	}))
	api.Get("/", monitor.New())
}
