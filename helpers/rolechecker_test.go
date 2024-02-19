// Brief   :
// Author  : Nicolas Baeriswyl <nicolasfabrice.baeriswyl@edu.hefr.ch>
// Created : 26.06.2023

package helpers

import (
	"fmt"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/athletistics/core-api/database"
	"github.com/athletistics/core-api/internal/models"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"testing"
)

var (
	userId    = "123"
	clubId, _ = uuid.NewUUID()
)

func TestGetAuthorisations_NoRoles(t *testing.T) {
	database.ConnectMockDB()

	database.Mock.ExpectQuery(`SELECT`).WillReturnRows(sqlmock.NewRows([]string{}))
	auth, err := GetAuthorisations(userId, clubId)

	assert.Nil(t, err)
	assert.Empty(t, auth)
}

func TestGetAuthorisations_SingleRole(t *testing.T) {
	database.ConnectMockDB()

	database.Mock.ExpectQuery(`SELECT`).WillReturnRows(sqlmock.NewRows([]string{}).AddRow())
	auth, err := GetAuthorisations(userId, clubId)

	assert.Nil(t, err)
	assert.Len(t, auth, 1)
}

func TestGetDistinctClubs(t *testing.T) {
	count := 5
	auth := make([]models.Authorisation, count)
	for i := 0; i < count; i++ {
		id, _ := uuid.NewUUID()
		auth[i] = models.Authorisation{
			ClubID: id,
		}
	}

	uniqueIds := getDistinctClubs(auth)
	assert.Len(t, uniqueIds, count)
}

func TestGetDistinctClubs_WithDuplicates(t *testing.T) {
	count := 5
	auth := make([]models.Authorisation, count*2)
	for i := 0; i < count; i++ {
		id, _ := uuid.NewUUID()
		auth[i] = models.Authorisation{ClubID: id}
	}

	for i := 0; i < count; i++ {
		id, _ := uuid.Parse(auth[i].ClubID.String())
		auth[i+count] = models.Authorisation{ClubID: id}
	}

	uniqueIds := getDistinctClubs(auth)
	assert.Len(t, uniqueIds, count)
}

func TestGetDistinctClubs_WithRolesNone(t *testing.T) {
	count := 5
	auth := make([]models.Authorisation, count)
	for i := 0; i < count; i++ {
		id, _ := uuid.NewUUID()
		auth[i] = models.Authorisation{
			ClubID:     id,
			AccessRole: models.AccessRole{I18nName: RoleCoach.String()},
		}
	}

	uniqueIds := getDistinctClubs(auth, RoleHeadCoach)
	assert.Len(t, uniqueIds, 0)
}

func TestGetDistinctClubs_WithRoles(t *testing.T) {
	count := 5
	auth := make([]models.Authorisation, count*2)
	for i := 0; i < count; i++ {
		id, _ := uuid.NewUUID()
		auth[i] = models.Authorisation{
			ClubID:     id,
			AccessRole: models.AccessRole{I18nName: RoleCoach.String()},
		}
	}
	for i := 0; i < count; i++ {
		id, _ := uuid.NewUUID()
		auth[i+count] = models.Authorisation{
			ClubID:     id,
			AccessRole: models.AccessRole{I18nName: RoleHeadCoach.String()},
		}
	}

	uniqueIds := getDistinctClubs(auth, RoleHeadCoach)
	assert.Len(t, uniqueIds, count)
}

func TestHasRoleOnOneOfTheClubs_TargetIsSameClub(t *testing.T) {
	database.ConnectMockDB()
	clubId, _ := uuid.NewUUID()
	sourceUser := models.User{ID: "1"}
	targetUser := models.User{ID: "2", ClubID: clubId}

	database.Mock.ExpectQuery(`SELECT`).WillReturnRows(sqlmock.NewRows([]string{"user_id"}).
		FromCSVString(fmt.Sprintf("%v", clubId)))
	database.Mock.ExpectQuery(`SELECT`).WillReturnRows(sqlmock.NewRows([]string{}).AddRow())

	assert.True(t, HasRoleOnOneOfTheClubs(&sourceUser, &targetUser))
}

