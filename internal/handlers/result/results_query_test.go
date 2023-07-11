// Brief   :
// Author  : Nicolas Baeriswyl <nicolasfabrice.baeriswyl@edu.hefr.ch>
// Created : 26.06.2023

package resultHandler

import (
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"gitlab.forge.hefr.ch/tb-athletistics/backend/database"
	"gitlab.forge.hefr.ch/tb-athletistics/backend/internal/models"
	"testing"
)

var user = &models.User{ID: "123"}

func TestFindResults(t *testing.T) {
	id, _ := uuid.NewUUID()
	user.ClubID = id

	database.ConnectMockDB()
	database.Mock.ExpectQuery(`SELECT .* FROM "results".*users\.club_id = \$1.*`).WillReturnRows(sqlmock.NewRows([]string{}))

	r := &ResultQuery{}

	var output = &models.Result{}
	err := FindResults(r, output, user)

	assert.Nil(t, err)
}

func TestFindResults_User(t *testing.T) {
	id := "id"

	database.ConnectMockDB()
	database.Mock.ExpectQuery(`SELECT .* "results"\."user_id" = \$3.*`).WillReturnRows(sqlmock.NewRows([]string{}))

	r := &ResultQuery{UserID: &id}

	var output = &models.Result{}
	err := FindResults(r, output, user)

	assert.Nil(t, err)
}

func TestFindResults_Event(t *testing.T) {
	id, _ := uuid.NewUUID()
	user.ClubID = id

	database.ConnectMockDB()
	database.Mock.ExpectQuery(`SELECT .* FROM "results".*"results"\."event_id" = \$3.*`).WillReturnRows(sqlmock.NewRows([]string{}))

	r := &ResultQuery{EventID: &id}

	var output = &models.Result{}
	err := FindResults(r, output, user)

	assert.Nil(t, err)
}

func TestFindResults_EventType(t *testing.T) {
	id := 1

	database.ConnectMockDB()
	database.Mock.ExpectQuery(`SELECT .* FROM "results".*inner join events.*events\.event_type_id = \$3.*`).WillReturnRows(sqlmock.NewRows([]string{}))

	r := &ResultQuery{EventTypeID: &id}

	var output = &models.Result{}
	err := FindResults(r, output, user)

	assert.Nil(t, err)
}

func TestFindResults_Discipline(t *testing.T) {
	id, _ := uuid.NewUUID()
	user.ClubID = id

	database.ConnectMockDB()
	database.Mock.ExpectQuery(`SELECT.*FROM "results".*"results"\."discipline_id" = \$3.*`).WillReturnRows(sqlmock.NewRows([]string{}))

	r := &ResultQuery{DisciplineID: &id}

	var output = &models.Result{}
	err := FindResults(r, output, user)

	assert.Nil(t, err)
}

func TestFindResults_Best(t *testing.T) {
	id, _ := uuid.NewUUID()
	user.ClubID = id

	database.ConnectMockDB()
	database.Mock.ExpectQuery(`SELECT.*UNION.*`).WillReturnRows(sqlmock.NewRows([]string{}))

	r := &ResultQuery{BestOnly: true}

	var output = &models.Result{}
	err := FindResults(r, output, user)

	assert.Nil(t, err)
}
