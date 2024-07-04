// Brief   : Managing disciplines
// Author  : Nicolas Baeriswyl <nicolasfabrice.baeriswyl@edu.hefr.ch>
// Created : 21.06.2023

package disciplineHandler

import (
	h "github.com/athletistics/core-api/helpers"
	"github.com/athletistics/core-api/internal/models"
	res "github.com/athletistics/core-api/internal/responses"
	"github.com/gofiber/fiber/v2"
)

const (
	ParamForce        = "force"
	ParamDisciplineId = "id"
)

// GetDisciplines lists all global disciplines as well as the club's custom disciplines.
//
//	@Summary	List disciplines
//	@Tags		Discipline
//	@Produce	json
//	@Success	200	{object}	[]models.DisciplineShort
//	@Failure	401	{object}	res.Error
//	@Failure	403	{object}	res.Error	"Not enough permissions"
//	@Security	ApiKeyAuth
//	@Router		/disciplines [get]
func GetDisciplines(c *fiber.Ctx) error {
	ctx := h.UserContext(c)

	disciplines, err := getDisciplines(ctx)

	if err != nil {
		return c.Status(err.Error).JSON(err)
	}

	return c.JSON(disciplines)
}

// CreateDiscipline creates a new custom club discipline.
//
//	@Summary	Create a new discipline
//	@Tags		Discipline
//	@Consumes	json
//	@Produce	json
//	@Param		request	body		models.DisciplineInput	true	"Discipline data"
//	@Success	200		{object}	res.UuidCreated
//	@Failure	400		{object}	res.Error
//	@Failure	401		{object}	res.Error
//	@Failure	403		{object}	res.Error	"Not enough permissions"
//	@Security	ApiKeyAuth
//	@Router		/disciplines [post]
func CreateDiscipline(c *fiber.Ctx) error {
	ctx := h.UserContext(c)

	// read input body
	reqBody := new(models.DisciplineInput)
	parseErr := c.BodyParser(reqBody)

	id, err := createDiscipline(ctx, reqBody, parseErr)

	if err != nil {
		return c.Status(err.Error).JSON(err)
	}

	// return id of new discipline
	return c.JSON(res.UuidCreated{Id: *id})
}

// UpdateDiscipline updates a specific discipline.
//
//	@Summary	Update a discipline. Can only be done on disciplines owned by the requesting user's club.
//	@Tags		Discipline
//	@Consumes	json
//	@Param		request			body		models.DisciplineInput	true	"Discipline data"
//	@Param		disciplineId	path		uuid.UUID				true	"Discipline ID"
//	@Success	200				{string}	OK
//	@Failure	400				{object}	res.Error
//	@Failure	401				{object}	res.Error
//	@Failure	403				{object}	res.Error	"Not enough permissions"
//	@Security	ApiKeyAuth
//	@Router		/disciplines [post]
func UpdateDiscipline(c *fiber.Ctx) error {
	ctx := h.UserContext(c)

	// read path parameter
	disciplineId := c.Params(ParamDisciplineId)

	// read input body
	reqBody := new(models.DisciplineInput)
	parseErr := c.BodyParser(reqBody)

	err := updateDiscipline(ctx, disciplineId, reqBody, parseErr)

	if err != nil {
		return c.Status(err.Error).JSON(err)
	}

	// return id of new discipline
	return c.SendStatus(fiber.StatusOK)
}

// DeleteDiscipline deletes a club discipline.
//
//	@Summary	Delete a specific discipline. Can only be done on disciplines owned by the requesting user's club.
//	@Tags		Discipline
//	@Param		disciplineId	path		uuid.UUID	true	"Discipline ID"
//	@Success	200				{string}	OK
//	@Failure	400				{object}	res.Error
//	@Failure	401				{object}	res.Error
//	@Failure	403				{object}	res.Error	"Not enough permissions"
//	@Security	ApiKeyAuth
//	@Router		/disciplines/{disciplineId} [delete]
func DeleteDiscipline(c *fiber.Ctx) error {
	ctx := h.UserContext(c)

	// read path parameter
	disciplineId := c.Params(ParamDisciplineId)
	force := c.QueryBool(ParamForce, false)

	err := deleteDiscipline(ctx, disciplineId, force)

	if err != nil {
		return c.Status(err.Error).JSON(err)
	}

	// return id of new discipline
	return c.SendStatus(fiber.StatusOK)
}
