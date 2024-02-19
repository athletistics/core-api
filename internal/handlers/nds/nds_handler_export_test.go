// Author  : Nicolas Baeriswyl <nicolasfabrice.baeriswyl@edu.hefr.ch>
// Created : 07.07.2023

package ndsHandler

import (
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/athletistics/core-api/database"
	"github.com/athletistics/core-api/internal/models"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestFindNdsIds(t *testing.T) {
	users := []models.User{
		{ID: "1"},
		{ID: "2"},
		{ID: "3"},
	}

	database.ConnectMockDB()
	database.Mock.ExpectQuery(`SELECT .* FROM "external_ids"`).WillReturnRows(sqlmock.NewRows([]string{"internal_id", "external_id"}).
		AddRow("1", "a").
		AddRow("2", "b"))

	ids, err := findNdsIds(users)
	assert.Nil(t, err)

	id, ok := ids[users[0].ID]
	assert.True(t, ok)
	assert.Equal(t, "a", id)

	id, ok = ids[users[1].ID]
	assert.True(t, ok)
	assert.Equal(t, "b", id)

	_, ok = ids[users[2].ID]
	assert.False(t, ok)
}
