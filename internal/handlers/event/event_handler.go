// Brief   : Event management handlers
// Author  : Nicolas Baeriswyl <nicolasfabrice.baeriswyl@edu.hefr.ch>
// Created : 21.06.2023

package eventHandler

import (
	"errors"
	"fmt"
	"github.com/athletistics/core-api/database"
	h "github.com/athletistics/core-api/helpers"
	"github.com/athletistics/core-api/internal/models"
	res "github.com/athletistics/core-api/internal/responses"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
	"net/http"
)

const (
	QueryForce   = "force"
	ParamEventId = "id"
)

// GetEvents lists all global events as well as the club's custom events.
//
//	@Summary	List events
//	@Tags		Event
//	@Produce	json
//	@Success	200	{object}	[]models.EventShort
//	@Failure	401	{object}	res.Error
//	@Failure	403	{object}	res.Error	"Not enough permissions"
//	@Security	ApiKeyAuth
//	@Router		/events [get]
func GetEvents(c *fiber.Ctx) error {
	user := h.UserContext(c).User()
	if !h.UserContext(c).HasRoleOnParentClub(h.RoleAny) {
		h.LogSecurity("tried accessing events", user.ID)
		return c.Status(http.StatusForbidden).JSON(res.ErrForbidden)
	}

	db := database.DB
	var events []models.EventShort
	err := db.Model(&models.Event{}).
		Where("club_id IS NULL").
		Or(&models.Event{ClubID: &user.ClubID}).Find(&events).Error

	if err != nil {
		h.LogError(err.Error(), h.GetUserId(c))
		return c.Status(http.StatusInternalServerError).JSON(res.ErrDB.Err(err))
	}

	return c.JSON(events)
}

// CreateEvent creates a new custom event.
//
//	@Summary	Create a new event
//	@Tags		Event
//	@Consume	json
//	@Produce	json
//	@Param		request	body		models.EventInput	true	"Event data"
//	@Success	200		{object}	res.UuidCreated
//	@Failure	400		{object}	res.Error
//	@Failure	401		{object}	res.Error
//	@Failure	403		{object}	res.Error	"Not enough permissions"
//	@Security	ApiKeyAuth
//	@Router		/events [post]
func CreateEvent(c *fiber.Ctx) error {
	user := h.UserContext(c).User()
	if !h.UserContext(c).HasRoleOnParentClub(h.RoleCoach) {
		return c.Status(http.StatusForbidden).JSON(res.ErrForbidden)
	}

	// read input body
	reqBody := new(models.EventInput)
	if err := c.BodyParser(reqBody); err != nil {
		h.LogInfo(fmt.Sprintf("bad request: %v", err), h.GetUserId(c))
		return c.Status(http.StatusBadRequest).JSON(res.ErrBadRequest.Err(err))
	}

	// get sport id of parent club
	db := database.DB
	club := models.Club{ID: user.ClubID}
	if err := db.First(&club).Error; errors.Is(err, gorm.ErrRecordNotFound) {
		return c.Status(http.StatusNotFound).JSON(res.ErrNotFound)
	} else if err != nil {
		h.LogError(err.Error(), h.GetUserId(c))
		return c.Status(http.StatusInternalServerError).JSON(res.ErrDB.Err(err))
	}

	// write to database
	event := models.Event{
		EventTypeID:   reqBody.EventTypeID,
		Name:          reqBody.Name,
		Place:         reqBody.Place,
		ClubID:        &club.ID,
		Date:          reqBody.Date,
		IsHomologated: reqBody.IsHomologated,
	}
	if err := db.Create(&event).Error; err != nil {
		h.LogError(err.Error(), h.GetUserId(c))
		return c.Status(http.StatusInternalServerError).JSON(res.ErrDB.Err(err))
	}

	// return id of new event
	return c.JSON(res.UuidCreated{Id: event.ID})
}

// UpdateEvent updates a custom event.
//
//	@Summary	Update an event
//	@Tags		Event
//	@Consume	json
//	@Param		eventId	path		string				true	"User ID"
//	@Param		request	body		models.EventInput	true	"Event data"
//	@Success	200		{object}	res.UuidCreated
//	@Failure	400		{object}	res.Error
//	@Failure	401		{object}	res.Error
//	@Failure	403		{object}	res.Error	"Not enough permissions"
//	@Failure	404		{object}	res.Error	"Event does not exist"
//	@Security	ApiKeyAuth
//	@Router		/events/{eventId} [post]
func UpdateEvent(c *fiber.Ctx) error {
	body := new(models.EventInput)
	parseErr := c.BodyParser(body)
	eventId := c.Params(ParamEventId)
	ctx := h.UserContext(c)

	return h.ProcessJsonResponse(c, nil, updateEvent(ctx, eventId, body, parseErr))
}

// DeleteEvent deletes a custom event.
//
//	@Summary	Deletes an event. If results are registered for this event, the deletion only
//					succeeds if the "force" query parameter is set to true.
//	@Tags		Event
//	@Consume	json
//	@Param		force	query		boolean	false	"Force child results to be deleted"
//	@Param		eventId	path		string	true	"User ID"
//	@Success	200		{string}	OK
//	@Failure	400		{object}	res.Error	"Bad request or Event has child results"
//	@Failure	401		{object}	res.Error
//	@Failure	403		{object}	res.Error	"Not enough permissions"
//	@Failure	404		{object}	res.Error	"Event does not exist"
//	@Security	ApiKeyAuth
//	@Router		/events/{eventId} [delete]
func DeleteEvent(c *fiber.Ctx) error {
	eventId := c.Params(ParamEventId)
	force := c.QueryBool(QueryForce)
	ctx := h.UserContext(c)

	return h.ProcessJsonResponse(c, nil, deleteEvent(ctx, eventId, force))
}
