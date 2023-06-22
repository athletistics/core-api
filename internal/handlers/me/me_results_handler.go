// Brief   :
// Author  : Nicolas Baeriswyl <nicolasfabrice.baeriswyl@edu.hefr.ch>
// Created : 22.06.2023

package meHandler

import (
	"github.com/gofiber/fiber/v2"
	"gitlab.forge.hefr.ch/tb-athletistics/backend/constants"
	"gitlab.forge.hefr.ch/tb-athletistics/backend/database"
	"gitlab.forge.hefr.ch/tb-athletistics/backend/helpers"
	resultHandler "gitlab.forge.hefr.ch/tb-athletistics/backend/internal/handlers/result"
	"gitlab.forge.hefr.ch/tb-athletistics/backend/internal/models"
	res "gitlab.forge.hefr.ch/tb-athletistics/backend/internal/responses"
	"net/http"
)

func GetResults(c *fiber.Ctx) error {
	userId := helpers.GetUserId(c)
	q := new(resultHandler.ResultQuery)
	if err := c.QueryParser(q); err != nil {
		return c.Status(http.StatusBadRequest).JSON(res.ErrBadRequest.Err(err))
	}

	q.UserID = &userId
	return resultHandler.RunQuery(c, q)
}

func CreateResult(c *fiber.Ctx) error {
	body := new(models.ResultInput)
	if err := c.BodyParser(body); err != nil {
		return c.Status(http.StatusBadRequest).JSON(res.ErrBadRequest.Err(err))
	}

	// set user id to authenticated user's id
	body.UserId = helpers.GetUserId(c)

	// make sure the event has type 4 (training)
	event := models.Event{ID: body.EventID, EventTypeID: constants.TrainingEventTypeId}
	if err := database.DB.First(&event).Error; err != nil {
		return c.Status(http.StatusForbidden).JSON(res.ErrForbidden.Err(err))
	}

	return resultHandler.RunCreateQuery(c, body)
}

func GetBestResults(c *fiber.Ctx) error {
	userId := helpers.GetUserId(c)
	q := new(resultHandler.ResultQuery)
	if err := c.QueryParser(q); err != nil {
		return c.Status(http.StatusBadRequest).JSON(res.ErrBadRequest.Err(err))
	}

	q.UserID = &userId
	q.BestOnly = true

	return resultHandler.RunQuery(c, q)
}
