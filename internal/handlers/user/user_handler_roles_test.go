// Author  : Nicolas Baeriswyl <nicolasfabrice.baeriswyl@edu.hefr.ch>
// Created : 11.07.2023

package userHandler

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

func TestGetRoles(t *testing.T) {
	database.ConnectMockDB()

	mockHasRoleOnParentClub(true)
	database.Mock.ExpectQuery(`SELECT \* FROM "authorisations" WHERE.*"authorisations"."user_id" = \$1 AND "authorisations"."club_id" = \$2`).
		WillReturnRows(sqlmock.NewRows(authCols).
			AddRow(userId, requestingUser.ClubID, models.AccessRoleIdAthlete).
			AddRow(userId, requestingUser.ClubID, models.AccessRoleIdCoach))
	database.Mock.ExpectQuery(`SELECT .* FROM "access_roles"`).
		WillReturnRows(sqlmock.NewRows([]string{"id", "i18n_name"}).
			AddRow(models.AccessRoleIdAthlete, helpers.RoleAthlete.String()).
			AddRow(models.AccessRoleIdCoach, helpers.RoleCoach.String()))

	data, err := getRoles(&requestingUser, userId)
	assert.Nil(t, err)
	assert.NotNil(t, data)
	assert.Len(t, data, 2)
}

func TestGetRoles_Forbidden(t *testing.T) {
	database.ConnectMockDB()
	mockHasRoleOnParentClub(false)

	_, err := getRoles(&requestingUser, userId)
	assert.NotNil(t, err)
}

func TestAddRole(t *testing.T) {
	database.ConnectMockDB()
	body := models.Authorisation{AccessRoleID: models.AccessRoleIdAthlete, UserID: userId, ClubID: requestingUser.ClubID}

	mockHasRoleOnParentClub(true)
	database.Mock.ExpectQuery(`SELECT \* FROM "users".*WHERE "users"."id" = \$1 AND "users"."club_id" = \$2`).
		WillReturnRows(sqlmock.NewRows(userCols).AddRow(userId, requestingUser.ClubID))

	database.Mock.ExpectBegin()
	database.Mock.ExpectExec(`INSERT INTO "authorisations"`).WillReturnResult(sqlmock.NewResult(0, 1))
	database.Mock.ExpectCommit()

	err := addRole(&requestingUser, userId, &body, nil)
	assert.Nil(t, err)
}

func TestAddRole_AmbiguousBody(t *testing.T) {
	err := addRole(&requestingUser, "", nil, fmt.Errorf(""))
	assert.NotNil(t, err)
	assert.Equal(t, http.StatusBadRequest, err.Error)

	body := models.Authorisation{
		UserID: userId,
		ClubID: uuid.New(),
	}
	err = addRole(&requestingUser, userId, &body, nil)
	assert.NotNil(t, err)
	assert.Equal(t, http.StatusBadRequest, err.Error)

	body.UserID = "1"
	body.ClubID = requestingUser.ClubID
	err = addRole(&requestingUser, userId, &body, nil)
	assert.NotNil(t, err)
	assert.Equal(t, http.StatusBadRequest, err.Error)
}

func TestAddRole_Forbidden(t *testing.T) {
	database.ConnectMockDB()
	mockHasRoleOnParentClub(false)

	body := models.Authorisation{UserID: userId, ClubID: requestingUser.ClubID}
	err := addRole(&requestingUser, userId, &body, nil)
	assert.NotNil(t, err)
	assert.Equal(t, http.StatusForbidden, err.Error)

	// no user found with given id and club
	database.ConnectMockDB()

	mockHasRoleOnParentClub(true)
	database.Mock.ExpectQuery(`SELECT \* FROM "users".*WHERE "users"."id" = \$1 AND "users"."club_id" = \$2`).
		WillReturnError(gorm.ErrRecordNotFound)
	err = addRole(&requestingUser, userId, &body, nil)
	assert.NotNil(t, err)
	assert.Equal(t, http.StatusNotFound, err.Error)
}
