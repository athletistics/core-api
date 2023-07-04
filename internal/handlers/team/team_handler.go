// Brief   :
// Author  : Nicolas Baeriswyl <nicolasfabrice.baeriswyl@edu.hefr.ch>
// Created : 20.06.2023

package teamHandler

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gitlab.forge.hefr.ch/tb-athletistics/backend/database"
	"gitlab.forge.hefr.ch/tb-athletistics/backend/helpers"
	"gitlab.forge.hefr.ch/tb-athletistics/backend/internal/models"
	res "gitlab.forge.hefr.ch/tb-athletistics/backend/internal/responses"
	"net/http"
)

const (
	ParamTeamId = "team"
)

// GetTeam returns detailed information about a specific team.
//
//	@Summary	Team details
//	@Tags		Team
//	@Produce	json
//	@Param		team_id	path		string	true	"Team identifier"
//	@Success	200		{object}	models.Team
//	@Failure	401		{object}	res.Error
//	@Failure	403		{object}	res.Error	"Not enough permissions"
//	@Security	ApiKeyAuth
//	@Router		/teams/{team_id} [get]
func GetTeam(c *fiber.Ctx) error {
	ok, team, err := canAccessTeam(c, []helpers.Role{helpers.RoleCoach, helpers.RoleHeadCoach})
	if !ok {
		return err
	}

	return c.JSON(team)
}

func canAccessTeam(c *fiber.Ctx, roles []helpers.Role) (bool, *models.Team, error) {
	teamId, err := uuid.Parse(c.Params(ParamTeamId))
	if err != nil {
		helpers.LogInfo(fmt.Sprintf("bad request: %v", err), helpers.GetUserId(c))
		return false, nil, c.Status(http.StatusBadRequest).JSON(res.ErrBadRequest.Err(err))
	}

	var team models.Team
	err = database.DB.Where(&models.Team{ID: teamId}).First(&team).Error
	if err != nil {
		helpers.LogError(err.Error(), helpers.GetUserId(c))
		return false, nil, c.Status(http.StatusInternalServerError).JSON(res.ErrDB.Err(err))
	}

	if !helpers.HasRoleOnClub(c, team.ClubID, roles...) {
		helpers.LogSecurity(fmt.Sprintf("tried accessing team %v", teamId), helpers.GetUserId(c))
		return false, nil, c.Status(http.StatusForbidden).JSON(res.ErrForbidden)
	}

	return true, &team, nil
}
