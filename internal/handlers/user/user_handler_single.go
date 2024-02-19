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
	"net/http"
	"strconv"
)

func getUser(requestingUserId string, userId string) (*models.User, int, *res.Error) {
	// read the requested user id, and verify that it is a numeric string
	if _, err := strconv.ParseInt(userId, 10, 64); err != nil {
		helpers.LogInfo(fmt.Sprintf("bad request: %v", err), requestingUserId)
		return nil, http.StatusBadRequest, res.ErrBadRequest.Err(err)
	}

	// search the user
	targetUser := models.User{ID: userId}
	db := database.DB
	err := db.Preload("Nationality").Preload("Country").Preload("Club").
		First(&targetUser).Error
	if err != nil {
		helpers.LogError(err.Error(), requestingUserId)
		return nil, http.StatusInternalServerError, res.ErrDB.Err(err)
	}

	return &targetUser, http.StatusOK, nil
}

func updateUser(user *models.User, userId string, body *models.User, parseErr error) *res.Error {
	db := database.DB

	// verify parameters
	if len(userId) == 0 {
		return res.ErrBadRequest
	}
	if parseErr != nil {
		helpers.LogInfo(fmt.Sprintf("bad request: %v", parseErr), user.ID)
		return res.ErrBadRequest.Err(parseErr)
	}

	// verify permissions. requesting user must be coach or head coach of the user that is to be updated
	userToModify := models.User{}
	if err := db.First(&userToModify, userId).Error; err != nil {
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
