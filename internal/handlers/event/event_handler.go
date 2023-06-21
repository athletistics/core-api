// Brief   :
// Author  : Nicolas Baeriswyl <nicolasfabrice.baeriswyl@edu.hefr.ch>
// Created : 21.06.2023

package eventHandler

import (
	"errors"
	"github.com/gofiber/fiber/v2"
	"gitlab.forge.hefr.ch/tb-athletistics/backend/database"
	"gitlab.forge.hefr.ch/tb-athletistics/backend/helpers"
	"gitlab.forge.hefr.ch/tb-athletistics/backend/internal/models"
	res "gitlab.forge.hefr.ch/tb-athletistics/backend/internal/responses"
	"gorm.io/gorm"
	"net/http"
)

func GetEvents(c *fiber.Ctx) error {
	if !helpers.HasAnyRoleOnParentClub(c) {
		return c.Status(http.StatusForbidden).JSON(res.ErrForbidden)
	}

	db := database.DB
	user := helpers.GetUser(c)
	var events []models.EventShort
	err := db.Model(&models.Event{}).
		Where("club_id IS NULL").
		Or(&models.Event{ClubID: &user.ClubID}).Find(&events).Error

	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(res.ErrDB.Err(err))
	}

	return c.JSON(events)
}

func CreateEvent(c *fiber.Ctx) error {
	if !helpers.HasRoleOnClub(c, []helpers.Role{helpers.RoleCoach}, helpers.GetUser(c).ClubID) {
		return c.Status(http.StatusForbidden).JSON(res.ErrForbidden)
	}

	// read input body
	reqBody := new(models.EventInput)
	if err := c.BodyParser(reqBody); err != nil {
		return c.Status(http.StatusBadRequest).JSON(res.ErrBadRequest.Err(err))
	}

	// get sport id of parent club
	db := database.DB
	club := models.Club{ID: helpers.GetUser(c).ClubID}
	if err := db.First(&club).Error; errors.Is(err, gorm.ErrRecordNotFound) {
		return c.Status(http.StatusNotFound).JSON(res.ErrNotFound)
	} else if err != nil {
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
		return c.Status(http.StatusInternalServerError).JSON(res.ErrDB.Err(err))
	}

	// return id of new event
	return c.JSON(res.UuidCreated{Id: event.ID})
}
