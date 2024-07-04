// Brief   : Managing activities in the database
// Author  : Nicolas Baeriswyl <nicolasfabrice.baeriswyl@edu.hefr.ch>
// Created : 06.07.2023

package ndsHandler

import (
	"fmt"
	"github.com/athletistics/core-api/database"
	"github.com/athletistics/core-api/helpers"
	"github.com/athletistics/core-api/internal/models"
	res "github.com/athletistics/core-api/internal/responses"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// selectTypes reads the complete list of nds activity types.
func selectTypes() ([]models.NdsActivityType, *res.Error) {
	db := database.DB

	var activityTypes []models.NdsActivityType
	err := db.Find(&activityTypes).Error

	return activityTypes, res.NewErrorFromDb(err)

}

// selectActivities reads all activities belonging to a specific team. This function
// fails in case the requesting user has not enough permissions on the club of the team.
func selectActivities(user *models.User, teamId string) ([]models.NdsActivity, *res.Error) {
	db := database.DB

	// verify parameters
	teamUuid, err := uuid.Parse(teamId)
	if err != nil {
		helpers.LogInfo(fmt.Sprintf("bad request: %v", err), user.ID)
		return nil, res.ErrBadRequest.Err(err)
	} else if teamUuid == uuid.Nil {
		return nil, res.ErrBadRequest.Err(err)
	}

	// find team and make sure the user can access said team
	if !helpers.HasRoleOnTeam(user, teamUuid, helpers.RoleHeadCoach, helpers.RoleCoach) {
		helpers.LogSecurity(fmt.Sprintf("tried reading activities of team  %v", teamId), user.ID)
		return nil, res.ErrForbidden
	}

	var activities []models.NdsActivity
	err = db.Where(&models.NdsActivity{TeamID: teamUuid}).Preload("Team").Find(&activities).Error

	return activities, res.NewErrorFromDb(err)

}

// selectActivities reads a single activity identified by its ID. This function
// fails in case the requesting user has not enough permissions on the club of the team owning the activity.
func selectActivity(user *models.User, activityId string) (*models.NdsActivity, *res.Error) {
	db := database.DB

	// verify parameters
	activityUuid, err := uuid.Parse(activityId)
	if err != nil {
		helpers.LogInfo(fmt.Sprintf("bad request: %v", err), user.ID)
		return nil, res.ErrBadRequest.Err(err)
	} else if activityUuid == uuid.Nil {
		return nil, res.ErrBadRequest.Err(err)
	}

	// find target activity
	var activity models.NdsActivity
	err = db.Where(&models.NdsActivity{ID: activityUuid}).Preload("Team").First(&activity).Error
	if err != nil {
		return nil, res.NewErrorFromDb(err)
	}

	// find corresponding team and make sure the user can access said team
	if !helpers.HasRoleOnTeam(user, activity.TeamID, helpers.RoleHeadCoach, helpers.RoleCoach) {
		helpers.LogSecurity(fmt.Sprintf("tried reading activity %v", activityUuid), user.ID)
		return nil, res.ErrForbidden
	}

	return &activity, nil
}

// createActivity writes a new activity to the database. This request fails in case
// the user has not enough permissions on the club of the specified team.
func createActivity(user *models.User, body *models.NdsActivityInput, parseErr error) (*res.UuidCreated, *res.Error) {
	db := database.DB

	if parseErr != nil {
		helpers.LogInfo(fmt.Sprintf("bad request: %v", parseErr), user.ID)
		return nil, res.ErrBadRequest.Err(parseErr)
	}
	if body.StartDate.T().After(body.EndDate.T()) {
		helpers.LogInfo("bad request: start date is after end date", user.ID)
		return nil, res.ErrBadRequest
	}

	// Search for target team. If target team does not exist, or the user has not enough permissions on said team,
	// mark response as forbidden
	if !helpers.HasRoleOnTeam(user, body.TeamID, helpers.RoleHeadCoach) {
		helpers.LogSecurity("tried creating activity", user.ID)
		return nil, res.ErrForbidden
	}

	dbEntry := models.NdsActivity{
		StartDate: body.StartDate,
		EndDate:   body.EndDate,
		StartTime: body.StartTime,
		Duration:  body.Duration,
		Weekday:   body.Weekday,
		Place:     body.Place,
		TeamID:    body.TeamID,
	}

	err := db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&dbEntry).Error; err != nil {
			return err
		}

		return instantiateSeries(tx, &dbEntry)
	})

	return &res.UuidCreated{Id: dbEntry.ID}, res.NewErrorFromDb(err)
}
