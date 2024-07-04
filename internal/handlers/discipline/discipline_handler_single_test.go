// Author  : Nicolas Baeriswyl <nicolas@baeriswyl.dev>
// Created : 28.06.2024

package disciplineHandler

import (
	"fmt"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/athletistics/core-api/database"
	h "github.com/athletistics/core-api/helpers"
	"github.com/athletistics/core-api/internal/models"
	res "github.com/athletistics/core-api/internal/responses"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
	"testing"
)

func TestCreateDiscipline_Permission(t *testing.T) {
	clubId := uuid.New()
	// case: user has no role
	ctx := h.CreateContext(models.User{ID: "123"}, nil, h.TwoFAState{})

	_, err := createDiscipline(ctx, nil, nil)
	assert.NotNil(t, err)
	assert.Equal(t, res.ErrForbidden.Error, err.Error)

	// case: user has wrong role
	authMap := getValidAuthMap(clubId)
	authMap[clubId][h.RoleCoach] = false
	authMap[clubId][h.RoleAthlete] = true
	ctx = h.CreateContext(models.User{ID: "123", ClubID: clubId}, authMap, h.TwoFAState{})

	_, err = createDiscipline(ctx, nil, nil)
	assert.NotNil(t, err)
	assert.Equal(t, res.ErrForbidden.Error, err.Error)
}

func TestCreateDiscipline_Input(t *testing.T) {
	clubId := uuid.New()
	authMap := getValidAuthMap(clubId)

	// case: the body parsing produced an error
	ctx := h.CreateContext(models.User{ID: "123", ClubID: clubId}, authMap, h.TwoFAState{})

	_, err := createDiscipline(ctx, nil, fmt.Errorf(""))
	assert.NotNil(t, err)
	assert.Equal(t, res.ErrBadRequest.Error, err.Error)
}

func TestCreateDiscipline(t *testing.T) {
	clubId := uuid.New()
	authMap := getValidAuthMap(clubId)
	ctx := h.CreateContext(models.User{ID: "123", ClubID: clubId}, authMap, h.TwoFAState{})
	body := &models.DisciplineInput{
		CustomName:   "some name",
		IsDescending: false,
	}
	database.ConnectMockDB()

	database.Mock.ExpectQuery(`SELECT \* FROM "clubs" WHERE.*"clubs"."id" = \$1*`).
		WillReturnRows(sqlmock.NewRows([]string{"id", "sport_id"}).AddRow(clubId.String(), 1))
	database.Mock.ExpectBegin()
	database.Mock.ExpectQuery(`INSERT INTO "disciplines"`).WillReturnRows(sqlmock.NewRows([]string{}).AddRow())
	database.Mock.ExpectCommit()

	id, err := createDiscipline(ctx, body, nil)
	assert.Nil(t, err)
	assert.IsType(t, &uuid.UUID{}, id)
}

func TestUpdateDiscipline_Permission(t *testing.T) {
	clubId := uuid.New()
	disciplineId := uuid.New()
	database.ConnectMockDB()
	// case: user has no role
	ctx := h.CreateContext(models.User{ID: "123"}, nil, h.TwoFAState{})
	database.Mock.ExpectQuery(`SELECT \* FROM "disciplines" WHERE.*"disciplines"."id" = \$1*`).
		WithArgs(disciplineId, 1).
		WillReturnRows(sqlmock.NewRows([]string{"id", "club_id"}).AddRow(disciplineId, clubId))

	err := updateDiscipline(ctx, disciplineId.String(), nil, nil)
	assert.NotNil(t, err)
	assert.Equal(t, res.ErrForbidden.Error, err.Error)

	// case: user has wrong role
	authMap := getValidAuthMap(clubId)
	authMap[clubId][h.RoleCoach] = false
	authMap[clubId][h.RoleAthlete] = true
	ctx = h.CreateContext(models.User{ID: "123", ClubID: clubId}, authMap, h.TwoFAState{})
	database.Mock.ExpectQuery(`SELECT \* FROM "disciplines" WHERE.*"disciplines"."id" = \$1*`).
		WithArgs(disciplineId, 1).
		WillReturnRows(sqlmock.NewRows([]string{"id", "club_id"}).AddRow(disciplineId, clubId))

	err = updateDiscipline(ctx, disciplineId.String(), nil, nil)
	assert.NotNil(t, err)
	assert.Equal(t, res.ErrForbidden.Error, err.Error)

	// case: discipline has no club id
	authMap = getValidAuthMap(clubId)
	ctx = h.CreateContext(models.User{ID: "123", ClubID: clubId}, authMap, h.TwoFAState{})
	database.Mock.ExpectQuery(`SELECT \* FROM "disciplines" WHERE.*"disciplines"."id" = \$1*`).
		WithArgs(disciplineId, 1).
		WillReturnRows(sqlmock.NewRows([]string{"id", "club_id"}).AddRow(disciplineId, uuid.Nil))

	err = updateDiscipline(ctx, disciplineId.String(), nil, nil)
	assert.NotNil(t, err)
	assert.Equal(t, res.ErrForbidden.Error, err.Error)
}

