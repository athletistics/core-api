// Brief   :
// Author  : Nicolas Baeriswyl <nicolasfabrice.baeriswyl@edu.hefr.ch>
// Created : 25.06.2023

//nolint:bodyclose
package test

import (
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/athletistics/core-api/constants"
	"github.com/athletistics/core-api/database"
	"github.com/athletistics/core-api/router"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
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
	database.Mock.ExpectQuery(`SELECT .* FROM "users"`).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(userId))

	req := httptest.NewRequest("GET", "/api/v1/me", nil)
	resp, err := app.Test(req, -1)

	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestGetMe_NoUser(t *testing.T) {
	userId := "123"
	app := prepareFiberWithUser(userId)

	_ = sqlmock.NewRows([]string{""})
	database.Mock.ExpectQuery(`SELECT .* FROM "users"`).WillReturnError(gorm.ErrRecordNotFound)

	req := httptest.NewRequest("GET", "/api/v1/me", nil)
	resp, err := app.Test(req, -1)

	assert.Nil(t, err)
	assert.Equal(t, http.StatusNotFound, resp.StatusCode)
}
