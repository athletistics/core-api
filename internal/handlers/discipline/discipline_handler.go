// Brief   :
// Author  : Nicolas Baeriswyl <nicolasfabrice.baeriswyl@edu.hefr.ch>
// Created : 21.06.2023

package disciplineHandler

import (
	"github.com/gofiber/fiber/v2"
	"gitlab.forge.hefr.ch/tb-athletistics/backend/helpers"
)

func GetDisciplines(c *fiber.Ctx) error {
	return nil
}

func CreateDisciplines(c *fiber.Ctx) error {
	// todo: create function that stores the entire user in "locals"
	if !helpers.HasRoleOnParentClub(c, []helpers.Role{helpers.RoleCoach}) {
		//return
	}
	return nil
}
