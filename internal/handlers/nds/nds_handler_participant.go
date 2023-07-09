// Brief   : Managing activity instance participants in the database
// Author  : Nicolas Baeriswyl <nicolasfabrice.baeriswyl@edu.hefr.ch>
// Created : 06.07.2023

package ndsHandler

import (
	"gitlab.forge.hefr.ch/tb-athletistics/backend/database"
	"gitlab.forge.hefr.ch/tb-athletistics/backend/internal/models"
	res "gitlab.forge.hefr.ch/tb-athletistics/backend/internal/responses"
	"gorm.io/gorm"
)

// selectAttendances returns a list of all users that are stored as attendees of the specified instance.
// This function fails if the requesting user has not enough permissions to manage the instance.
func selectAttendances(user *models.User, activityId string, instanceId string) ([]models.NdsAttendance, *res.Error) {
	db := database.DB

	// verify parameters and access permission
	_, instanceUuid, errParam := instanceExistsAndUserCanManage(user, activityId, instanceId)
	if errParam != nil {
		return nil, errParam
	}

	var attendances []models.NdsAttendance
	err := db.Where(&models.NdsAttendance{NdsActivityInstanceID: *instanceUuid}).Find(&attendances).Error

	return attendances, res.NewErrorFromDb(err)
}

// setAttendances updates the list of users attending an activity instance. In a transaction, all attendees are first
// removed, and the new list is then added. All users to be added must have at least one role on the club owning the
// activity instance.
func setAttendances(user *models.User, activityId string, instanceId string, body *[]models.NdsAttendanceInput, parseError error) *res.Error {
	db := database.DB

	// verify parameters and access permission
	if parseError != nil {
		return res.ErrBadRequest.Err(parseError)
	}
	activityUuid, instanceUuid, errParam := instanceExistsAndUserCanManage(user, activityId, instanceId)
	if errParam != nil {
		return errParam
	}

	// in a transaction, remove all existing attendances, and add the newly requested ones
	err := db.Transaction(func(tx *gorm.DB) error {
		// get the owner team of the activity
		team := models.Team{}
		err := db.Joins("inner join nds_activities a ON a.team_id = teams.id AND a.id = ?", activityUuid).First(&team).Error
		if err != nil {
			return err
		}

		// get list of all users to be added
		userIds := make([]string, len(*body))
		for i := 0; i < len(userIds); i++ {
			userIds[i] = (*body)[i].UserID
		}

		// count the number of users that have authorisations on the owner club of the activity instance
		var nbRecords int64
		query := models.Authorisation{ClubID: team.ClubID}
		db.Model(&models.Authorisation{}).Distinct("user_id").Where(&query).Where("user_id IN (?)", userIds).Count(&nbRecords)

		// if the counted number is equal to the number of users to add, all users are related to the club,
		// and thus they all can be added as attendees.
		if int64(len(*body)) != nbRecords {
			return gorm.ErrInvalidTransaction
		}

		// delete old entries
		if err := tx.Where(&models.NdsAttendance{NdsActivityInstanceID: *instanceUuid}).Unscoped().Delete(&models.NdsAttendance{}).Error; err != nil {
			return err
		}

		// prepare entries to be added
		if len(*body) == 0 {
			return nil
		}
		attendances := make([]models.NdsAttendance, len(*body))
		for i := 0; i < len(attendances); i++ {
			attendances[i] = models.NdsAttendance{UserID: (*body)[i].UserID, NdsActivityInstanceID: *instanceUuid, IsCoach: (*body)[i].IsCoach}
		}

		// create entries, and verify that all went well
		if dbRes := tx.Create(&attendances); dbRes.Error != nil {
			return dbRes.Error
		} else if dbRes.RowsAffected != int64(len(attendances)) {
			return gorm.ErrInvalidTransaction
		}

		return nil
	})

	return res.NewErrorFromDb(err)
}