func TestUpdateDiscipline_Input(t *testing.T) {
	clubId := uuid.New()
	disciplineId := uuid.New()
	authMap := getValidAuthMap(clubId)

	// case: the body parsing produced an error
	ctx := h.CreateContext(models.User{ID: "123", ClubID: clubId}, authMap, h.TwoFAState{})

	err := updateDiscipline(ctx, disciplineId.String(), nil, fmt.Errorf(""))
	assert.NotNil(t, err)
	assert.Equal(t, res.ErrBadRequest.Error, err.Error)

	// case: bad uuid
	err = updateDiscipline(ctx, "", nil, nil)
	assert.NotNil(t, err)
	assert.Equal(t, res.ErrBadRequest.Error, err.Error)
	err = updateDiscipline(ctx, "123", nil, nil)
	assert.NotNil(t, err)
	assert.Equal(t, res.ErrBadRequest.Error, err.Error)
	err = updateDiscipline(ctx, uuid.Nil.String(), nil, nil)
	assert.NotNil(t, err)
	assert.Equal(t, res.ErrBadRequest.Error, err.Error)
}

func TestUpdateDiscipline(t *testing.T) {
	clubId := uuid.New()
	disciplineId := uuid.New()
	authMap := getValidAuthMap(clubId)
	ctx := h.CreateContext(models.User{ID: "123", ClubID: clubId}, authMap, h.TwoFAState{})
	body := &models.DisciplineInput{
		CustomName:   "some name",
		IsDescending: false,
	}
	database.ConnectMockDB()

	database.Mock.ExpectQuery(`SELECT \* FROM "disciplines" WHERE.*"disciplines"."id" = \$1*`).
		WithArgs(disciplineId, 1).
		WillReturnRows(sqlmock.NewRows([]string{"id", "club_id"}).AddRow(disciplineId, clubId.String()))
	database.Mock.ExpectBegin()
	database.Mock.ExpectExec(`UPDATE "disciplines".* WHERE.*"id"`).
		WillReturnResult(sqlmock.NewResult(1, 1))
	database.Mock.ExpectCommit()

	err := updateDiscipline(ctx, disciplineId.String(), body, nil)
	assert.Nil(t, err)
}

func TestDeleteDiscipline_Permission(t *testing.T) {
	clubId := uuid.New()
	disciplineId := uuid.New()
	authMap := getValidAuthMap(clubId)
	ctx := h.CreateContext(models.User{ID: "123", ClubID: clubId}, authMap, h.TwoFAState{})

	database.ConnectMockDB()

	// case: permission on wrong club
	database.Mock.ExpectQuery(`SELECT \* FROM "disciplines" WHERE.*"disciplines"."id" = \$1`).
		WithArgs(disciplineId, 1).
		WillReturnRows(sqlmock.NewRows([]string{"id", "club_id"}).AddRow(disciplineId, uuid.New()))

	err := deleteDiscipline(ctx, disciplineId.String(), false)
	assert.NotNil(t, err)
	assert.Equal(t, res.ErrForbidden.Error, err.Error)

	// case: missing permission
	authMap[clubId][h.RoleCoach] = false
	authMap[clubId][h.RoleAthlete] = true
	ctx = h.CreateContext(models.User{ID: "123", ClubID: clubId}, authMap, h.TwoFAState{})

	database.Mock.ExpectQuery(`SELECT \* FROM "disciplines" WHERE.*"disciplines"."id" = \$1`).
		WithArgs(disciplineId, 1).
		WillReturnRows(sqlmock.NewRows([]string{"id", "club_id"}).AddRow(disciplineId, uuid.New()))

	err = deleteDiscipline(ctx, disciplineId.String(), false)
	assert.NotNil(t, err)
	assert.Equal(t, res.ErrForbidden.Error, err.Error)

	// case: discipline has no club id
	authMap = getValidAuthMap(clubId)
	ctx = h.CreateContext(models.User{ID: "123", ClubID: clubId}, authMap, h.TwoFAState{})
	database.Mock.ExpectQuery(`SELECT \* FROM "disciplines" WHERE.*"disciplines"."id" = \$1*`).
		WithArgs(disciplineId, 1).
		WillReturnRows(sqlmock.NewRows([]string{"id", "club_id"}).AddRow(disciplineId, uuid.Nil))

	err = deleteDiscipline(ctx, disciplineId.String(), false)
	assert.NotNil(t, err)
	assert.Equal(t, res.ErrForbidden.Error, err.Error)
}

