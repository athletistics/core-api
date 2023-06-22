// Brief   :
// Author  : Nicolas Baeriswyl <nicolasfabrice.baeriswyl@edu.hefr.ch>
// Created : 22.06.2023

package resultHandler

import (
	"github.com/gofiber/fiber/v2"
	"gitlab.forge.hefr.ch/tb-athletistics/backend/helpers"
	"gitlab.forge.hefr.ch/tb-athletistics/backend/internal/models"
	res "gitlab.forge.hefr.ch/tb-athletistics/backend/internal/responses"
	"net/http"
)

func GetResults(c *fiber.Ctx) error {
	if !helpers.HasRoleOnParentClub(c, helpers.RoleCoach) {
		return c.Status(http.StatusForbidden).JSON(res.ErrForbidden)
	}

	q := new(ResultQuery)
	if err := c.QueryParser(q); err != nil {
		return c.Status(http.StatusBadRequest).JSON(res.ErrBadRequest.Err(err))
	}

	return RunQuery(c, q)
}

func CreateResult(c *fiber.Ctx) error {
	if !helpers.HasRoleOnParentClub(c, helpers.RoleCoach) {
		return c.Status(http.StatusForbidden).JSON(res.ErrForbidden)
	}

	body := new(models.ResultInput)
	if err := c.BodyParser(body); err != nil {
		return c.Status(http.StatusBadRequest).JSON(res.ErrBadRequest.Err(err))
	}

	return RunCreateQuery(c, body)
}

func GetBestResults(c *fiber.Ctx) error {
	if !helpers.HasRoleOnParentClub(c, helpers.RoleCoach) {
		return c.Status(http.StatusForbidden).JSON(res.ErrForbidden)
	}

	q := new(ResultQuery)
	if err := c.QueryParser(q); err != nil {
		return c.Status(http.StatusBadRequest).JSON(res.ErrBadRequest.Err(err))
	}

	q.BestOnly = true

	return RunQuery(c, q)
}
