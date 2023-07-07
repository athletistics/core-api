// Author  : Nicolas Baeriswyl <nicolasfabrice.baeriswyl@edu.hefr.ch>
// Created : 07.07.2023

package ndsHandler

import (
	"fmt"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"gitlab.forge.hefr.ch/tb-athletistics/backend/database"
	"gitlab.forge.hefr.ch/tb-athletistics/backend/database/types"
	"gitlab.forge.hefr.ch/tb-athletistics/backend/helpers"
	"gitlab.forge.hefr.ch/tb-athletistics/backend/internal/models"
	"net/http"
	"testing"
	"time"
)

func TestSelectInstances_BadParam(t *testing.T) {
	_, err := selectInstances(nil, "")

	assert.NotNil(t, err)
	assert.Equal(t, err.Error, http.StatusBadRequest)
}

func TestSelectInstances_Forbidden(t *testing.T) {
	database.ConnectMockDB()

	activityId, _ := uuid.NewUUID()
	clubId, _ := uuid.NewUUID()

	mockActivityAuthorisations(clubId, helpers.RoleAthlete)

	_, err := selectInstances(nil, activityId.String())

	assert.NotNil(t, err)
	assert.Equal(t, err.Error, http.StatusForbidden)
}

func TestSelectInstances(t *testing.T) {
	database.ConnectMockDB()

	expectedNbInstances := 5
	activityId, _ := uuid.NewUUID()
	clubId, _ := uuid.NewUUID()
	user := models.User{ID: "1"}

	mockActivityAuthorisations(clubId, helpers.RoleCoach)

	columns := []string{"id", "nds_activity_type_id"}
	database.Mock.ExpectQuery(`SELECT .* FROM "nds_activity_instances"`).WillReturnRows(sqlmock.NewRows(columns).
		AddRow(activityId, 1).AddRow(activityId, 1).AddRow(activityId, 1).AddRow(activityId, 1).AddRow(activityId, 1))
	database.Mock.ExpectQuery(`SELECT .* FROM "nds_activity_types"`).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

	instances, err := selectInstances(&user, activityId.String())

	assert.Nil(t, err)
	assert.Len(t, instances, expectedNbInstances)
}

func TestSelectInstance(t *testing.T) {
	activityId, _ := uuid.NewUUID()
	instanceId, _ := uuid.NewUUID()
	clubId, _ := uuid.NewUUID()
	user := models.User{ID: "1"}

	// Unauthorized
	database.ConnectMockDB()
	mockActivityAuthorisations(clubId, helpers.RoleAthlete)
	_, err := selectInstance(&user, activityId.String(), instanceId.String())
	assert.NotNil(t, err)
	assert.Equal(t, err.Error, http.StatusForbidden)

	// Success
	database.ConnectMockDB()
	mockActivityAuthorisations(clubId, helpers.RoleCoach)
	database.Mock.ExpectQuery(`SELECT .* FROM "nds_activity_instances"`).WillReturnRows(sqlmock.NewRows([]string{}).AddRow())
	database.Mock.ExpectQuery(`SELECT .* FROM "nds_activity_types"`).WillReturnRows(sqlmock.NewRows([]string{}).AddRow())
	_, err = selectInstance(&user, activityId.String(), instanceId.String())
	assert.Nil(t, err)
}

func TestCreateInstance(t *testing.T) {
	activityId, _ := uuid.NewUUID()
	clubId, _ := uuid.NewUUID()
	user := models.User{ID: "1"}
	body := models.NdsActivityInstanceInput{
		Date: types.NewDate(2023, time.July, 7),
	}

	// Unauthorized
	database.ConnectMockDB()
	mockActivityAuthorisations(clubId, helpers.RoleAthlete)
	_, err := createInstance(&user, activityId.String(), nil, nil)
	assert.NotNil(t, err)
	assert.Equal(t, err.Error, http.StatusForbidden)

	// Success
	database.ConnectMockDB()
	mockActivityAuthorisations(clubId, helpers.RoleCoach)
	database.Mock.ExpectQuery(`SELECT .* FROM "nds_activities"`).WillReturnRows(sqlmock.NewRows([]string{"start_date", "end_date"}).
		AddRow(types.NewDate(2023, time.July, 5), types.NewDate(2023, time.July, 10)))
	database.Mock.ExpectBegin()
	database.Mock.ExpectQuery(`INSERT INTO "nds_activity_instances"`).WillReturnRows(sqlmock.NewRows([]string{}).AddRow())
	database.Mock.ExpectCommit()

	_, err = createInstance(&user, activityId.String(), &body, nil)
	assert.Nil(t, err)
}

