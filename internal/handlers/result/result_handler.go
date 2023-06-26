// Brief   :
// Author  : Nicolas Baeriswyl <nicolasfabrice.baeriswyl@edu.hefr.ch>
// Created : 22.06.2023

package resultHandler

import (
	"github.com/gofiber/fiber/v2"
	"gitlab.forge.hefr.ch/tb-athletistics/backend/helpers"
	"gitlab.forge.hefr.ch/tb-athletistics/backend/internal/models"
	res "gitlab.forge.hefr.ch/tb-athletistics/backend/internal/responses"
	"net/http"
)

// GetResults lists the result according to the query parameters.
//
//	@Summary	List results.
//	@Tags		Results
//	@Produce	json
//	@Param		request	query		resultHandler.ResultQuery	true	"Result query"
//	@Success	200		{object}	[]models.ResultString
//	@Failure	401		{object}	res.Error
//	@Failure	403		{object}	res.Error	"Not enough permissions"
//	@Security	ApiKeyAuth
//	@Router		/results [get]
func GetResults(c *fiber.Ctx) error {
	if !helpers.HasRoleOnParentClub(c, helpers.RoleCoach) {
		return c.Status(http.StatusForbidden).JSON(res.ErrForbidden)
	}

	q := new(ResultQuery)
	if err := c.QueryParser(q); err != nil {
		return c.Status(http.StatusBadRequest).JSON(res.ErrBadRequest.Err(err))
	}

	return RunQuery(c, q)
}

// CreateResult creates a new training result.
//
//	@Summary	Create a new result.
//	@Tags		Results
//	@Consume	json
//	@Produce	json
//	@Param		request	body		models.ResultInput	true	"Result data"
//	@Success	200		{object}	res.UuidCreated
//	@Failure	400		{object}	res.Error
//	@Failure	401		{object}	res.Error
//	@Failure	403		{object}	res.Error
//	@Security	ApiKeyAuth
//	@Router		/results [post]
func CreateResult(c *fiber.Ctx) error {
	if !helpers.HasRoleOnParentClub(c, helpers.RoleCoach) {
		return c.Status(http.StatusForbidden).JSON(res.ErrForbidden)
	}

	body := new(models.ResultInput)
	if err := c.BodyParser(body); err != nil {
		return c.Status(http.StatusBadRequest).JSON(res.ErrBadRequest.Err(err))
	}

	return RunCreateQuery(c, body)
}

// GetBestResults lists the best result by discipline.
//
//	@Summary	Read personal bests of authenticated user.
//	@Tags		Results
//	@Produce	json
//	@Param		request	query		resultHandler.ResultQuery	true	"Result query"
//	@Success	200		{object}	[]models.ResultString
//	@Failure	401		{object}	res.Error
//	@Failure	403		{object}	res.Error	"Not enough permissions"
//	@Security	ApiKeyAuth
//	@Router		/results/best [get]
func GetBestResults(c *fiber.Ctx) error {
	if !helpers.HasRoleOnParentClub(c, helpers.RoleCoach) {
		return c.Status(http.StatusForbidden).JSON(res.ErrForbidden)
	}

	q := new(ResultQuery)
	if err := c.QueryParser(q); err != nil {
		return c.Status(http.StatusBadRequest).JSON(res.ErrBadRequest.Err(err))
	}

	q.BestOnly = true

	return RunQuery(c, q)
}

func GetStaticResults(c *fiber.Ctx) error {
	return nil
}
