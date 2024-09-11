// Author  : Nicolas Baeriswyl <nicolasfabrice.baeriswyl@edu.hefr.ch>
// Created : 16.06.2024

// nolint:all
package eventHandler

import (
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/athletistics/core-api/database"
	h "github.com/athletistics/core-api/helpers"
	"github.com/athletistics/core-api/internal/models"
	res "github.com/athletistics/core-api/internal/responses"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestPermissionCheck_NoUser(t *testing.T) {
	ctx := h.Context{}
	_, err := permissionCheck(&ctx, uuid.New().String())

	assert.NotNil(t, err)
	assert.Equal(t, res.ErrForbidden.Error, err.Error)
}

func TestPermissionCheck_InvalidId(t *testing.T) {
	ctx := h.CreateContext(models.User{ID: "123"}, nil, h.TwoFAState{})
	_, err := permissionCheck(ctx, "")

	assert.NotNil(t, err)
	assert.Equal(t, res.ErrBadRequest.Error, err.Error)
}

func TestPermissionCheck_DoesNotHaveRole(t *testing.T) {
	clubUuid := uuid.New()

	database.ConnectMockDB()
	database.Mock.ExpectQuery(`SELECT \* FROM "events".*`).WillReturnRows(
		sqlmock.NewRows([]string{"id", "club_id"}).AddRow(uuid.New().String(), clubUuid))

	ctx := h.CreateContext(models.User{ID: "123"}, nil, h.TwoFAState{})
	_, err := permissionCheck(ctx, clubUuid.String())

	assert.NotNil(t, err)
	assert.Equal(t, res.ErrForbidden.Error, err.Error)

	database.Mock.ExpectQuery(`SELECT \* FROM "events".*`).WillReturnRows(
		sqlmock.NewRows([]string{"id", "club_id"}).AddRow(uuid.New().String(), clubUuid))

	authByClub := make(map[uuid.UUID]map[h.Role]bool)
	authByClub[clubUuid] = make(map[h.Role]bool)
	authByClub[clubUuid][h.RoleAthlete] = true
	ctx = h.CreateContext(models.User{ID: "123"}, authByClub, h.TwoFAState{})
	_, err = permissionCheck(ctx, clubUuid.String())

	assert.NotNil(t, err)
	assert.Equal(t, res.ErrForbidden.Error, err.Error)
}

func TestPermissionCheck_OK(t *testing.T) {
	ctx, _, clubId, _ := getAuthorizedContext()

	_, err := permissionCheck(ctx, clubId.String())

	assert.Nil(t, err)
}

func TestUpdateEvent_Unauthorized(t *testing.T) {
	ctx, _, clubId, eventId := getAuthorizedContext()

	authByClub := make(map[uuid.UUID]map[h.Role]bool)
	authByClub[clubId] = make(map[h.Role]bool)
	authByClub[clubId][h.RoleAthlete] = true
	ctx.SetAuthMap(authByClub)

	err := updateEvent(ctx, eventId.String(), &models.EventInput{}, nil)

	assert.NotNil(t, err)
	assert.Equal(t, res.ErrForbidden.Error, err.Error)
}

func TestUpdateEvent_OK(t *testing.T) {
	ctx, _, _, eventId := getAuthorizedContext()

	database.Mock.ExpectBegin()
	database.Mock.ExpectExec(`UPDATE "events".*`).WillReturnResult(sqlmock.NewResult(1, 1))
	database.Mock.ExpectCommit()
	err := updateEvent(ctx, eventId.String(), &models.EventInput{}, nil)

	assert.Nil(t, err)
}

func TestDeleteEvent_Unauthorized(t *testing.T) {
	ctx, _, clubId, eventId := getAuthorizedContext()

	authByClub := make(map[uuid.UUID]map[h.Role]bool)
	authByClub[clubId] = make(map[h.Role]bool)
	authByClub[clubId][h.RoleAthlete] = true
	ctx.SetAuthMap(authByClub)

	err := deleteEvent(ctx, eventId.String(), false)

	assert.NotNil(t, err)
	assert.Equal(t, res.ErrForbidden.Error, err.Error)
}

func TestDeleteEvent_HasChildren(t *testing.T) {
	ctx, _, _, eventId := getAuthorizedContext()

	database.Mock.ExpectBegin()
	database.Mock.ExpectQuery(`SELECT count\(\*\) FROM "results".* WHERE event_id = \$1.*`).
		WillReturnRows(sqlmock.NewRows([]string{"count(*)"}).AddRow(8))
	database.Mock.ExpectRollback()
	err := deleteEvent(ctx, eventId.String(), false)

	assert.NotNil(t, err)
	assert.Equal(t, res.ErrBadRequest.Error, err.Error)
	assert.Equal(t, res.ErrEventHasResults.ErrorSubCode, err.ErrorSubCode)
}

func TestDeleteEvent_Force(t *testing.T) {
	ctx, _, _, eventId := getAuthorizedContext()

	database.Mock.ExpectQuery(`SELECT count\(\*\) FROM "results".* WHERE "results"."event_id" = \$1.*`).
		WillReturnRows(sqlmock.NewRows([]string{"count(*)"}).AddRow(8))
	database.Mock.ExpectExec(`DELETE FROM "results".* WHERE "results"."event_id" = \$1.*`).
		WillReturnResult(sqlmock.NewResult(1, 1))
	err := deleteEvent(ctx, eventId.String(), true)

	assert.NotNil(t, err)
}

func TestDeleteEvent_OK(t *testing.T) {
	ctx, _, _, eventId := getAuthorizedContext()

	database.Mock.ExpectQuery(`SELECT count\(\*\) FROM "results".* WHERE "results"."event_id" = \$1.*`).
		WillReturnRows(sqlmock.NewRows([]string{"count(*)"}).AddRow(0))
	database.Mock.ExpectExec(`DELETE FROM "results".* WHERE "results"."event_id" = \$1.*`).
		WillReturnResult(sqlmock.NewResult(1, 1))
	err := deleteEvent(ctx, eventId.String(), false)

	assert.NotNil(t, err)
}

func getAuthorizedContext() (*h.Context, string, uuid.UUID, uuid.UUID) {
	userId := "123"
	clubId := uuid.New()
	eventId := uuid.New()

	database.ConnectMockDB()
	database.Mock.ExpectQuery(`SELECT \* FROM "events".*`).WillReturnRows(
		sqlmock.NewRows([]string{"id", "club_id"}).AddRow(uuid.New().String(), clubId))

	authByClub := make(map[uuid.UUID]map[h.Role]bool)
	authByClub[clubId] = make(map[h.Role]bool)
	authByClub[clubId][h.RoleCoach] = true
	ctx := h.CreateContext(models.User{ID: userId}, authByClub, h.TwoFAState{})

	return ctx, userId, clubId, eventId
}
