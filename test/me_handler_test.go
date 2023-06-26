// Brief   :
// Author  : Nicolas Baeriswyl <nicolasfabrice.baeriswyl@edu.hefr.ch>
// Created : 25.06.2023

//nolint:bodyclose
package test

import (
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"gitlab.forge.hefr.ch/tb-athletistics/backend/constants"
	"gitlab.forge.hefr.ch/tb-athletistics/backend/database"
	"gitlab.forge.hefr.ch/tb-athletistics/backend/router"
	"net/http"
	"net/http/httptest"
	"testing"
)

func prepareFiberWithUser(userId string) *fiber.App {
	database.ConnectMockDB()
	app := fiber.New()

	app.Use("/", func(c *fiber.Ctx) error {
		c.Locals(constants.LocalsUserId, userId)
		return c.Next()
	})
	router.SetupRoutes(app)

	return app
}

func TestGetMe(t *testing.T) {
	userId := "123"
	app := prepareFiberWithUser(userId)

	_ = sqlmock.NewRows([]string{""})
	database.Mock.ExpectQuery(`SELECT \* FROM "users"`).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(userId))

	req := httptest.NewRequest("GET", "/api/v1/me", nil)
	resp, err := app.Test(req, -1)

	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestGetMe_NoUser(t *testing.T) {
	userId := "123"
	app := prepareFiberWithUser(userId)

	_ = sqlmock.NewRows([]string{""})
	database.Mock.ExpectQuery(`SELECT \* FROM "users"`).WillReturnRows(sqlmock.NewRows([]string{"id"}))

	req := httptest.NewRequest("GET", "/api/v1/me", nil)
	resp, err := app.Test(req, -1)

	assert.Nil(t, err)
	assert.Equal(t, http.StatusNotFound, resp.StatusCode)
}
