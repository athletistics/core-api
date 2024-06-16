// Author  : Nicolas Baeriswyl <nicolasfabrice.baeriswyl@edu.hefr.ch>
// Created : 11.07.2023

package userHandler

import (
	"database/sql/driver"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/athletistics/core-api/database"
	"github.com/athletistics/core-api/helpers"
	"github.com/athletistics/core-api/internal/models"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
	"net/http"
	"testing"
	"time"
)

const userId = "123"

var clubId = uuid.New()
var code = uuid.New()
var requestingUser = models.User{ID: "321", ClubID: clubId}
var userCols = []string{"id", "club_id"}
var shareCols = []string{"id", "user_id"}
var authCols = []string{"user_id", "club_id", "access_role_id"}

func TestCreateShareCode_Success(t *testing.T) {
	database.ConnectMockDB()

	database.Mock.ExpectQuery(`SELECT \* FROM "users" WHERE.*"users"."id" = \$1`).WithArgs(userId, 1).
		WillReturnRows(sqlmock.NewRows(userCols).AddRow(userId, clubId))
	mockHasRoleOnParentClub(true)

	// code creation
	database.Mock.ExpectBegin()
	database.Mock.ExpectQuery(`INSERT INTO "user_shares"`).WillReturnRows(sqlmock.NewRows([]string{}).AddRow())
	database.Mock.ExpectCommit()

	_, err := createShareCode(&requestingUser, userId)
	assert.Nil(t, err)
}

func TestCreateShareCode_Forbidden(t *testing.T) {
	database.ConnectMockDB()

	database.Mock.ExpectQuery(`SELECT \* FROM "users" WHERE.*"users"."id" = \$1`).WithArgs(userId, 1).
		WillReturnRows(sqlmock.NewRows(userCols).AddRow(userId, clubId))
	mockHasRoleOnParentClub(false)

	_, err := createShareCode(&requestingUser, userId)
	assert.NotNil(t, err)
	assert.Equal(t, http.StatusForbidden, err.Error)
}

func TestCreateShareCode_UnknownUser(t *testing.T) {
	database.ConnectMockDB()

	database.Mock.ExpectQuery(`SELECT \* FROM "users" WHERE.*"users"."id" = \$1`).WithArgs(userId, 1).
		WillReturnRows(sqlmock.NewRows(userCols))

	_, err := createShareCode(&requestingUser, userId)
	assert.NotNil(t, err)
	assert.Equal(t, http.StatusNotFound, err.Error)
}

func TestRedeemShareCode_Success(t *testing.T) {
	mockHasRoleOnParentClub(true)
	database.Mock.ExpectQuery(`SELECT \* FROM "user_shares" WHERE.*"user_shares"."id" = \$2`).WithArgs(AnyTime{}, code, 1).
		WillReturnRows(sqlmock.NewRows(shareCols).AddRow(code, userId))

	// auth creation
	database.Mock.ExpectBegin()
	database.Mock.ExpectExec(`INSERT INTO "authorisations"`).WithArgs(models.AccessRoleIdAthlete, userId, clubId, AnyTime{}, AnyTime{}, nil).
		WillReturnResult(sqlmock.NewResult(0, 1))
	database.Mock.ExpectExec(`DELETE FROM "user_shares"`).WillReturnResult(sqlmock.NewResult(0, 1))
	database.Mock.ExpectCommit()

	// read user for return
	database.Mock.ExpectQuery(`SELECT.*FROM "users" WHERE.*"users"."id" = \$1`).WithArgs(userId, 1).
		WillReturnRows(sqlmock.NewRows(userCols).AddRow(userId, clubId))

	user, err := redeemShareCode(&requestingUser, code.String())
	assert.Nil(t, err)
	assert.Equal(t, userId, user.ID)
}

func TestRedeemShareCode_SuccessDuplicate(t *testing.T) {
	mockHasRoleOnParentClub(true)
	database.Mock.ExpectQuery(`SELECT \* FROM "user_shares" WHERE.*"user_shares"."id" = \$2`).WithArgs(AnyTime{}, code, 1).
		WillReturnRows(sqlmock.NewRows(shareCols).AddRow(code, userId))

	// auth creation
	database.Mock.ExpectBegin()
	database.Mock.ExpectExec(`INSERT INTO "authorisations"`).WillReturnError(gorm.ErrDuplicatedKey)
	database.Mock.ExpectRollback()

	// read user for return
	database.Mock.ExpectQuery(`SELECT.*FROM "users" WHERE.*"users"."id" = \$1`).WithArgs(userId, 1).
		WillReturnRows(sqlmock.NewRows(userCols).AddRow(userId, clubId))

	user, err := redeemShareCode(&requestingUser, code.String())
	assert.Nil(t, err)
	assert.Equal(t, userId, user.ID)
}

func TestRedeemShareCode_Forbidden(t *testing.T) {
	mockHasRoleOnParentClub(false)
	_, err := redeemShareCode(&requestingUser, code.String())
	assert.NotNil(t, err)
	assert.Equal(t, http.StatusForbidden, err.Error)
}

func TestRedeemShareCode_Expired(t *testing.T) {
	mockHasRoleOnParentClub(true)
	database.Mock.ExpectQuery(`SELECT \* FROM "user_shares" WHERE.*"user_shares"."id" = \$2`).
		WillReturnError(gorm.ErrRecordNotFound)

	_, err := redeemShareCode(&requestingUser, code.String())
	assert.NotNil(t, err)
	assert.Equal(t, http.StatusNotFound, err.Error)
}

type AnyTime struct{}

// Match satisfies sqlmock.Argument interface
func (a AnyTime) Match(v driver.Value) bool {
	_, ok := v.(time.Time)
	return ok
}

func mockHasRoleOnParentClub(success bool) {
	if success {
		database.Mock.ExpectQuery(`SELECT \* FROM "authorisations"`).
			WillReturnRows(sqlmock.NewRows(authCols).AddRow(requestingUser.ID, clubId, int(helpers.RoleHeadCoach)))
		database.Mock.ExpectQuery(`SELECT .* FROM "access_roles"`).
			WillReturnRows(sqlmock.NewRows([]string{"id", "i18n_name"}).AddRow(models.AccessRoleIdHeadCoach, helpers.RoleHeadCoach.String()))
	} else {
		database.Mock.ExpectQuery(`SELECT \* FROM "authorisations"`).
			WillReturnRows(sqlmock.NewRows(authCols))
	}
}
