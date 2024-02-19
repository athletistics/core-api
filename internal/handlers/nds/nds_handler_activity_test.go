// Author  : Nicolas Baeriswyl <nicolasfabrice.baeriswyl@edu.hefr.ch>
// Created : 07.07.2023

package ndsHandler

import (
	"fmt"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/athletistics/core-api/database"
	"github.com/athletistics/core-api/database/types"
	"github.com/athletistics/core-api/helpers"
	"github.com/athletistics/core-api/internal/models"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
	"time"
)

var user = models.User{ID: "123"}

func TestSelectActivities_BadParams(t *testing.T) {
	_, err := selectActivities(&user, "")
	assert.NotNil(t, err)
	assert.Equal(t, err.Error, http.StatusBadRequest)

	_, err = selectActivities(&user, "\"")
	assert.NotNil(t, err)
	assert.Equal(t, err.Error, http.StatusBadRequest)

	id, _ := uuid.NewUUID()
	_, err = selectActivities(&user, fmt.Sprintf("%vb", id))
	assert.NotNil(t, err)
	assert.Equal(t, err.Error, http.StatusBadRequest)
}

func TestSelectActivities_Forbidden(t *testing.T) {
	database.ConnectMockDB()
	clubId, _ := uuid.NewUUID()

	database.Mock.ExpectQuery(`SELECT .* FROM "teams"`).WillReturnRows(sqlmock.NewRows([]string{"club_id"}).
		FromCSVString(fmt.Sprintf("%v", clubId)))
	database.Mock.ExpectQuery(`SELECT .* FROM "authorisations"`).WillReturnRows(sqlmock.NewRows([]string{"club_id", "access_role_id"}))

	id, _ := uuid.NewUUID()
	_, err := selectActivities(&user, id.String())
	assert.NotNil(t, err)
	assert.Equal(t, err.Error, http.StatusForbidden)
}

func TestSelectActivities(t *testing.T) {
	database.ConnectMockDB()
	clubId, _ := uuid.NewUUID()
	user := models.User{ID: "1"}

	mockTeamAuthorisations(clubId, helpers.RoleHeadCoach)

	database.Mock.ExpectQuery(`SELECT .* FROM "nds_activities"`).WillReturnRows(sqlmock.NewRows([]string{}).AddRow())
	database.Mock.ExpectQuery(`SELECT .* FROM "teams"`).WillReturnRows(sqlmock.NewRows([]string{}).AddRow())

	id, _ := uuid.NewUUID()
	_, err := selectActivities(&user, id.String())
	assert.Nil(t, err)
}

func TestSelectActivity(t *testing.T) {
	database.ConnectMockDB()
	clubId, _ := uuid.NewUUID()
	teamId, _ := uuid.NewUUID()
	user := models.User{ID: "1"}

	// activity queries
	database.Mock.ExpectQuery(`SELECT .* FROM "nds_activities"`).WillReturnRows(sqlmock.NewRows([]string{"team_id"}).
		FromCSVString(fmt.Sprintf("%v", teamId)))
	database.Mock.ExpectQuery(`SELECT .* FROM "teams"`).WillReturnRows(sqlmock.NewRows([]string{"id"}).FromCSVString(teamId.String()))

	mockTeamAuthorisations(clubId, helpers.RoleCoach)

	id, _ := uuid.NewUUID()
	_, err := selectActivity(&user, id.String())
	assert.Nil(t, err)
}

func TestCreateActivity_BadParams(t *testing.T) {
	_, err := createActivity(&user, nil, fmt.Errorf(""))
	assert.NotNil(t, err)
	assert.Equal(t, err.Error, http.StatusBadRequest)
}

func TestCreateActivity_Forbidden(t *testing.T) {
	database.ConnectMockDB()
	clubId, _ := uuid.NewUUID()
	teamId, _ := uuid.NewUUID()
	body := models.NdsActivityInput{TeamID: teamId}
	user := models.User{ID: "1"}

	mockTeamAuthorisations(clubId, helpers.RoleCoach)

	_, err := createActivity(&user, &body, nil)
	assert.NotNil(t, err)
	assert.Equal(t, err.Error, http.StatusForbidden)
}

func TestCreateActivity_NoInstances(t *testing.T) {
	database.ConnectMockDB()
	clubId, _ := uuid.NewUUID()
	teamId, _ := uuid.NewUUID()
	user := models.User{ID: "1"}

	mockTeamAuthorisations(clubId, helpers.RoleHeadCoach)

	body := models.NdsActivityInput{
		StartDate: types.NewDate(2023, time.July, 7),
		EndDate:   types.NewDate(2023, time.July, 10),
		StartTime: types.NewTime(10, 0, 0),
		Duration:  90,
		Weekday:   time.Wednesday,
		Place:     "Plaffeien",
		TeamID:    teamId,
	}

	database.Mock.ExpectBegin()
	database.Mock.ExpectQuery(`INSERT INTO "nds_activities"`).WillReturnRows(sqlmock.NewRows([]string{}).AddRow())
	database.Mock.ExpectCommit()

	_, err := createActivity(&user, &body, nil)
	assert.Nil(t, err)
}

func TestCreateActivity_WithInstances(t *testing.T) {
	database.ConnectMockDB()
	expectedNbInstances := 5
	clubId, _ := uuid.NewUUID()
	teamId, _ := uuid.NewUUID()
	user := models.User{ID: "1"}

	mockTeamAuthorisations(clubId, helpers.RoleHeadCoach)

	body := models.NdsActivityInput{
		StartDate: types.NewDate(2023, time.July, 5),
		EndDate:   types.NewDate(2023, time.August, 2),
		StartTime: types.NewTime(10, 0, 0),
		Duration:  90,
		Weekday:   time.Wednesday,
		Place:     "Plaffeien",
		TeamID:    teamId,
	}

	database.Mock.ExpectBegin()
	database.Mock.ExpectQuery(`INSERT INTO "nds_activities"`).WillReturnRows(sqlmock.NewRows([]string{}).AddRow())

	for i := 0; i < expectedNbInstances; i++ {
		database.Mock.ExpectQuery(`INSERT INTO "nds_activity_instances"`).WillReturnRows(sqlmock.NewRows([]string{}).AddRow())
	}

	database.Mock.ExpectCommit()

	_, err := createActivity(&user, &body, nil)
	assert.Nil(t, err)
}

func mockTeamAuthorisations(clubId uuid.UUID, role helpers.Role) {
	// access permission queries
	database.Mock.ExpectQuery(`SELECT .* FROM "teams"`).WillReturnRows(sqlmock.NewRows([]string{"club_id"}).
		FromCSVString(fmt.Sprintf("%v", clubId)))
	database.Mock.ExpectQuery(`SELECT .* FROM "authorisations"`).WillReturnRows(sqlmock.NewRows([]string{"club_id", "access_role_id"}).
		FromCSVString(fmt.Sprintf("%v,%v", clubId, 2)))
	database.Mock.ExpectQuery(`SELECT .* FROM "access_roles"`).WillReturnRows(sqlmock.NewRows([]string{"id", "i18n_name"}).
		FromCSVString(fmt.Sprintf("%v,%v", 2, role)))
}
