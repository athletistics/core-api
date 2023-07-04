// Brief   : Middleware extracting JWT claims and putting them into the fiber context as locals.
//		   : Insipired by https://strapengine.com/build-a-custom-jwt-go-fiber-middleware/
// Author  : Nicolas Baeriswyl <nicolasfabrice.baeriswyl@edu.hefr.ch>
// Created : 17.06.2023

package jwtclaim

import (
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"gitlab.forge.hefr.ch/tb-athletistics/backend/constants"
	"gitlab.forge.hefr.ch/tb-athletistics/backend/helpers"
	res "gitlab.forge.hefr.ch/tb-athletistics/backend/internal/responses"
	"net/http"
	"strings"
)

const (
	AuthHeaderName        = "Authorization"
	AuthHeaderSplitLength = 2
)

// -----------------------------------
// Config is currently not in use, but
// prepared in case it is required in
// the future
// -----------------------------------

type Config struct {
	Filter func(c *fiber.Ctx) bool
}

var ConfigDefault = Config{
	Filter: nil,
}

/*
configDefault creates a new config based on the config given as parameter. The new config takes all values from the
parameter and fills in missing values.
*/
//nolint:unused
func configDefault(config []Config) Config {
	if len(config) < 1 {
		return ConfigDefault
	}

	cfg := config[0]

	if cfg.Filter == nil {
		cfg.Filter = ConfigDefault.Filter
	}

	return cfg
}

/*
New returns a middleware that reads claims in the JWT token, and stores them in Fiber's Locals.

Warning: The token is not verified in this step, and must thus have been verified in a previous middleware!!
*/
func New(config ...Config) fiber.Handler {
	// cfg := configDefault(config) // uncomment this, as soon as the config has values

	return func(c *fiber.Ctx) error {
		// get the auth header
		authHeader, ok := c.GetReqHeaders()[AuthHeaderName]
		if !ok {
			helpers.LogSecurity("missing auth header", "")
			return c.Status(http.StatusUnauthorized).JSON(res.ErrMissingAuth)
		}

		// split header into parts, as the token is prefixed with "Bearer "
		authHeaderParts := strings.Fields(authHeader)
		if len(authHeaderParts) != AuthHeaderSplitLength {
			helpers.LogSecurity("bad auth header", "")
			return c.Status(http.StatusUnauthorized).JSON(res.ErrBadAuthHeader)
		}

		rawToken := authHeaderParts[AuthHeaderSplitLength-1]

		// Extract claims without verifying the token
		claims := jwt.RegisteredClaims{}
		_, _, err := jwt.NewParser().ParseUnverified(rawToken, &claims)
		if err != nil {
			helpers.LogSecurity("bad auth token", "")
			return c.Status(http.StatusUnauthorized).JSON(res.ErrBadToken)
		}

		// make sure the subject exists, and store it in the locals
		if len(claims.Subject) == 0 {
			helpers.LogSecurity("missing subject in auth token", "")
			return c.Status(http.StatusUnauthorized).JSON(res.ErrMissingSubject)
		}
		c.Locals(constants.LocalsUserId, claims.Subject)

		return c.Next()
	}
}
