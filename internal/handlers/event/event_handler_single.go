// Brief   : Actions on a single event
// Author  : Nicolas Baeriswyl <nicolasfabrice.baeriswyl@edu.hefr.ch>
// Created : 16.06.2024

package eventHandler

import (
	"errors"
	"fmt"
	"github.com/athletistics/core-api/database"
	h "github.com/athletistics/core-api/helpers"
	"github.com/athletistics/core-api/internal/models"
	res "github.com/athletistics/core-api/internal/responses"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func permissionCheck(ctx *h.Context, eventId string) (models.Event, *res.Error) {
	event := models.Event{}
	if ctx.User() == nil {
		return event, res.ErrForbidden
	}

	eventUuid, err := uuid.Parse(eventId)
	if err != nil {
		h.LogInfo(fmt.Sprintf("bad request: %v", err), ctx.User().ID)
		return event, res.ErrBadRequest.Err(err)
	} else if eventUuid == uuid.Nil {
		return event, res.ErrBadRequest.Err(err)
	}
	event.ID = eventUuid

	// Permission check
	if !ctx.HasRoleOnClubOwned(&event, h.RoleCoach) {
		h.LogSecurity("tried updating user", ctx.User().ID)
		return event, res.ErrForbidden
	}

	return event, nil
}

func updateEvent(ctx *h.Context, eventId string, body *models.EventInput, parseErr error) *res.Error {
	db := database.DB

	if parseErr != nil {
		h.LogInfo(fmt.Sprintf("bad request: %v", parseErr), ctx.User().ID)
		return res.ErrBadRequest.Err(parseErr)
	}

	event, err := permissionCheck(ctx, eventId)
	if err != nil {
		return err
	}

	event.Name = body.Name
	event.Place = body.Place
	event.Date = body.Date
	event.EventTypeID = body.EventTypeID
	event.IsHomologated = body.IsHomologated
	event.IsOutdoor = body.IsOutdoor

	if err := db.Save(event).Error; err != nil {
		return res.ErrDB.Err(err)
	}

	return nil
}

func deleteEvent(ctx *h.Context, eventId string, force bool) *res.Error {
	db := database.DB
	event, errPermissions := permissionCheck(ctx, eventId)
	if errPermissions != nil {
		return errPermissions
	}

	err := db.Transaction(func(tx *gorm.DB) error {
		// Make sure the event can be deleted. If the event has results, force must be true.
		var count int64
		if err := tx.Model(&models.Result{}).Where("event_id = ?", event.ID).Count(&count).Error; err != nil {
			return err
		}

		if count > 0 && !force {
			return gorm.ErrForeignKeyViolated
		}

		if err := tx.Unscoped().Where("event_id = ?", event.ID).Delete(&models.Result{}).Error; err != nil {
			return err
		}

		// Delete
		if err := tx.Unscoped().Delete(&event).Error; err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return res.ErrNotFound
		} else if errors.Is(err, gorm.ErrForeignKeyViolated) {
			return res.ErrEventHasResults
		}
		return res.ErrDB.Err(err)
	}
	return nil
}
