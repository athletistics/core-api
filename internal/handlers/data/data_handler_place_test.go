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
