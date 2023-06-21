package main

import (
	"fmt"
	"github.com/baeriswyln/goseeder"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/adaptor"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/swagger"
	http_mw "github.com/zitadel/zitadel-go/v2/pkg/api/middleware/http"
	"gitlab.forge.hefr.ch/tb-athletistics/backend/config"
	"gitlab.forge.hefr.ch/tb-athletistics/backend/database"
	"gitlab.forge.hefr.ch/tb-athletistics/backend/helpers"
	"gitlab.forge.hefr.ch/tb-athletistics/backend/internal/middlewares/jwtclaim"
	"gitlab.forge.hefr.ch/tb-athletistics/backend/internal/middlewares/totpverify"
	"gitlab.forge.hefr.ch/tb-athletistics/backend/router"
	"log"
	"os"

	"github.com/zitadel/zitadel-go/v2/pkg/client/middleware"

	_ "gitlab.forge.hefr.ch/tb-athletistics/backend/database/seeds"
	_ "gitlab.forge.hefr.ch/tb-athletistics/backend/docs"
)

//	@title			Core Athletistics API
//	@version		1.0
//	@description	API that manages TOTPs and

//	@contact.name	API Support
//	@contact.email	nicolasfabrice.baeriswyl@edu.hefr.ch

//	@license.name	GPL-3.0
//	@license.url	http://www.apache.org/licenses/LICENSE-2.0.html

//	@host		athletistics.kube.isc.heia-fr.ch
//	@schemes	https
//	@BasePath	/api/v1

// @securityDefinitions.apikey	Zitadel JWT token
// @in							header
// @name						Authorization
func main() {
	database.ConnectDb()

	goseeder.WithSeeder(database.GetSqlDB, func() {
		fiberSetup()
	})
}

func fiberSetup() {
	// write the API Key to a file, allowing Zitadel Interceptor to use it
	if err := CreateKeyFile(); err != nil {
		log.Fatal(err)
	}

	introspection, err := http_mw.NewIntrospectionInterceptor(os.Getenv("ZITADEL_INSTANCE"), middleware.OSKeyPath())
	if err != nil {
		log.Fatal(err)
	}

	// Create connection to the Zitadel Management API
	helpers.InitMgmtClient()
	defer helpers.CloseMgmtClient()

	// clean up key file, as it is no longer needed as file
	defer DeleteKeyFile()
	//if err := DeleteKeyFile(); err != nil {
	//	log.Fatal(err)
	//}

	app := fiber.New()
	app.Use(cors.New())

	// attach swagger doc
	app.Get(fmt.Sprintf("/api/v%v/docs/*", config.ApiVersion), swagger.HandlerDefault)
	app.Get(fmt.Sprintf("/api/v%v/docs/*", config.ApiVersion), swagger.New(swagger.Config{ // custom
		URL:         "https://athletistics.com/doc.json",
		DeepLinking: false,
		// Expand ("list") or Collapse ("none") tag groups by default
		DocExpansion: "none",
	}))

	// check auth token on each call
	app.Use(adaptor.HTTPMiddleware(introspection.Handler))

	// extract user id from JWT token
	app.Use(jwtclaim.New())

	// store the user's 2FA state
	app.Use(totpverify.New())

	router.SetupRoutes(app)

	err = app.Listen(":3000")
	if err != nil {
		log.Fatal(err)
	}
}

func CreateKeyFile() error {
	if err := os.WriteFile(os.Getenv("ZITADEL_MGMT_KEY_PATH"), []byte(os.Getenv("ZITADEL_MGMT_API_KEY")), 0755); err != nil {
		return err
	}
	return os.WriteFile(middleware.OSKeyPath(), []byte(os.Getenv("ZITADEL_API_KEY")), 0755) //nolint:gomnd
}

func DeleteKeyFile() error {
	if err := os.Remove(os.Getenv("ZITADEL_MGMT_KEY_PATH")); err != nil {
		return err
	}
	return os.Remove(middleware.OSKeyPath())
}
