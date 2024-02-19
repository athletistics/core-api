// Brief   :
// Author  : Nicolas Baeriswyl <nicolasfabrice.baeriswyl@edu.hefr.ch>
// Created : 27.06.2023

package dataHandler

import (
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/athletistics/core-api/database"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSelectIdTypes(t *testing.T) {
	database.ConnectMockDB()
	database.Mock.ExpectQuery(`SELECT \* FROM "external_id_types"`).WillReturnRows(sqlmock.NewRows([]string{}).AddRow())

	c, err := SelectIdTypes()
	assert.Nil(t, err)
	assert.Len(t, c, 1)
}

func TestSelectAccessRoles(t *testing.T) {
	database.ConnectMockDB()
	database.Mock.ExpectQuery(`SELECT \* FROM "access_roles"`).WillReturnRows(sqlmock.NewRows([]string{}).AddRow())

	c, err := SelectAccessRoles()
	assert.Nil(t, err)
	assert.Len(t, c, 1)
}

func TestSelectLanguages(t *testing.T) {
	database.ConnectMockDB()
	database.Mock.ExpectQuery(`SELECT \* FROM "languages"`).WillReturnRows(sqlmock.NewRows([]string{}).AddRow())

	c, err := SelectLanguages()
	assert.Nil(t, err)
	assert.Len(t, c, 1)
}
