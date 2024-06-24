package main

import (
	"fmt"
	"github.com/athletistics/core-api/config"
	"github.com/athletistics/core-api/constants"
	"github.com/athletistics/core-api/database"
	"github.com/athletistics/core-api/helpers"
	"github.com/athletistics/core-api/internal/middlewares/contextprovider"
	"github.com/athletistics/core-api/internal/middlewares/jwtclaim"
	"github.com/athletistics/core-api/internal/middlewares/totpverify"
	"github.com/athletistics/core-api/internal/middlewares/userprovider"
	"github.com/athletistics/core-api/router"
	"github.com/baeriswyln/goseeder"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/adaptor"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/requestid"
	"github.com/gofiber/swagger"
	http_mw "github.com/zitadel/zitadel-go/v2/pkg/api/middleware/http"
	"log"
	"os"

	"github.com/zitadel/zitadel-go/v2/pkg/client/middleware"

	_ "github.com/athletistics/core-api/database/seeds"
	_ "github.com/athletistics/core-api/docs"
)

//	@title			Core Athletistics API
//	@version		v1.1.3
//	@description	API that manages the core Athletistics database containing athletes and their sports performances.

//	@contact.name	API Support
//	@contact.email	nicolasfabrice.baeriswyl@edu.hefr.ch

//	@license.name	AGPL-3.0
//	@license.url	https://www.gnu.org/licenses/agpl-3.0.html

//	@host		athletistics.kube.isc.heia-fr.ch
//	@schemes	https
//	@BasePath	/api/v1

// @securityDefinitions.apikey	Zitadel JWT token
// @in							header
// @name						Authorization
func main() {
	database.ConnectDb()

	// prepare logger
	f, err := os.OpenFile("./core-api.log", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644) //nolint:mnd
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	log.SetOutput(f)

	goseeder.WithSeeder(database.GetSqlDB, func() {
		fiberSetup()
	})
}

func fiberSetup() {
	// write the API Key to a file, allowing Zitadel Interceptor to use it
	if err := CreateKeyFile(); err != nil {
		helpers.LogFatal(err)
	}

	introspection, err := http_mw.NewIntrospectionInterceptor(os.Getenv("ZITADEL_INSTANCE"), middleware.OSKeyPath())
	if err != nil {
		helpers.LogFatal(err)
	}

	// Create connection to the Zitadel Management API
	helpers.InitZitadelClients()
	defer helpers.CloseZitadelClients()

	// Prepare the ID generator
	helpers.InitIdGen()

	// clean up key file, as it is no longer needed as file
	if err := DeleteKeyFile(); err != nil {
		helpers.LogFatal(err)
	}

	app := fiber.New()

	file := useFiberLogger(app)
	defer file.Close()

	// set up the monitor that is protected through BasicAuth
	router.SetupMonitor(app)

	useCors(app)

	// attach swagger doc
	mountDocs(app)

	// check auth token on each call
	app.Use(adaptor.HTTPMiddleware(introspection.Handler))

	// extract user id from JWT token
	app.Use(jwtclaim.New())

	// load the requesting users's context
	app.Use(contextprovider.New())

	// extract user from db using id that was previously extracted
	app.Use(userprovider.New())

	// store the user's 2FA state
	app.Use(totpverify.New())

	router.SetupRoutes(app)

	err = app.Listen(":3000")
	if err != nil {
		helpers.LogFatal(err)
	}
}

// mountDocs listens for calls on the documentation route and serves the documentation
func mountDocs(app *fiber.App) {
	// attach swagger doc
	app.Get(fmt.Sprintf("/api/v%v/docs/*", config.ApiVersion), swagger.HandlerDefault)
	app.Get(fmt.Sprintf("/api/v%v/docs/*", config.ApiVersion), swagger.New(swagger.Config{ // custom
		URL:         "https://athletistics.com/doc.json",
		DeepLinking: false,
		// Expand ("list") or Collapse ("none") tag groups by default
		DocExpansion: "none",
	}))
}

// useFiberLogger configures Fiber to keep a log of all requests that are processed.
// the log includes the endpoint as well as the response code and the request duration
func useFiberLogger(app *fiber.App) *os.File {
	// prepare the logger, and add a requestid to each request
	file, err := os.OpenFile("./request.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644) //nolint:mnd
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	app.Use(requestid.New())
	app.Use(logger.New(logger.Config{Output: file}))

	return file
}

func useCors(app *fiber.App) {
	switch os.Getenv("ENV") {
	case constants.EnvLocal:
		// serve requests from anywhere in local environments
		app.Use(cors.New())
	case constants.EnvProd:
		// secure production configuration
		app.Use(cors.New(cors.Config{
			AllowOrigins:     "https://app.athletistics.ch",
			AllowHeaders:     "",
			AllowCredentials: false,
			ExposeHeaders:    "",
		}))
	default:
		// default to secure configuration under subdomain
		app.Use(cors.New(cors.Config{
			AllowOrigins:     fmt.Sprintf("https://%v.athletistics.ch", os.Getenv("ENV")),
			AllowHeaders:     "",
			AllowCredentials: false,
			ExposeHeaders:    "",
		}))
	}
}

func CreateKeyFile() error {
	if err := os.WriteFile(os.Getenv("ZITADEL_MGMT_KEY_PATH"), []byte(os.Getenv("ZITADEL_MGMT_API_KEY")), 0755); err != nil { //nolint:mnd
		return err
	}
	return os.WriteFile(middleware.OSKeyPath(), []byte(os.Getenv("ZITADEL_API_KEY")), 0755) //nolint:mnd
}

func DeleteKeyFile() error {
	if err := os.Remove(os.Getenv("ZITADEL_MGMT_KEY_PATH")); err != nil {
		return err
	}
	return os.Remove(middleware.OSKeyPath())
}
