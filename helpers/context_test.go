// Author  : Nicolas Baeriswyl <nicolasfabrice.baeriswyl@edu.hefr.ch>
// Created : 14.06.2024

package helpers

import (
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/athletistics/core-api/database"
	"github.com/athletistics/core-api/internal/models"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetUser_UnknownId(t *testing.T) {
	database.ConnectMockDB()
	database.Mock.ExpectQuery(`SELECT \* FROM "users" WHERE .*"id" = \$1.*`).
		WillReturnRows(sqlmock.NewRows([]string{}))

	u := getUser(userId)

	assert.Nil(t, u)
}

func TestGetUser_Normal(t *testing.T) {
	database.ConnectMockDB()
	database.Mock.ExpectQuery(`SELECT \* FROM "users" WHERE .*"id" = \$1.*`).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(userId))

	u := getUser(userId)

	assert.NotNil(t, u)
	assert.Equal(t, u.ID, userId)
}

func TestGetAuthsByClub_NoUser(t *testing.T) {
	roles := getAuthsByClub(nil)

	assert.Len(t, roles, 0)
}

func TestGetAuthsByClub_NoAuths(t *testing.T) {
	database.ConnectMockDB()
	database.Mock.ExpectQuery(`SELECT \* FROM "authorisations" WHERE .*"user_id" = \$1.*`).WillReturnRows(sqlmock.NewRows([]string{}))
	database.Mock.ExpectQuery(`SELECT \* FROM "access_roles".*`).WillReturnRows(sqlmock.NewRows([]string{}))

	roles := getAuthsByClub(&models.User{ID: userId})

	assert.Len(t, roles, 0)
}

func TestGetAuthsByClub_ManyAuths(t *testing.T) {
	clubId2, _ := uuid.NewUUID()
	clubId2Copy, _ := uuid.Parse(clubId2.String())
	unusedClubId, _ := uuid.NewUUID()

	database.ConnectMockDB()
	database.Mock.ExpectQuery(`SELECT \* FROM "authorisations" WHERE .*"user_id" = \$1.*`).
		WillReturnRows(sqlmock.NewRows([]string{"user_id", "club_id", "access_role_id"}).
			AddRow(userId, clubId, "1").
			AddRow(userId, clubId, "2").
			AddRow(userId, clubId, "3").
			AddRow(userId, clubId2, "1").
			AddRow(userId, clubId2Copy, "2"))
	database.Mock.ExpectQuery(`SELECT \* FROM "access_roles".*`).
		WillReturnRows(sqlmock.NewRows([]string{"id", "i18n_name"}).
			AddRow("1", RoleAthlete.String()).
			AddRow("2", RoleHeadCoach.String()).
			AddRow("3", RoleCoach.String()))

	roles := getAuthsByClub(&models.User{ID: userId})

	assert.Len(t, roles, 2)
	assert.Len(t, roles[clubId], 3)
	assert.Len(t, roles[clubId2], 2)
	assert.Len(t, roles[unusedClubId], 0)
}

func TestHasRoleOnClub_NotAuthenticated(t *testing.T) {
	c := &Context{}
	assert.False(t, c.HasRoleOnClub(clubId, RoleAny))
}

func TestHasRoleOnClub_NoRole(t *testing.T) {
	c := &Context{user: &models.User{ID: userId}}
	assert.False(t, c.HasRoleOnClub(clubId, RoleAny))
}

func TestHasRoleOnClub_AnyRole(t *testing.T) {
	clubId2, _ := uuid.NewUUID()

	authByClub := make(map[uuid.UUID]map[Role]bool)
	authByClub[clubId] = make(map[Role]bool)

	c := &Context{user: &models.User{ID: userId}, authByClub: authByClub}
	assert.False(t, c.HasRoleOnClub(clubId, RoleAny))

	authByClub[clubId][RoleAthlete] = true

	c = &Context{user: &models.User{ID: userId}, authByClub: authByClub}
	assert.False(t, c.HasRoleOnClub(clubId2, RoleAny))
	assert.True(t, c.HasRoleOnClub(clubId, RoleAny))
}

func TestHasRoleOnClub_SpecificRoles(t *testing.T) {
	clubId2, _ := uuid.NewUUID()

	authByClub := make(map[uuid.UUID]map[Role]bool)
	authByClub[clubId] = make(map[Role]bool)
	authByClub[clubId2] = make(map[Role]bool)
	authByClub[clubId][RoleAthlete] = true
	authByClub[clubId2][RoleCoach] = true

	c := &Context{user: &models.User{ID: userId}, authByClub: authByClub}
	assert.True(t, c.HasRoleOnClub(clubId, RoleAny))
	assert.True(t, c.HasRoleOnClub(clubId, RoleAthlete, RoleCoach))
	assert.True(t, c.HasRoleOnClub(clubId, RoleAthlete))
	assert.False(t, c.HasRoleOnClub(clubId, RoleCoach))
	assert.True(t, c.HasRoleOnClub(clubId2, RoleCoach))
}
