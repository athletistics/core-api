// Brief   :
// Author  : Nicolas Baeriswyl <nicolasfabrice.baeriswyl@edu.hefr.ch>
// Created : 27.06.2023

package dataHandler

import (
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"gitlab.forge.hefr.ch/tb-athletistics/backend/database"
	"testing"
)

const countryId = 1
const stateId = 2

func TestSelectCountries(t *testing.T) {
	database.ConnectMockDB()
	database.Mock.ExpectQuery(`SELECT \* FROM "countries"`).WillReturnRows(sqlmock.NewRows([]string{}))
	database.Mock.ExpectQuery(`SELECT \* FROM "countries"`).WillReturnRows(sqlmock.NewRows([]string{}).AddRow())

	c, err := SelectCountries()
	assert.Nil(t, err)
	assert.Empty(t, c)

	c, err = SelectCountries()
	assert.Nil(t, err)
	assert.Len(t, c, 1)
}

func TestSelectStatesByCountry(t *testing.T) {
	database.ConnectMockDB()
	database.Mock.ExpectQuery(`SELECT \* FROM "states" WHERE "states"."country_id" = \$1`).WillReturnRows(sqlmock.NewRows([]string{}))
	database.Mock.ExpectQuery(`SELECT \* FROM "states" WHERE "states"."country_id" = \$1`).WillReturnRows(sqlmock.NewRows([]string{}).AddRow())

	c, err := SelectStatesByCountry(countryId)
	assert.Nil(t, err)
	assert.Empty(t, c)

	c, err = SelectStatesByCountry(countryId)
	assert.Nil(t, err)
	assert.Len(t, c, 1)
}

func TestSearchForPlaces(t *testing.T) {
	database.ConnectMockDB()
	database.Mock.ExpectQuery(`SELECT.*FROM "places" INNER JOIN.*ON.*"State"."country_id" = \$2`).WillReturnRows(sqlmock.NewRows([]string{}))
	database.Mock.ExpectQuery(`SELECT.*FROM "places" INNER JOIN.*WHERE.*"post_code" LIKE \$3`).WillReturnRows(sqlmock.NewRows([]string{}).AddRow())

	c, err := SearchForPlaces("", countryId, stateId)
	assert.Nil(t, err)
	assert.Empty(t, c)

	c, err = SearchForPlaces("q", countryId, stateId)
	assert.Nil(t, err)
	assert.Len(t, c, 1)
}