func TestHasRoleOnOneOfTheClubs_SourceHasNotEnoughPermissions(t *testing.T) {
	database.ConnectMockDB()
	clubId, _ := uuid.NewUUID()
	sourceUser := models.User{ID: "1"}
	targetUser := models.User{ID: "2", ClubID: clubId}

	database.Mock.ExpectQuery(`SELECT .* FROM "authorisations"`).WillReturnRows(sqlmock.NewRows([]string{"club_id", "access_role_id"}).
		FromCSVString(fmt.Sprintf("%v,%v", clubId, 1)))
	database.Mock.ExpectQuery(`SELECT .* FROM "access_roles"`).WillReturnRows(sqlmock.NewRows([]string{"id", "i18n_name"}).
		FromCSVString(fmt.Sprintf("%v,%v", 1, RoleCoach)))
	database.Mock.ExpectQuery(`SELECT`).WillReturnRows(sqlmock.NewRows([]string{}).AddRow())

	assert.False(t, HasRoleOnOneOfTheClubs(&sourceUser, &targetUser, RoleHeadCoach))
}

func TestHasRoleOnOneOfTheClubs_SourceHasEnoughPermissions(t *testing.T) {
	database.ConnectMockDB()
	clubId, _ := uuid.NewUUID()
	sourceUser := models.User{ID: "1"}
	targetUser := models.User{ID: "2", ClubID: clubId}

	database.Mock.ExpectQuery(`SELECT .* FROM "authorisations"`).WillReturnRows(sqlmock.NewRows([]string{"club_id", "access_role_id"}).
		FromCSVString(fmt.Sprintf("%v,%v", clubId, 1)))
	database.Mock.ExpectQuery(`SELECT .* FROM "access_roles"`).WillReturnRows(sqlmock.NewRows([]string{"id", "i18n_name"}).
		FromCSVString(fmt.Sprintf("%v,%v", 1, RoleHeadCoach)))
	database.Mock.ExpectQuery(`SELECT`).WillReturnRows(sqlmock.NewRows([]string{}).AddRow())

	assert.True(t, HasRoleOnOneOfTheClubs(&sourceUser, &targetUser, RoleHeadCoach))
}

func TestHasRoleOnOneOfTheClubs_SourceHasEnoughPermissions_SharedAthlete(t *testing.T) {
	database.ConnectMockDB()
	clubId, _ := uuid.NewUUID()
	sourceUser := models.User{ID: "1"}
	targetUser := models.User{ID: "2"}

	database.Mock.ExpectQuery(`SELECT .* FROM "authorisations"`).WillReturnRows(sqlmock.NewRows([]string{"club_id", "access_role_id"}).
		FromCSVString(fmt.Sprintf("%v,%v", clubId, 1)))
	database.Mock.ExpectQuery(`SELECT .* FROM "access_roles"`).WillReturnRows(sqlmock.NewRows([]string{"id", "i18n_name"}).
		FromCSVString(fmt.Sprintf("%v,%v", 1, RoleHeadCoach)))
	database.Mock.ExpectQuery(`SELECT .* FROM "authorisations"`).WillReturnRows(sqlmock.NewRows([]string{"club_id", "access_role_id"}).
		FromCSVString(fmt.Sprintf("%v,%v", clubId, 2)))
	database.Mock.ExpectQuery(`SELECT .* FROM "access_roles"`).WillReturnRows(sqlmock.NewRows([]string{"id", "i18n_name"}).
		FromCSVString(fmt.Sprintf("%v,%v", 2, RoleAthlete)))

	assert.True(t, HasRoleOnOneOfTheClubs(&sourceUser, &targetUser, RoleHeadCoach))
}
