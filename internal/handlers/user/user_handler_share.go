// Brief   : Share users between clubs.
// Author  : Nicolas Baeriswyl <nicolasfabrice.baeriswyl@edu.hefr.ch>
// Created : 11.07.2023

package userHandler

import (
	"errors"
	"fmt"
	"github.com/athletistics/core-api/database"
	"github.com/athletistics/core-api/helpers"
	"github.com/athletistics/core-api/internal/models"
	res "github.com/athletistics/core-api/internal/responses"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

// createShareCode creates a new entry in the UserShare table and returns the created object.
// The object's ID can then be used to add the user on another club.
func createShareCode(user *models.User, targetUserId string) (*models.UserShare, *res.Error) {
	db := database.DB
	if len(targetUserId) == 0 {
		helpers.LogInfo("bad request: missing user id", user.ID)
		return nil, res.ErrBadRequest
	}

	targetUser := models.User{ID: targetUserId}
	if err := db.First(&targetUser).Error; err != nil {
		return nil, res.NewErrorFromDb(err)
	}

	// Make sure requesting user is head coach on target user's parent club
	if !helpers.HasRoleOnClub(user, targetUser.ClubID, helpers.RoleHeadCoach) {
		helpers.LogSecurity(fmt.Sprintf("tried sharing user %v", targetUserId), user.ID)
		return nil, res.ErrForbidden
	}

	// Create and return share code
	shareCode := models.UserShare{UserID: targetUser.ID}
	if err := db.Create(&shareCode).Error; err != nil {
		return nil, res.NewErrorFromDb(err)
	}

	return &shareCode, nil
}

// redeemShareCode adds a user to the requesting user's parent club based on a share code. If the user is already added
// to the target club, the request is considered successful. The added user is returned.
func redeemShareCode(user *models.User, code string) (*models.UserShort, *res.Error) {
	db := database.DB

	// verify the code is in a valid format
	codeUuid, err := uuid.Parse(code)
	if err != nil {
		helpers.LogInfo(fmt.Sprintf("bad request: %v", err), user.ID)
		return nil, res.ErrBadRequest.Err(err)
	} else if codeUuid == uuid.Nil {
		return nil, res.ErrBadRequest
	}

	// make sure user has required permissions
	if !helpers.HasRoleOnParentClub(user, helpers.RoleHeadCoach) {
		helpers.LogSecurity(fmt.Sprintf("tried redeeming code %v", code), user.ID)
		return nil, res.ErrForbidden
	}

	// find the corresponding user share
	share := models.UserShare{ID: codeUuid}
	yesterday := time.Now().AddDate(0, 0, -1)
	if err := db.Where("created_at > ?", yesterday).First(&share).Error; err != nil {
		return nil, res.NewErrorFromDb(err)
	}

	err = db.Transaction(func(tx *gorm.DB) error {
		// add the user to the requesting user's parent club as athlete
		auth := models.Authorisation{AccessRoleID: models.AccessRoleIdAthlete, UserID: share.UserID, ClubID: user.ClubID}
		if err := tx.Create(&auth).Error; err != nil {
			return err
		}

		if err := tx.Unscoped().Delete(&share).Error; err != nil {
			return err
		}

		return nil
	})
	if err != nil && !errors.Is(err, gorm.ErrDuplicatedKey) {
		return nil, res.NewErrorFromDb(err)
	}

	// return the added user
	var addedUser models.UserShort
	if err := db.Model(&models.User{}).First(&addedUser, share.UserID).Error; err != nil {
		return nil, res.NewErrorFromDb(err)
	}
	return &addedUser, nil
}
