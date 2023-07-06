// Brief   :
// Author  : Nicolas Baeriswyl <nicolasfabrice.baeriswyl@edu.hefr.ch>
// Created : 06.07.2023

package helpers

import (
	"github.com/gofiber/fiber/v2"
	res "gitlab.forge.hefr.ch/tb-athletistics/backend/internal/responses"
	"net/http"
)

// ProcessResponse decides how to send data back to the requester. If err is nil,
// the data will be sent in JSON format with status code 200. Otherwise, a log is
// created (log level depends on error's http status code), and the error is sent
// back to the requester, with the status code changed to the error's code.
func ProcessResponse(c *fiber.Ctx, data interface{}, err *res.Error) error {
	if err == nil {
		if data == nil {
			return c.SendStatus(http.StatusOK)
		} else {
			return c.JSON(data)
		}
	}

	switch err.Error {
	case http.StatusBadRequest, http.StatusNotFound:
		LogInfo(err.String(), GetUserId(c))
		return c.Status(err.Error).JSON(err)
	case http.StatusUnauthorized, http.StatusForbidden:
		LogSecurity(err.String(), GetUserId(c))
		return c.Status(err.Error).JSON(err)
	default:
		LogError(err.String(), GetUserId(c))
		return c.Status(err.Error).JSON(err)
	}
}
