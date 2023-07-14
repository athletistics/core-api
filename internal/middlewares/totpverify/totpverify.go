// Brief   : Middleware reading the state of the second factor authentication
// Author  : Nicolas Baeriswyl <nicolasfabrice.baeriswyl@edu.hefr.ch>
// Created : 18.06.2023

package totpverify

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"gitlab.forge.hefr.ch/tb-athletistics/backend/config"
	"gitlab.forge.hefr.ch/tb-athletistics/backend/constants"
	"net/http"
	"os"
)

const (
	AuthHeaderName = "Authorization"
)

func New() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// get the auth header
		authHeader, ok := c.GetReqHeaders()[AuthHeaderName]
		if !ok {
			c.Locals(constants.LocalsHasVerified2FA, false)
			return c.Next()
		}

		client := &http.Client{}
		req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%v/verify?%v=%v", os.Getenv("TOTP_ENDPOINT"), constants.KeySecondsAgo, config.TotpSecondsLimit), nil)
		if err != nil {
			c.Locals(constants.LocalsHasVerified2FA, false)
			c.Locals(constants.LocalsHas2FA, false)
			return c.Next()
		}
		req.Header.Set(AuthHeaderName, authHeader)
		res, err := client.Do(req)

		// mark 2FA as verified if response has OK code
		c.Locals(constants.LocalsHasVerified2FA, err == nil && res.StatusCode == http.StatusOK)
		c.Locals(constants.LocalsHas2FA, err == nil && (res.StatusCode == http.StatusExpectationFailed || res.StatusCode == http.StatusOK))

		if res != nil {
			_ = res.Body.Close()
		}

		return c.Next()
	}
}
