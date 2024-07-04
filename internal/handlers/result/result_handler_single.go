// Brief   : Managing a single result
// Author  : Nicolas Baeriswyl <nicolasfabrice.baeriswyl@edu.hefr.ch>
// Created : 31.08.2023

package resultHandler

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

func updateResult(user *models.User, resultId string, body *models.ResultUpdateInput, parseErr error) (*models.ResultString, *res.Error) {
	db := database.DB

	// verify parameters
	resultUuid, uuidErr := uuid.Parse(resultId)
	if uuidErr != nil {
		h.LogInfo(fmt.Sprintf("bad request: %v", uuidErr), user.ID)
		return nil, res.ErrBadRequest.Err(uuidErr)
	} else if resultUuid == uuid.Nil {
		return nil, res.ErrBadRequest
	}

	if parseErr != nil {
		h.LogInfo(fmt.Sprintf("bad request: %v", parseErr), user.ID)
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
	if err := db.Model(&models.Result{}).Select("results.*, false is_personal_training_result").
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
		h.LogInfo(fmt.Sprintf("bad request: %v", uuidErr), user.ID)
		return res.ErrBadRequest.Err(uuidErr)
	} else if resultUuid == uuid.Nil {
		return res.ErrBadRequest
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

func createTrainingResult(ctx *h.Context, body *models.ResultTrainingInput) *res.Error {
	db := database.DB
	if !ctx.HasRoleOnParentClub(h.RoleAthlete) {
		return res.ErrForbidden
	}

	d := models.Discipline{ID: body.DisciplineID}
	if err := db.Where("club_id IS NULL OR club_id = ?", ctx.User().ClubID).First(&d).Error; err != nil {
		h.LogSecurity(fmt.Sprintf("tried using discipline %v", body.DisciplineID), ctx.User().ID)
		return res.ErrForbidden.Err(err)
	}

	result := body.NewDBEntry(ctx.User().ID)
	if err := db.Create(&result).Error; err != nil {
		return res.ErrDB.Err(err)
	}

	return nil
}

func updateTrainingResult(ctx *h.Context, resultId string, body *models.ResultTrainingInput) *res.Error {
	db := database.DB
	resultUuid, uuidErr := uuid.Parse(resultId)
	if uuidErr != nil {
		return res.ErrBadRequest.Err(uuidErr)
	} else if resultUuid == uuid.Nil {
		return res.ErrBadRequest
	}

	result := &models.ResultTraining{ID: resultUuid}
	if err := db.Preload("User").First(&result).Error; err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		return res.ErrNotFound
	} else if err != nil {
		return res.ErrDB.Err(err)
	}

	d := models.Discipline{ID: body.DisciplineID}
	if err := db.Where("club_id IS NULL OR club_id = ?", ctx.User().ClubID).First(&d).Error; err != nil {
		h.LogSecurity(fmt.Sprintf("tried using discipline %v", body.DisciplineID), ctx.User().ID)
		return res.ErrForbidden.Err(err)
	}

	if !ctx.HasRoleOnClub(result.User.ClubID, h.RoleCoach) && ctx.User().ID != result.User.ID {
		return res.ErrForbidden
	}

	result.User = models.User{} // prevent updating of the user due to preloading
	body.UpdateDBEntry(result)
	if err := db.Save(&result).Error; err != nil {
		return res.ErrDB.Err(err)
	}

	return nil
}

func deleteTrainingResult(ctx *h.Context, resultId string) *res.Error {
	db := database.DB
	resultUuid, uuidErr := uuid.Parse(resultId)
	if uuidErr != nil {
		return res.ErrBadRequest.Err(uuidErr)
	} else if resultUuid == uuid.Nil {
		return res.ErrBadRequest
	}

	result := &models.ResultTraining{ID: resultUuid}
	if err := db.Preload("User").First(&result).Error; err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		return res.ErrNotFound
	} else if err != nil {
		return res.ErrDB.Err(err)
	}

	if !ctx.HasRoleOnClub(result.User.ClubID, h.RoleCoach) && ctx.User().ID != result.User.ID {
		return res.ErrForbidden
	}

	result.User = models.User{} // prevent updating of the user due to preloading
	if err := db.Unscoped().Delete(&result).Error; err != nil {
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

	if !h.HasRoleOnOneOfTheClubs(user, &parentUser, h.RoleCoach) {
		h.LogSecurity("tried updating user", user.ID)
		return nil, res.ErrForbidden
	}

	return &resultToModify, nil
}