func TestDeleteDiscipline_Input(t *testing.T) {
	clubId := uuid.New()
	authMap := getValidAuthMap(clubId)
	ctx := h.CreateContext(models.User{ID: "123", ClubID: clubId}, authMap, h.TwoFAState{})

	err := deleteDiscipline(ctx, "", false)
	assert.NotNil(t, err)
	assert.Equal(t, res.ErrBadRequest.Error, err.Error)
	err = deleteDiscipline(ctx, "123", false)
	assert.NotNil(t, err)
	assert.Equal(t, res.ErrBadRequest.Error, err.Error)
	err = deleteDiscipline(ctx, uuid.Nil.String(), false)
	assert.NotNil(t, err)
	assert.Equal(t, res.ErrBadRequest.Error, err.Error)
}

func TestDeleteDiscipline_Force(t *testing.T) {
	clubId := uuid.New()
	disciplineId := uuid.New()
	authMap := getValidAuthMap(clubId)
	ctx := h.CreateContext(models.User{ID: "123", ClubID: clubId}, authMap, h.TwoFAState{})

	database.ConnectMockDB()
	// case: deletion is not forced, but there are connections to results
	database.Mock.ExpectQuery(`SELECT \* FROM "disciplines" WHERE.*"disciplines"."id" = \$1`).
		WithArgs(disciplineId, 1).
		WillReturnRows(sqlmock.NewRows([]string{"id", "club_id"}).AddRow(disciplineId, clubId))
	database.Mock.ExpectBegin()
	database.Mock.ExpectExec(`DELETE FROM "disciplines" WHERE "disciplines"."id" = \$1`).
		WithArgs(disciplineId).
		WillReturnError(gorm.ErrForeignKeyViolated)
	database.Mock.ExpectRollback()

	err := deleteDiscipline(ctx, disciplineId.String(), false)
	assert.NotNil(t, err)
	assert.Equal(t, res.ErrDisciplineHasResults.ErrorSubCode, err.ErrorSubCode)

	// case: force
	database.Mock.ExpectQuery(`SELECT \* FROM "disciplines" WHERE.*"disciplines"."id" = \$1`).
		WithArgs(disciplineId, 1).
		WillReturnRows(sqlmock.NewRows([]string{"id", "club_id"}).AddRow(disciplineId, clubId))
	database.Mock.ExpectBegin()
	database.Mock.ExpectExec(`DELETE FROM "results" WHERE discipline_id = \$1`).
		WithArgs(disciplineId).
		WillReturnResult(sqlmock.NewResult(1, 1))
	database.Mock.ExpectExec(`DELETE FROM "result_trainings" WHERE discipline_id = \$1`).
		WithArgs(disciplineId).
		WillReturnResult(sqlmock.NewResult(1, 1))
	database.Mock.ExpectExec(`DELETE FROM "result_statics" WHERE discipline_id = \$1`).
		WithArgs(disciplineId).
		WillReturnResult(sqlmock.NewResult(1, 1))
	database.Mock.ExpectExec(`DELETE FROM "disciplines" WHERE "disciplines"."id" = \$1`).
		WithArgs(disciplineId).
		WillReturnResult(sqlmock.NewResult(1, 1))
	database.Mock.ExpectCommit()

	err = deleteDiscipline(ctx, disciplineId.String(), true)
	assert.Nil(t, err)
}

func TestDeleteDiscipline(t *testing.T) {
	clubId := uuid.New()
	disciplineId := uuid.New()
	authMap := getValidAuthMap(clubId)
	ctx := h.CreateContext(models.User{ID: "123", ClubID: clubId}, authMap, h.TwoFAState{})

	database.ConnectMockDB()

	database.Mock.ExpectQuery(`SELECT \* FROM "disciplines" WHERE.*"disciplines"."id" = \$1`).
		WithArgs(disciplineId, 1).
		WillReturnRows(sqlmock.NewRows([]string{"id", "club_id"}).AddRow(disciplineId, clubId))
	database.Mock.ExpectBegin()
	database.Mock.ExpectExec(`DELETE FROM "disciplines" WHERE "disciplines"."id" = \$1`).
		WithArgs(disciplineId).
		WillReturnResult(sqlmock.NewResult(1, 1))
	database.Mock.ExpectCommit()

	err := deleteDiscipline(ctx, disciplineId.String(), false)
	assert.Nil(t, err)
}

func getValidAuthMap(clubId uuid.UUID) map[uuid.UUID]map[h.Role]bool {
	authMap := make(map[uuid.UUID]map[h.Role]bool)
	authMap[clubId] = make(map[h.Role]bool)
	authMap[clubId][h.RoleCoach] = true

	return authMap
}
