// Brief   : Choosing response status and sending fiber response based on the error state and the data.
// Author  : Nicolas Baeriswyl <nicolasfabrice.baeriswyl@edu.hefr.ch>
// Created : 06.07.2023

package helpers

import (
	"github.com/gofiber/fiber/v2"
	res "gitlab.forge.hefr.ch/tb-athletistics/backend/internal/responses"
	"net/http"
)

// ProcessJsonResponse decides how to send data back to the requester. If err is nil,
// the data will be sent in JSON format with status code 200. Otherwise, a log is
// created (log level depends on error's http status code), and the error is sent
// back to the requester, with the status code changed to the error's code.
func ProcessJsonResponse(c *fiber.Ctx, data interface{}, err *res.Error) error {
	if err == nil {
		if data == nil {
			return c.SendStatus(http.StatusOK)
		} else {
			return c.JSON(data)
		}
	}

	return processError(c, err)
}

// ProcessCsvResponse decides how to send data back to the requester. If err is nil,
// the data will be sent in CSV format with status code 200. Otherwise, a log is
// created (log level depends on error's http status code), and the error is sent
// back to the requester, with the status code changed to the error's code.
func ProcessCsvResponse(c *fiber.Ctx, data string, err *res.Error) error {
	if err == nil {
		return c.SendString(data)
	}

	return processError(c, err)
}

// processError logs an event based on the error status code, and answer the request
// with an error message
func processError(c *fiber.Ctx, err *res.Error) error {
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
