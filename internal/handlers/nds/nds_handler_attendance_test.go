// Author  : Nicolas Baeriswyl <nicolasfabrice.baeriswyl@edu.hefr.ch>
// Created : 07.07.2023

package ndsHandler

import (
	"fmt"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/athletistics/core-api/database"
	"github.com/athletistics/core-api/helpers"
	"github.com/athletistics/core-api/internal/models"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
	"net/http"
	"testing"
)

func TestSelectAttendances(t *testing.T) {
	activityId, _ := uuid.NewUUID()
	instanceId, _ := uuid.NewUUID()
	clubId, _ := uuid.NewUUID()
	user := models.User{ID: "1"}

	database.ConnectMockDB()
	mockActivityAuthorisations(clubId, helpers.RoleCoach)
	database.Mock.ExpectQuery(`SELECT .* FROM "nds_activity_instances"`).WillReturnRows(sqlmock.NewRows([]string{}).AddRow())
	database.Mock.ExpectQuery(`SELECT .* FROM "nds_attendances" WHERE "nds_attendances"."nds_activity_instance_id"`).WillReturnRows(sqlmock.NewRows([]string{}).AddRow())
	_, err := selectAttendances(&user, activityId.String(), instanceId.String())
	assert.Nil(t, err)
}

func TestSetParticipants(t *testing.T) {
	activityId, _ := uuid.NewUUID()
	instanceId, _ := uuid.NewUUID()
	clubId, _ := uuid.NewUUID()
	user := models.User{ID: "1"}
	body := make([]models.NdsAttendanceInput, 2)

	for i := 0; i < len(body); i++ {
		body[i] = models.NdsAttendanceInput{UserID: fmt.Sprintf("%v", i), IsCoach: true}
	}

	// Unauthorized
	database.ConnectMockDB()
	mockActivityAuthorisations(clubId, helpers.RoleAthlete)
	err := setAttendances(&user, activityId.String(), instanceId.String(), &body, nil)
	assert.NotNil(t, err)
	assert.Equal(t, err.Error, http.StatusForbidden)

	// Regular
	database.ConnectMockDB()
	mockActivityAuthorisations(clubId, helpers.RoleCoach)
	database.Mock.ExpectQuery(`SELECT .* FROM "nds_activity_instances"`).WillReturnRows(sqlmock.NewRows([]string{}).AddRow())
	database.Mock.ExpectBegin()
	database.Mock.ExpectQuery(`SELECT .* FROM "teams" inner join nds_activities`).WillReturnRows(sqlmock.NewRows([]string{}).AddRow())
	database.Mock.ExpectQuery(`SELECT .* FROM "authorisations"`).WillReturnRows(sqlmock.NewRows([]string{`COUNT(DISTINCT("user_id"))`}).AddRow(len(body)))
	database.Mock.ExpectExec(`DELETE FROM "nds_attendances" WHERE "nds_attendances"."nds_activity_instance_id"`).WillReturnResult(sqlmock.NewResult(0, 1))
	database.Mock.ExpectExec(`INSERT INTO "nds_attendances"`).WillReturnResult(sqlmock.NewResult(int64(len(body)), int64(len(body))))
	database.Mock.ExpectCommit()

	err = setAttendances(&user, activityId.String(), instanceId.String(), &body, nil)
	assert.Nil(t, err)
}

func TestSetParticipants_ErrorInTransaction(t *testing.T) {
	activityId, _ := uuid.NewUUID()
	instanceId, _ := uuid.NewUUID()
	clubId, _ := uuid.NewUUID()
	user := models.User{ID: "1"}
	body := []models.NdsAttendanceInput{
		models.NdsAttendanceInput{},
		models.NdsAttendanceInput{},
	}

	database.ConnectMockDB()
	mockActivityAuthorisations(clubId, helpers.RoleCoach)
	database.Mock.ExpectQuery(`SELECT .* FROM "nds_activity_instances"`).WillReturnRows(sqlmock.NewRows([]string{}).AddRow())
	database.Mock.ExpectBegin()
	database.Mock.ExpectQuery(`DELETE .* FROM "nds_attendances" WHERE "nds_attendances"."nds_activity_instance_id"`).WillReturnError(gorm.ErrDuplicatedKey)
	database.Mock.ExpectRollback()

	err := setAttendances(&user, activityId.String(), instanceId.String(), &body, nil)
	assert.NotNil(t, err)

	database.ConnectMockDB()
	mockActivityAuthorisations(clubId, helpers.RoleCoach)
	database.Mock.ExpectQuery(`SELECT .* FROM "nds_activity_instances"`).WillReturnRows(sqlmock.NewRows([]string{}).AddRow())
	database.Mock.ExpectBegin()
	database.Mock.ExpectExec(`DELETE FROM "nds_attendances" WHERE "nds_attendances"."nds_activity_instance_id"`).WillReturnResult(sqlmock.NewResult(0, 1))
	database.Mock.ExpectExec(`INSERT INTO "nds_attendances"`).WillReturnResult(sqlmock.NewResult(int64(len(body)-1), 0))
	database.Mock.ExpectRollback()

	err = setAttendances(&user, activityId.String(), instanceId.String(), &body, nil)
	assert.NotNil(t, err)
}