func TestDeleteInstance(t *testing.T) {
	activityId, _ := uuid.NewUUID()
	instanceId, _ := uuid.NewUUID()
	clubId, _ := uuid.NewUUID()
	user := models.User{ID: "1"}

	// Unauthorized
	database.ConnectMockDB()
	mockActivityAuthorisations(clubId, helpers.RoleAthlete)
	err := deleteInstance(&user, activityId.String(), instanceId.String())
	assert.NotNil(t, err)
	assert.Equal(t, err.Error, http.StatusForbidden)

	// Success
	database.ConnectMockDB()
	mockActivityAuthorisations(clubId, helpers.RoleCoach)
	database.Mock.ExpectQuery(`SELECT .* FROM "nds_activity_instances"`).WillReturnRows(sqlmock.NewRows([]string{}).AddRow())
	database.Mock.ExpectBegin()
	database.Mock.ExpectExec(`UPDATE "nds_attendances" SET "deleted_at"=\$1 WHERE "nds_attendances"\."nds_activity_instance_id"`).WillReturnResult(sqlmock.NewResult(0, 5))
	database.Mock.ExpectExec(`UPDATE "nds_activity_instances" SET "deleted_at"=\$1`).WillReturnResult(sqlmock.NewResult(0, 1))
	database.Mock.ExpectCommit()

	err = deleteInstance(&user, activityId.String(), instanceId.String())
	assert.Nil(t, err)
}

func TestInstanceExistsAndUserCanManage_Errors(t *testing.T) {
	clubId, _ := uuid.NewUUID()
	activityId, _ := uuid.NewUUID()
	instanceId, _ := uuid.NewUUID()
	user := models.User{ID: "1"}

	_, _, err := instanceExistsAndUserCanManage(&user, "", "")
	assert.NotNil(t, err)
	assert.Equal(t, http.StatusBadRequest, err.Error)

	database.ConnectMockDB()
	mockActivityAuthorisations(clubId, helpers.RoleAthlete)

	_, _, err = instanceExistsAndUserCanManage(&user, activityId.String(), instanceId.String())
	assert.NotNil(t, err)
	assert.Equal(t, http.StatusForbidden, err.Error)
}

func TestInstanceExistsAndUserCanManage(t *testing.T) {
	clubId, _ := uuid.NewUUID()
	activityId, _ := uuid.NewUUID()
	instanceId, _ := uuid.NewUUID()
	user := models.User{ID: "1"}

	database.ConnectMockDB()
	mockActivityAuthorisations(clubId, helpers.RoleCoach)
	database.Mock.ExpectQuery(`SELECT .* FROM "nds_activity_instances"`).WillReturnRows(sqlmock.NewRows([]string{}).AddRow())

	aId, iId, err := instanceExistsAndUserCanManage(&user, activityId.String(), instanceId.String())
	assert.Nil(t, err)
	assert.Equal(t, activityId.String(), aId.String())
	assert.Equal(t, instanceId.String(), iId.String())
}

func TestParseUuids(t *testing.T) {
	_, err := parseUuids("")
	assert.NotNil(t, err)

	_, err = parseUuids("'")
	assert.NotNil(t, err)

	id, _ := uuid.NewUUID()
	ids, err := parseUuids(id.String())
	assert.Nil(t, err)
	assert.Len(t, ids, 1)

	id1, _ := uuid.NewUUID()
	id2, _ := uuid.NewUUID()
	ids, err = parseUuids(id1.String(), id2.String())
	assert.Nil(t, err)
	assert.Len(t, ids, 2)
	assert.Equal(t, id1.String(), ids[0].String())
	assert.Equal(t, id2.String(), ids[1].String())
}

func mockActivityAuthorisations(clubId uuid.UUID, role helpers.Role) {
	// access permission queries
	database.Mock.ExpectQuery(`SELECT .* FROM "teams" inner join nds_activities`).WillReturnRows(sqlmock.NewRows([]string{"club_id"}).
		FromCSVString(fmt.Sprintf("%v", clubId)))
	database.Mock.ExpectQuery(`SELECT .* FROM "authorisations"`).WillReturnRows(sqlmock.NewRows([]string{"club_id", "access_role_id"}).
		FromCSVString(fmt.Sprintf("%v,%v", clubId, 2)))
	database.Mock.ExpectQuery(`SELECT .* FROM "access_roles"`).WillReturnRows(sqlmock.NewRows([]string{"id", "i18n_name"}).
		FromCSVString(fmt.Sprintf("%v,%v", 2, role)))
}
