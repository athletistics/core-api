// Brief   : Manage access roles of users.
// Author  : Nicolas Baeriswyl <nicolasfabrice.baeriswyl@edu.hefr.ch>
// Created : 11.07.2023

package userHandler

import (
	"fmt"
	"gitlab.forge.hefr.ch/tb-athletistics/backend/database"
	"gitlab.forge.hefr.ch/tb-athletistics/backend/helpers"
	"gitlab.forge.hefr.ch/tb-athletistics/backend/internal/models"
	res "gitlab.forge.hefr.ch/tb-athletistics/backend/internal/responses"
)

// getRoles returns all roles of the target user he possesses on the requesting user's parent club. The requesting
// user requires to be coach or head coach on his parent club.
func getRoles(user *models.User, targetUserId string) ([]models.Authorisation, *res.Error) {
	db := database.DB
	// make sure requesting user is head coach or coach
	if !helpers.HasRoleOnParentClub(user, helpers.RoleHeadCoach, helpers.RoleCoach) {
		helpers.LogSecurity(fmt.Sprintf("tried accessing roles of user %v", targetUserId), user.ID)
		return nil, res.ErrForbidden
	}

	// return all authorisations the target user has on the requesting user's parent club
	var auths []models.Authorisation
	if err := db.Where(&models.Authorisation{UserID: targetUserId, ClubID: user.ClubID}).Preload("AccessRole").Find(&auths).Error; err != nil {
		return nil, res.NewErrorFromDb(err)
	}
	return auths, nil
}

// addRole adds a role to the target user. Roles can only be added to the target user's parent club.
func addRole(user *models.User, targetUserId string, body *models.Authorisation, parseErr error) *res.Error {
	db := database.DB
	if parseErr != nil {
		helpers.LogInfo(fmt.Sprintf("bad request: %v", parseErr), user.ID)
		return res.ErrBadRequest.Err(parseErr)
	}
	if body.UserID != targetUserId || body.ClubID != user.ClubID {
		helpers.LogInfo("bad request: ambiguous data", user.ID)
		return res.ErrBadRequest
	}

	// make sure requesting user is head coach
	if !helpers.HasRoleOnParentClub(user, helpers.RoleHeadCoach) {
		helpers.LogSecurity(fmt.Sprintf("tried managing roles of user %v", targetUserId), user.ID)
		return res.ErrForbidden
	}

	// make sure user exists and has same parent club
	targetUser := models.User{ID: ParamUserId, ClubID: user.ClubID}
	if err := db.Where(&targetUser).First(&targetUser).Error; err != nil {
		return res.NewErrorFromDb(err)
	}

	// add the new authorisation
	if err := db.Create(&body).Error; err != nil {
		return res.NewErrorFromDb(err)
	}
	return nil
}

// removeRole deletes a role. The requesting user must be a head coach, and he can only remove roles the target user has
// on the requesting user's parent club.
func removeRole(user *models.User, targetUserId string, body *models.Authorisation, parseErr error) *res.Error {
	db := database.DB
	if parseErr != nil {
		helpers.LogInfo(fmt.Sprintf("bad request: %v", parseErr), user.ID)
		return res.ErrBadRequest.Err(parseErr)
	}
	if body.UserID != targetUserId || body.ClubID != user.ClubID {
		helpers.LogInfo("bad request: ambiguous data", user.ID)
		return res.ErrBadRequest
	}

	// make sure requesting user is head coach
	if !helpers.HasRoleOnParentClub(user, helpers.RoleHeadCoach) {
		helpers.LogSecurity(fmt.Sprintf("tried managing roles of user %v", targetUserId), user.ID)
		return res.ErrForbidden
	}

	// delete the authorisation
	if err := db.Unscoped().Delete(&body).Error; err != nil {
		return res.NewErrorFromDb(err)
	}
	return nil
}
