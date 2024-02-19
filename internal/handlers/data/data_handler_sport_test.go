// Brief   :
// Author  : Nicolas Baeriswyl <nicolasfabrice.baeriswyl@edu.hefr.ch>
// Created : 27.06.2023

package dataHandler

import (
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/athletistics/core-api/database"
	"github.com/athletistics/core-api/internal/models"
	"github.com/stretchr/testify/assert"
	"testing"
)

const sportId = 1

func TestSelectSports(t *testing.T) {
	database.ConnectMockDB()
	database.Mock.ExpectQuery(`SELECT \* FROM "sports"`).WillReturnRows(sqlmock.NewRows([]string{}).AddRow())

	c, err := SelectSports()
	assert.Nil(t, err)
	assert.Len(t, c, 1)
}

func TestSelectCategoriesBySport(t *testing.T) {
	database.ConnectMockDB()
	database.Mock.ExpectQuery(`SELECT \* FROM "categories" WHERE "categories"."sport_id" = \$1`).WillReturnRows(sqlmock.NewRows([]string{}).AddRow())

	c, err := SelectCategoriesBySport(sportId)
	assert.Nil(t, err)
	assert.Len(t, c, 1)
}

func TestSelectEventTypes(t *testing.T) {
	database.ConnectMockDB()
	database.Mock.ExpectQuery(`SELECT .* FROM "event_types".*INNER JOIN`).WillReturnRows(sqlmock.NewRows([]string{}))
	database.Mock.ExpectQuery(`SELECT .* FROM "event_types".*INNER JOIN`).WillReturnRows(sqlmock.NewRows([]string{}).AddRow())

	c, err := SelectEventTypes(true, []models.Authorisation{})
	assert.Nil(t, err)
	assert.Empty(t, c)

	c, err = SelectEventTypes(true, []models.Authorisation{{}})
	assert.Nil(t, err)
	assert.Len(t, c, 0)

	c, err = SelectEventTypes(true, []models.Authorisation{{}})
	assert.Nil(t, err)
	assert.Len(t, c, 1)
}
