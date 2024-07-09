// Brief   : Perform input check, permission checks, and run queries for disciplines
// Author  : Nicolas Baeriswyl <nicolas@baeriswyl.dev>
// Created : 28.06.2024

package disciplineHandler

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

// getDisciplines returns a list of disciplines. This list contains disciplines that are (1) public and (2) belong to
// the athlete's parent club.
//
// Required role: AnyRole
func getDisciplines(ctx *h.Context) ([]models.DisciplineShort, *res.Error) {
	if !ctx.HasRoleOnParentClub(h.RoleAny) {
		h.LogSecurity("tried accessing disciplines", ctx.User().ID)
		return nil, res.ErrForbidden
	}

	db := database.DB
	var disciplines []models.DisciplineShort
	err := db.Model(&models.Discipline{}).
		Where("club_id IS NULL and sport_id = (?)",
			db.Table("clubs").Select("sport_id").Where("id = ?", ctx.User().ClubID)).
		Or(&models.Discipline{ClubID: &ctx.User().ClubID}).Find(&disciplines).Error

	if err != nil {
		h.LogError(err.Error(), ctx.User().ID)
		return nil, res.ErrDB.Err(err)
	}

	for i := 0; i < len(disciplines); i++ {
		isOfficial := disciplines[i].ClubId == nil
		disciplines[i].IsOfficial = &isOfficial
	}

	return disciplines, nil
}

func createDiscipline(ctx *h.Context, body *models.DisciplineInput, parseErr error) (*uuid.UUID, *res.Error) {
	// permission check
	if !ctx.HasRoleOnParentClub(h.RoleCoach) {
		return nil, res.ErrForbidden
	}

	// check input body validity
	if parseErr != nil {
		h.LogInfo(fmt.Sprintf("bad request: %v", parseErr), ctx.User().ID)
		return nil, res.ErrBadRequest.Err(parseErr)
	}

	// get sport id of parent club
	db := database.DB
	club := models.Club{ID: ctx.User().ClubID}
	if err := db.First(&club).Error; errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, res.ErrNotFound
	} else if err != nil {
		h.LogError(err.Error(), ctx.User().ID)
		return nil, res.ErrDB.Err(err)
	}

	// write to database
	discipline := models.Discipline{
		SportID:      club.SportID,
		ClubID:       &club.ID,
		CustomName:   &body.CustomName,
		IsDescending: body.IsDescending,
	}
	if err := db.Create(&discipline).Error; err != nil {
		h.LogError(err.Error(), ctx.User().ID)
		return nil, res.ErrDB.Err(err)
	}

	// return id of new discipline
	return &discipline.ID, nil
}

func updateDiscipline(ctx *h.Context, disciplineId string, body *models.DisciplineInput, parseErr error) *res.Error {
	// check discipline id
	id, err := uuid.Parse(disciplineId)
	if err != nil {
		return res.ErrBadRequest.Err(err)
	} else if id == uuid.Nil {
		return res.ErrBadRequest
	}

	// check input body validity
	if parseErr != nil {
		h.LogInfo(fmt.Sprintf("bad request: %v", parseErr), ctx.User().ID)
		return res.ErrBadRequest.Err(parseErr)
	}

	db := database.DB
	discipline := &models.Discipline{ID: id}
	if err := db.First(discipline).Error; err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		return res.ErrNotFound
	} else if err != nil {
		return res.ErrDB.Err(err)
	}

	// permission check
	if discipline.ClubID == nil || !ctx.HasRoleOnClub(*discipline.ClubID, h.RoleCoach) {
		return res.ErrForbidden
	}

	discipline.CustomName = &body.CustomName
	discipline.IsDescending = body.IsDescending

	if err := db.Save(discipline).Error; err != nil {
		return res.ErrDB.Err(err)
	}

	return nil
}

func deleteDiscipline(ctx *h.Context, disciplineId string, force bool) *res.Error {
	// check discipline id
	id, err := uuid.Parse(disciplineId)
	if err != nil {
		return res.ErrBadRequest.Err(err)
	} else if id == uuid.Nil {
		return res.ErrBadRequest
	}

	db := database.DB
	discipline := &models.Discipline{ID: id}
	if err := db.First(discipline).Error; err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		return res.ErrNotFound
	} else if err != nil {
		return res.ErrDB.Err(err)
	}

	// permission check
	if discipline.ClubID == nil || !ctx.HasRoleOnClub(*discipline.ClubID, h.RoleCoach) {
		return res.ErrForbidden
	}

	// delete the discipline. if the deletion is forced, also delete all results, training results and static results
	// connected to the given discipline
	err = db.Transaction(func(tx *gorm.DB) error {
		if force {
			if err := tx.Unscoped().Where("discipline_id = ?", id).Delete(&models.Result{}).Error; err != nil {
				return err
			}
			if err := tx.Unscoped().Where("discipline_id = ?", id).Delete(&models.ResultTraining{}).Error; err != nil {
				return err
			}
			if err := tx.Unscoped().Where("discipline_id = ?", id).Delete(&models.ResultStatic{}).Error; err != nil {
				return err
			}
		} else if countDisciplineDependencies(id, tx) != 0 {
			return gorm.ErrForeignKeyViolated
		}

		if err := tx.Unscoped().Delete(discipline).Error; err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return res.ErrNotFound
		}
		if errors.Is(err, gorm.ErrForeignKeyViolated) {
			return res.ErrDisciplineHasResults
		}
		return res.ErrDB.Err(err)
	}

	return nil
}

func countDisciplineDependencies(id uuid.UUID, db *gorm.DB) int64 {
	var countResult int64
	var countResultTraining int64
	var countResultStatic int64
	db.Model(&models.Result{}).Where("discipline_id = ?", id).Count(&countResult)
	db.Model(&models.ResultTraining{}).Where("discipline_id = ?", id).Count(&countResultTraining)
	db.Model(&models.ResultStatic{}).Where("discipline_id = ?", id).Count(&countResultStatic)

	return countResult + countResultTraining + countResultStatic
}
