// Brief   : Setting up results handlers
// Author  : Nicolas Baeriswyl <nicolasfabrice.baeriswyl@edu.hefr.ch>
// Created : 22.06.2023

package resultHandler

import (
	"fmt"
	"github.com/athletistics/core-api/database"
	"github.com/athletistics/core-api/helpers"
	"github.com/athletistics/core-api/internal/models"
	res "github.com/athletistics/core-api/internal/responses"
	"github.com/gofiber/fiber/v2"
	"net/http"
)

const QuerySportId = "sport"
const ParamResultId = "resultId"

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
	user := helpers.GetUser(c)
	if !helpers.HasRoleOnParentClub(&user, helpers.RoleCoach) {
		helpers.LogSecurity("tried accessing results", helpers.GetUserId(c))
		return c.Status(http.StatusForbidden).JSON(res.ErrForbidden)
	}

	q := new(ResultQuery)
	if err := c.QueryParser(q); err != nil {
		helpers.LogError(err.Error(), helpers.GetUserId(c))
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
	user := helpers.GetUser(c)
	if !helpers.HasRoleOnParentClub(&user, helpers.RoleCoach) {
		helpers.LogSecurity("tried creating result", helpers.GetUserId(c))
		return c.Status(http.StatusForbidden).JSON(res.ErrForbidden)
	}

	body := new(models.ResultInput)
	if err := c.BodyParser(body); err != nil {
		helpers.LogInfo(fmt.Sprintf("bad request: %v", err), helpers.GetUserId(c))
		return c.Status(http.StatusBadRequest).JSON(res.ErrBadRequest.Err(err))
	}

	return RunCreateQuery(c, body)
}

// UpdateResult modifies the values of a result, without updating event, user, and discipline.
//
//	@Summary	Modify the values of a result, without updating event, user, and discipline.
//	@Tags		Results
//	@Consume	json
//	@Produce	json
//	@Param		resultId	path		string						true	"Result UUID to modify"
//	@Param		request		body		models.ResultUpdateInput	true	"Updated result data"
//	@Success	200			{object}	models.ResultString
//	@Failure	400			{object}	res.Error
//	@Failure	401			{object}	res.Error
//	@Failure	403			{object}	res.Error
//	@Security	ApiKeyAuth
//	@Router		/results/{resultId} [put]
func UpdateResult(c *fiber.Ctx) error {
	b := new(models.ResultUpdateInput)
	parseErr := c.BodyParser(b)
	user := helpers.GetUser(c)

	resultId := c.Params(ParamResultId)

	result, err := updateResult(&user, resultId, b, parseErr)
	return helpers.ProcessJsonResponse(c, result, err)
}

// DeleteResult deletes a specific result.
//
//	@Summary	Delete a specific result.
//	@Tags		Results
//	@Param		resultId	path		string	true	"Result UUID to delete"
//	@Success	200			{string}	string
//	@Failure	400			{object}	res.Error
//	@Failure	401			{object}	res.Error
//	@Failure	403			{object}	res.Error
//	@Security	ApiKeyAuth
//	@Router		/results/{resultId} [delete]
func DeleteResult(c *fiber.Ctx) error {
	user := helpers.GetUser(c)
	resultId := c.Params(ParamResultId)

	return helpers.ProcessJsonResponse(c, nil, deleteResult(&user, resultId))
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
	user := helpers.GetUser(c)
	if !helpers.HasRoleOnParentClub(&user, helpers.RoleCoach) {
		helpers.LogSecurity("tried accessing results", helpers.GetUserId(c))
		return c.Status(http.StatusForbidden).JSON(res.ErrForbidden)
	}

	q := new(ResultQuery)
	if err := c.QueryParser(q); err != nil {
		helpers.LogInfo(fmt.Sprintf("bad request: %v", err), helpers.GetUserId(c))
		return c.Status(http.StatusBadRequest).JSON(res.ErrBadRequest.Err(err))
	}

	q.BestOnly = true

	return RunQuery(c, q)
}

