// Brief   : Managing a single result
// Author  : Nicolas Baeriswyl <nicolasfabrice.baeriswyl@edu.hefr.ch>
// Created : 31.08.2023

package resultHandler

import (
	"fmt"
	"github.com/google/uuid"
	"gitlab.forge.hefr.ch/tb-athletistics/backend/database"
	"gitlab.forge.hefr.ch/tb-athletistics/backend/helpers"
	"gitlab.forge.hefr.ch/tb-athletistics/backend/internal/models"
	res "gitlab.forge.hefr.ch/tb-athletistics/backend/internal/responses"
)

func updateResult(user *models.User, resultId string, body *models.ResultUpdateInput, parseErr error) (*models.ResultString, *res.Error) {
	db := database.DB

	// verify parameters
	resultUuid, uuidErr := uuid.Parse(resultId)
	if uuidErr != nil {
		helpers.LogInfo(fmt.Sprintf("bad request: %v", uuidErr), user.ID)
		return nil, res.ErrBadRequest.Err(uuidErr)
	}
	if parseErr != nil {
		helpers.LogInfo(fmt.Sprintf("bad request: %v", parseErr), user.ID)
		return nil, res.ErrBadRequest.Err(parseErr)
	}

	// verify permissions. requesting user must be coach of the result's user
	resultToModify, err := verifyPermissions(user, resultUuid)
	if err != nil {
		return nil, err
	}

	resultToModify.Update(body)
	if err := db.Save(&resultToModify).Error; err != nil {
		return nil, res.ErrDB.Err(err)
	}
	var updatedResult models.ResultString
	if err := db.Model(&models.Result{}).
		Where(`results.id = ?`, resultToModify.ID).
		Preload("User").Preload("Discipline").Preload("Event").Preload("Category").
		First(&updatedResult).Error; err != nil {
		return nil, res.ErrDB.Err(err)
	}

	updatedResult.Prepare()

	return &updatedResult, nil
}

func deleteResult(user *models.User, resultId string) *res.Error {
	db := database.DB

	// verify parameters
	resultUuid, uuidErr := uuid.Parse(resultId)
	if uuidErr != nil {
		helpers.LogInfo(fmt.Sprintf("bad request: %v", uuidErr), user.ID)
		return res.ErrBadRequest.Err(uuidErr)
	}

	// verify permissions. requesting user must be coach of the result's user
	resultToDelete, err := verifyPermissions(user, resultUuid)
	if err != nil {
		return err
	}

	if err := db.Unscoped().Delete(&resultToDelete).Error; err != nil {
		return res.ErrDB.Err(err)
	}
	return nil
}

func verifyPermissions(user *models.User, resultId uuid.UUID) (*models.Result, *res.Error) {
	db := database.DB

	resultToModify := models.Result{ID: resultId}
	if err := db.First(&resultToModify).Error; err != nil {
		return nil, res.ErrDB.Err(err)
	}
	parentUser := models.User{ID: resultToModify.UserID}
	if err := db.First(&parentUser).Error; err != nil {
		return nil, res.ErrDB.Err(err)
	}

	if !helpers.HasRoleOnOneOfTheClubs(user, &parentUser, helpers.RoleCoach) {
		helpers.LogSecurity("tried updating user", user.ID)
		return nil, res.ErrForbidden
	}

	return &resultToModify, nil
}
