// Brief   : Result handlers for the authenticated user
// Author  : Nicolas Baeriswyl <nicolasfabrice.baeriswyl@edu.hefr.ch>
// Created : 22.06.2023

package meHandler

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"gitlab.forge.hefr.ch/tb-athletistics/backend/constants"
	"gitlab.forge.hefr.ch/tb-athletistics/backend/database"
	"gitlab.forge.hefr.ch/tb-athletistics/backend/helpers"
	resultHandler "gitlab.forge.hefr.ch/tb-athletistics/backend/internal/handlers/result"
	"gitlab.forge.hefr.ch/tb-athletistics/backend/internal/models"
	res "gitlab.forge.hefr.ch/tb-athletistics/backend/internal/responses"
	"net/http"
)

// GetResults lists all global disciplines as well as the club's custom disciplines.
//
//	@Summary	List all personal results
//	@Tags		Personal results
//	@Produce	json
//	@Success	200	{object}	[]models.ResultString
//	@Failure	401	{object}	res.Error
//	@Failure	403	{object}	res.Error	"Not enough permissions"
//	@Security	ApiKeyAuth
//	@Router		/me/results [get]
func GetResults(c *fiber.Ctx) error {
	userId := helpers.GetUserId(c)
	q := new(resultHandler.ResultQuery)
	if err := c.QueryParser(q); err != nil {
		helpers.LogInfo(fmt.Sprintf("bad request: %v", err), helpers.GetUserId(c))
		return c.Status(http.StatusBadRequest).JSON(res.ErrBadRequest.Err(err))
	}

	q.UserID = &userId
	return resultHandler.RunQuery(c, q)
}

// CreateResult creates a new training result.
//
//	@Summary	Create a new training result for the authenticated user.
//	@Tags		Personal results
//	@Consume	json
//	@Produce	json
//	@Param		request	body		models.ResultInput	true	"Result data"
//	@Success	200		{object}	res.UuidCreated
//	@Failure	400		{object}	res.Error
//	@Failure	401		{object}	res.Error
//	@Failure	403		{object}	res.Error	"Selected event is not a training"
//	@Security	ApiKeyAuth
//	@Router		/me/results [post]
func CreateResult(c *fiber.Ctx) error {
	body := new(models.ResultInput)
	if err := c.BodyParser(body); err != nil {
		helpers.LogInfo(fmt.Sprintf("bad request: %v", err), helpers.GetUserId(c))
		return c.Status(http.StatusBadRequest).JSON(res.ErrBadRequest.Err(err))
	}

	// set user id to authenticated user's id
	body.UserId = helpers.GetUserId(c)

	// make sure the event has type 4 (training)
	event := models.Event{ID: body.EventID}
	if err := database.DB.Where(&models.Event{EventTypeID: constants.TrainingEventTypeId}).First(&event).Error; err != nil {
		helpers.LogError(err.Error(), helpers.GetUserId(c))
		return c.Status(http.StatusForbidden).JSON(res.ErrForbidden.Err(err))
	}

	return resultHandler.RunCreateQuery(c, body)
}

// GetBestResults lists the best result by discipline.
//
//	@Summary	Read personal bests of authenticated user.
//	@Tags		Personal results
//	@Produce	json
//	@Success	200	{object}	[]models.ResultString
//	@Failure	401	{object}	res.Error
//	@Failure	403	{object}	res.Error	"Not enough permissions"
//	@Security	ApiKeyAuth
//	@Router		/me/results/best [get]
func GetBestResults(c *fiber.Ctx) error {
	userId := helpers.GetUserId(c)
	q := new(resultHandler.ResultQuery)
	if err := c.QueryParser(q); err != nil {
		helpers.LogInfo(fmt.Sprintf("bad request: %v", err), helpers.GetUserId(c))
		return c.Status(http.StatusBadRequest).JSON(res.ErrBadRequest.Err(err))
	}

	q.UserID = &userId
	q.BestOnly = true

	return resultHandler.RunQuery(c, q)
}
