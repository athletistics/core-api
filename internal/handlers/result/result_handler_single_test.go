// Author  : Nicolas Baeriswyl <nicolasfabrice.baeriswyl@edu.hefr.ch>
// Created : 21.06.2024

package resultHandler

import (
	"fmt"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/athletistics/core-api/database"
	h "github.com/athletistics/core-api/helpers"
	"github.com/athletistics/core-api/internal/models"
	res "github.com/athletistics/core-api/internal/responses"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCreatePersonalTrainingResult_NotAthlete(t *testing.T) {
	u := models.User{ID: "123"}
	authMap := make(map[uuid.UUID]map[h.Role]bool)
	ctx := h.CreateContext(u, authMap, h.TwoFAState{})
	body := &models.ResultTrainingInput{}

	err := createTrainingResult(ctx, body)

	assert.NotNil(t, err)
	assert.Equal(t, res.ErrForbidden.Error, err.Error)
}

func TestCreatePersonalTrainingResult(t *testing.T) {
	clubId := uuid.New()
	discId := uuid.New()
	u := models.User{ID: "123", ClubID: clubId}
	authMap := make(map[uuid.UUID]map[h.Role]bool)
	authMap[clubId] = make(map[h.Role]bool)
	authMap[clubId][h.RoleAthlete] = true

	ctx := h.CreateContext(u, authMap, h.TwoFAState{})
	body := &models.ResultTrainingInput{DisciplineID: discId}

	database.ConnectMockDB()
	database.Mock.ExpectQuery(`SELECT .* FROM "disciplines" WHERE.*club_id = \$1.*`).WithArgs(clubId, discId, 1).
		WillReturnRows(sqlmock.NewRows([]string{"id", "club_id"}).AddRow(discId, ""))
	database.Mock.ExpectBegin()
	database.Mock.ExpectQuery(`INSERT INTO "result_trainings".*`).WillReturnRows(sqlmock.NewRows([]string{}).AddRow())
	database.Mock.ExpectCommit()

	err := createTrainingResult(ctx, body)

	assert.Nil(t, err)
}

func TestUpdatePersonalTrainingResult_InvalidID(t *testing.T) {
	u := models.User{ID: "123"}
	authMap := make(map[uuid.UUID]map[h.Role]bool)
	ctx := h.CreateContext(u, authMap, h.TwoFAState{})
	body := &models.ResultTrainingInput{}

	assert.NotNil(t, updateTrainingResult(ctx, "", body))
	assert.NotNil(t, updateTrainingResult(ctx, "123", body))
	assert.NotNil(t, updateTrainingResult(ctx, fmt.Sprintf("%s0", uuid.New().String()), body))
}

func TestUpdatePersonalTrainingResult_NotFound(t *testing.T) {
	clubId := uuid.New()
	discId := uuid.New()
	resultId := uuid.New().String()
	u := models.User{ID: "123", ClubID: clubId}
	authMap := make(map[uuid.UUID]map[h.Role]bool)
	authMap[clubId] = make(map[h.Role]bool)
	authMap[clubId][h.RoleAthlete] = true

	ctx := h.CreateContext(u, authMap, h.TwoFAState{})
	body := &models.ResultTrainingInput{DisciplineID: discId}

	database.ConnectMockDB()
	database.Mock.ExpectQuery(`SELECT .* FROM "result_trainings".*`).WillReturnRows(sqlmock.NewRows([]string{}))

	err := updateTrainingResult(ctx, resultId, body)

	assert.NotNil(t, err)
	assert.Equal(t, res.ErrNotFound.Error, err.Error)
}

func TestUpdatePersonalTrainingResult_Athlete(t *testing.T) {
	testUpdatePersonalTrainingResult(t, h.RoleAthlete)
}

func TestUpdatePersonalTrainingResult_Coach(t *testing.T) {
	testUpdatePersonalTrainingResult(t, h.RoleCoach)
}

func testUpdatePersonalTrainingResult(t *testing.T, r h.Role) {
	clubId := uuid.New()
	discId := uuid.New()
	resultId := uuid.New().String()
	u := models.User{ID: "123", ClubID: clubId}
	userId := u.ID
	authMap := make(map[uuid.UUID]map[h.Role]bool)
	authMap[clubId] = make(map[h.Role]bool)
	authMap[clubId][r] = true

	if r == h.RoleCoach {
		userId = "234"
	}

	ctx := h.CreateContext(u, authMap, h.TwoFAState{})
	body := &models.ResultTrainingInput{DisciplineID: discId}

	database.ConnectMockDB()
	database.Mock.ExpectQuery(`SELECT .* FROM "result_trainings".*`).
		WillReturnRows(sqlmock.NewRows([]string{"user_id"}).AddRow(userId))
	database.Mock.ExpectQuery(`SELECT .* FROM "users"`).
		WillReturnRows(sqlmock.NewRows([]string{"id", "club_id"}).AddRow(userId, clubId))
	database.Mock.ExpectQuery(`SELECT .* FROM "disciplines"`).
		WillReturnRows(sqlmock.NewRows([]string{"id", "club_id"}).AddRow(discId, ""))
	database.Mock.ExpectBegin()
	database.Mock.ExpectExec(`UPDATE "result_trainings".*`).WillReturnResult(sqlmock.NewResult(1, 1))
	database.Mock.ExpectCommit()

	err := updateTrainingResult(ctx, resultId, body)

	assert.Nil(t, err)
}

func TestDeletePersonalTrainingResult_InvalidID(t *testing.T) {
	u := models.User{ID: "123"}
	authMap := make(map[uuid.UUID]map[h.Role]bool)
	ctx := h.CreateContext(u, authMap, h.TwoFAState{})

	assert.NotNil(t, deleteTrainingResult(ctx, ""))
	assert.NotNil(t, deleteTrainingResult(ctx, "123"))
	assert.NotNil(t, deleteTrainingResult(ctx, fmt.Sprintf("%s0", uuid.New().String())))
}

func TestDeletePersonalTrainingResult_Athlete(t *testing.T) {
	testDeletePersonalTrainingResult(t, h.RoleAthlete)
}

func TestDeletePersonalTrainingResult_Coach(t *testing.T) {
	testDeletePersonalTrainingResult(t, h.RoleCoach)
}

func testDeletePersonalTrainingResult(t *testing.T, r h.Role) {
	clubId := uuid.New()
	resultId := uuid.New().String()
	u := models.User{ID: "123", ClubID: clubId}
	userId := u.ID
	authMap := make(map[uuid.UUID]map[h.Role]bool)
	authMap[clubId] = make(map[h.Role]bool)
	authMap[clubId][r] = true

	if r == h.RoleCoach {
		userId = "234"
	}

	ctx := h.CreateContext(u, authMap, h.TwoFAState{})

	database.ConnectMockDB()
	database.Mock.ExpectQuery(`SELECT .* FROM "result_trainings".*`).
		WillReturnRows(sqlmock.NewRows([]string{"user_id"}).AddRow(userId))
	database.Mock.ExpectQuery(`SELECT .* FROM "users"`).
		WillReturnRows(sqlmock.NewRows([]string{"id", "club_id"}).AddRow(userId, clubId))
	database.Mock.ExpectBegin()
	database.Mock.ExpectExec(`DELETE FROM "result_trainings".*`).WillReturnResult(sqlmock.NewResult(1, 1))
	database.Mock.ExpectCommit()

	err := deleteTrainingResult(ctx, resultId)

	assert.Nil(t, err)
}
