// Brief   : Functions that manage a single user
// Author  : Nicolas Baeriswyl <nicolasfabrice.baeriswyl@edu.hefr.ch>
// Created : 31.08.2023

package userHandler

import (
	"fmt"
	"gitlab.forge.hefr.ch/tb-athletistics/backend/database"
	"gitlab.forge.hefr.ch/tb-athletistics/backend/helpers"
	"gitlab.forge.hefr.ch/tb-athletistics/backend/internal/models"
	res "gitlab.forge.hefr.ch/tb-athletistics/backend/internal/responses"
)

func updateUser(user *models.User, body *models.User, parseErr error) *res.Error {
	db := database.DB

	// verify parameters
	if parseErr != nil {
		helpers.LogInfo(fmt.Sprintf("bad request: %v", parseErr), user.ID)
		return res.ErrBadRequest.Err(parseErr)
	}

	// verify permissions. requesting user must be coach or head coach of the user that is to be updated
	userToModify := models.User{}
	if err := db.First(&userToModify, body.ID).Error; err != nil {
		return res.ErrDB.Err(err)
	}
	if !helpers.HasRoleOnOneOfTheClubs(user, &userToModify, helpers.RoleCoach, helpers.RoleHeadCoach) {
		helpers.LogSecurity("tried updating user", user.ID)
		return res.ErrForbidden
	}

	// update the values of the user that are allowed to be updated
	userToModify.Update(body)
	if err := db.Save(&userToModify).Error; err != nil {
		return res.ErrDB.Err(err)
	}

	return nil
}
