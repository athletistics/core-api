// Brief   :
// Author  : Nicolas Baeriswyl <nicolasfabrice.baeriswyl@edu.hefr.ch>
// Created : 26.06.2023

package helpers

import (
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"gitlab.forge.hefr.ch/tb-athletistics/backend/database"
	"testing"
)

var (
	userId    = "123"
	clubId, _ = uuid.NewUUID()
)

func TestGetAuthorisations_NoRoles(t *testing.T) {
	database.ConnectMockDB()

	database.Mock.ExpectQuery(`SELECT`).WillReturnRows(sqlmock.NewRows([]string{}))
	auth, err := getAuthorisations(userId, clubId)

	assert.Nil(t, err)
	assert.Empty(t, auth)
}

func TestGetAuthorisations_SingleRole(t *testing.T) {
	database.ConnectMockDB()

	database.Mock.ExpectQuery(`SELECT`).WillReturnRows(sqlmock.NewRows([]string{}).AddRow())
	auth, err := getAuthorisations(userId, clubId)

	assert.Nil(t, err)
	assert.Len(t, auth, 1)
}
