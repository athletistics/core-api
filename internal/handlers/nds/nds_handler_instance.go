// Brief   : Managing activity instances in the database
// Author  : Nicolas Baeriswyl <nicolasfabrice.baeriswyl@edu.hefr.ch>
// Created : 06.07.2023

package ndsHandler

import (
	"fmt"
	"github.com/athletistics/core-api/database"
	"github.com/athletistics/core-api/database/types"
	"github.com/athletistics/core-api/helpers"
	"github.com/athletistics/core-api/internal/models"
	res "github.com/athletistics/core-api/internal/responses"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

const daysInWeek = 7

// selectInstances returns a list of instances belonging to the requested activity. An error is returned
// in case the requesting user has not enough permissions to view the activity.
func selectInstances(user *models.User, activityId string) ([]models.NdsActivityInstance, *res.Error) {
	db := database.DB

	// verify parameters
	activityUuid, err := uuid.Parse(activityId)
	if err != nil {
		helpers.LogInfo(fmt.Sprintf("bad request: %v", err), user.ID)
		return nil, res.ErrBadRequest.Err(err)
	} else if activityUuid == uuid.Nil {
		return nil, res.ErrBadRequest.Err(err)
	}

	// make sure the user is allowed to access instances of the activity
	if !helpers.HasRoleOnActivity(user, activityUuid, helpers.RoleCoach) {
		helpers.LogSecurity(fmt.Sprintf("tried reading activity instances from %v", activityId), user.ID)
		return nil, res.ErrForbidden
	}

	var instances []models.NdsActivityInstance
	err = db.Where(&models.NdsActivityInstance{NdsActivityID: activityUuid}).Preload("NdsActivityType").Find(&instances).Error

	return instances, res.NewErrorFromDb(err)
}

// selectInstance returns a single instance belonging to a specific activity. An error is returned
// in case the requesting user has not enough permissions to view the activity.
func selectInstance(user *models.User, activityId string, instanceId string) (*models.NdsActivityInstance, *res.Error) {
	db := database.DB

	// verify parameters
	uuids, err := parseUuids(activityId, instanceId)
	if err != nil {
		helpers.LogInfo(fmt.Sprintf("bad request: %v", err), user.ID)
		return nil, res.ErrBadRequest.Err(err)
	}

	// make sure the user is allowed to access the instance of the activity
	if !helpers.HasRoleOnActivity(user, uuids[0], helpers.RoleCoach) {
		helpers.LogSecurity(fmt.Sprintf("tried reading activity instance %v", instanceId), user.ID)
		return nil, res.ErrForbidden
	}

	var instance models.NdsActivityInstance
	err = db.Where(&models.NdsActivityInstance{ID: uuids[1], NdsActivityID: uuids[0]}).Preload("NdsActivityType").First(&instance).Error

	return &instance, res.NewErrorFromDb(err)
}

// createInstance creates a new instance for a specific activity. An error is returned
// in case the requesting user has not enough permissions to manage the activity.
func createInstance(user *models.User, activityId string, body *models.NdsActivityInstanceInput, parseErr error) (*res.UuidCreated, *res.Error) {
	db := database.DB

	// verify parameters
	if parseErr != nil {
		helpers.LogInfo(fmt.Sprintf("bad request: %v", parseErr), user.ID)
		return nil, res.ErrBadRequest.Err(parseErr)
	}
	uuids, err := parseUuids(activityId)
	if err != nil {
		helpers.LogInfo(fmt.Sprintf("bad request: %v", err), user.ID)
		return nil, res.ErrBadRequest.Err(err)
	}

	// make sure the user is allowed to access the instance of the activity
	if !helpers.HasRoleOnActivity(user, uuids[0], helpers.RoleCoach) {
		helpers.LogSecurity("tried creating an instance", user.ID)
		return nil, res.ErrForbidden
	}

	// make sure the given date is between the start and end date of the parent instance
	activity := models.NdsActivity{ID: uuids[0]}
	if err = db.First(&activity).Error; err != nil {
		return nil, res.NewErrorFromDb(err)
	}
	if body.Date.T().Before(activity.StartTime.T()) || body.Date.T().After(activity.EndDate.T()) {
		helpers.LogInfo("bad request: date out of range", user.ID)
		return nil, res.ErrBadRequest
	}

	instance := models.NdsActivityInstance{
		Date:              body.Date,
		StartTime:         body.StartTime,
		Duration:          body.Duration,
		Place:             body.Place,
		NdsActivityID:     uuids[0],
		NdsActivityTypeID: body.NdsActivityTypeID,
	}

	err = db.Create(&instance).Error

	return &res.UuidCreated{Id: instance.ID}, res.NewErrorFromDb(err)
}

// deleteInstances deletes an instance and all attendance entries associated with the specified
// instance. An error is returned in case the requesting user has not enough permissions to manage the activity.
func deleteInstance(user *models.User, activityId string, instanceId string) *res.Error {
	db := database.DB

	// verify parameters
	_, instanceUuid, errParam := instanceExistsAndUserCanManage(user, activityId, instanceId)
	if errParam != nil {
		helpers.LogInfo(fmt.Sprintf("bad request: %v", errParam), user.ID)
		return errParam
	}

	err := db.Transaction(func(tx *gorm.DB) error {
		// delete all attendance entries for the specified instance
		if err := tx.Where(&models.NdsAttendance{NdsActivityInstanceID: *instanceUuid}).Delete(&models.NdsAttendance{}).Error; err != nil {
			return err
		}

		// delete the instance itself
		if err := tx.Delete(&models.NdsActivityInstance{ID: *instanceUuid}).Error; err != nil {
			return err
		}

		// commit the transaction by returning nil
		return nil
	})

	return res.NewErrorFromDb(err)
}

// instanceExistsAndUserCanManage verifies that an instance with given ID and activity ID exists, and
// that the requesting user has enough permissions to access the data. Once verified, the function
// returns the parsed UUIDs with the result (activityId, instanceId, nil)
func instanceExistsAndUserCanManage(user *models.User, activityId string, instanceId string) (*uuid.UUID, *uuid.UUID, *res.Error) {
	db := database.DB

	// verify parameters
	uuids, err := parseUuids(activityId, instanceId)
	if err != nil {
		return nil, nil, res.ErrBadRequest.Err(err)
	}

	// make sure the user is allowed to access the instance of the activity
	if !helpers.HasRoleOnActivity(user, uuids[0], helpers.RoleCoach) {
		return nil, nil, res.ErrForbidden
	}

	// check if and instance with given id and activity id exists
	var instance models.NdsActivityInstance
	if err = db.Where(&models.NdsActivityInstance{ID: uuids[1], NdsActivityID: uuids[0]}).First(&instance).Error; err != nil {
		return nil, nil, res.NewErrorFromDb(err)
	}

	return &uuids[0], &uuids[1], nil
}

// parseUuids reads a series of strings and tries to convert all of them to UUIDs.
// If one of them fails, the function will fail with an error. Otherwise, the
// outputted array contains the converted UUIDs in the order they were passed to
// the function.
func parseUuids(rawIds ...string) ([]uuid.UUID, error) {
	ids := make([]uuid.UUID, len(rawIds))

	for i := 0; i < len(rawIds); i++ {
		parsedId, err := uuid.Parse(rawIds[i])
		if err != nil {
			return nil, err
		} else if parsedId == uuid.Nil {
			return nil, fmt.Errorf("")
		}

		ids[i] = parsedId
	}

	return ids, nil
}

// instantiateSeries creates activity instances once a week between (and including)
// StartDate and EndDate on the specified weekday. This should be executed as part
// of a transaction.
func instantiateSeries(tx *gorm.DB, activity *models.NdsActivity) error {
	// find the first date to instantiate
	nextDate := time.Time(activity.StartDate)
	lastDate := time.Time(activity.EndDate)
	for nextDate.Weekday() != activity.Weekday {
		nextDate = nextDate.AddDate(0, 0, 1)
	}

	// instantiate the activity once a week between the start and end date
	for !nextDate.After(lastDate) {
		nextInstance := models.NdsActivityInstance{
			Date:              types.Date(nextDate),
			StartTime:         &activity.StartTime,
			Duration:          &activity.Duration,
			Place:             &activity.Place,
			NdsActivityID:     activity.ID,
			NdsActivityTypeID: DefaultActivityTypeId,
		}

		if err := tx.Create(&nextInstance).Error; err != nil {
			return err
		}

		// skip to next week
		nextDate = nextDate.AddDate(0, 0, daysInWeek)
	}

	// commit by returning nil
	return nil
}