// GetStaticResults lists static results (e.g. national limits) stored in the database.
//
//	@Summary	Get static results
//	@Tags		Results
//	@Produce	json
//	@Param		sport	query		integer	true	"Sport ID"
//	@Success	200		{object}	models.ResultStaticOut
//	@Failure	401		{object}	res.Error
//	@Security	ApiKeyAuth
//	@Router		/results/static [get]
func GetStaticResults(c *fiber.Ctx) error {
	db := database.DB
	var results []models.ResultStatic
	var types []models.ResultStaticType

	sportId := c.QueryInt(QuerySportId)
	if sportId == 0 {
		helpers.LogInfo("bad request: missing parameter 'sportId'", helpers.GetUserId(c))
		return c.Status(http.StatusBadRequest).JSON(res.ErrBadRequest)
	}

	if err := db.Find(&results).Error; err != nil {
		helpers.LogError(err.Error(), helpers.GetUserId(c))
		return c.Status(http.StatusInternalServerError).JSON(res.ErrDB.Err(err))
	}
	if err := db.Find(&types).Error; err != nil {
		helpers.LogError(err.Error(), helpers.GetUserId(c))
		return c.Status(http.StatusInternalServerError).JSON(res.ErrDB.Err(err))
	}

	return c.JSON(models.ResultStaticOut{
		Results: results,
		Types:   types,
	})
}

// CreatePersonalTrainingResult creates a new training result for the requesting user.
//
//	@Summary	Create a new result.
//	@Tags		Results
//	@Consume	json
//	@Produce	json
//	@Param		request	body		models.ResultTrainingInput	true	"Result data"
//	@Success	200		{object}	res.UuidCreated
//	@Failure	400		{object}	res.Error
//	@Failure	401		{object}	res.Error
//	@Security	ApiKeyAuth
//	@Router		/results/training/ [post]
func CreatePersonalTrainingResult(c *fiber.Ctx) error {
	body := new(models.ResultTrainingInput)
	ctx := helpers.LoadContext(c)
	if err := c.BodyParser(body); err != nil {
		helpers.LogInfo(fmt.Sprintf("bad request: %v", err), ctx.User().ID)
		return c.Status(http.StatusBadRequest).JSON(res.ErrBadRequest.Err(err))
	}

	return helpers.ProcessJsonResponse(c, nil, createTrainingResult(ctx, body))
}

// UpdatePersonalTrainingResult modifies the values of a result, without updating user and discipline.
//
//	@Summary	Modify the values of a result, without updating event, user, and discipline.
//	@Tags		Results
//	@Consume	json
//	@Produce	json
//	@Param		resultId	path		string						true	"Result UUID to modify"
//	@Param		request		body		models.ResultTrainingInput	true	"Updated result data"
//	@Success	200			{object}	models.ResultString
//	@Failure	400			{object}	res.Error
//	@Failure	401			{object}	res.Error
//	@Security	ApiKeyAuth
//	@Router		/results/training/{resultId} [put]
func UpdatePersonalTrainingResult(c *fiber.Ctx) error {
	body := new(models.ResultTrainingInput)
	ctx := helpers.LoadContext(c)
	resultId := c.Params(ParamResultId)

	if err := c.BodyParser(body); err != nil {
		helpers.LogInfo(fmt.Sprintf("bad request: %v", err), ctx.User().ID)
		return c.Status(http.StatusBadRequest).JSON(res.ErrBadRequest.Err(err))
	}

	return helpers.ProcessJsonResponse(c, nil, updateTrainingResult(ctx, resultId, body))
}

// DeletePersonalTrainingResult deletes a specific training result.
//
//	@Summary	Delete a specific result.
//	@Tags		Results
//	@Param		resultId	path		string	true	"Result UUID to delete"
//	@Success	200			{string}	string
//	@Failure	400			{object}	res.Error
//	@Failure	401			{object}	res.Error
//	@Security	ApiKeyAuth
//	@Router		/results/training/{resultId} [delete]
func DeletePersonalTrainingResult(c *fiber.Ctx) error {
	ctx := helpers.LoadContext(c)
	resultId := c.Params(ParamResultId)

	return helpers.ProcessJsonResponse(c, nil, deleteTrainingResult(ctx, resultId))
}
