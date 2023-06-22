// Brief   :
// Author  : Nicolas Baeriswyl <nicolasfabrice.baeriswyl@edu.hefr.ch>
// Created : 20.06.2023

package clubHandler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gitlab.forge.hefr.ch/tb-athletistics/backend/database"
	"gitlab.forge.hefr.ch/tb-athletistics/backend/helpers"
	"gitlab.forge.hefr.ch/tb-athletistics/backend/internal/models"
	res "gitlab.forge.hefr.ch/tb-athletistics/backend/internal/responses"
	"net/http"
)

const (
	ParamClubUuid = "uuid"
)

// GetClub returns details about a club.
//
//	@Summary	Get club details
//	@Tags		Club
//	@Produce	json
//	@Param		club_id	path		string	true	"Club identifier"
//	@Success	200		{object}	models.Club
//	@Failure	401		{object}	res.Error
//	@Failure	403		{object}	res.Error	"Not enough permissions"
//	@Security	ApiKeyAuth
//	@Router		/club/{club_id} [get]
func GetClub(c *fiber.Ctx) error {
	clubId, err := uuid.Parse(c.Params(ParamClubUuid))
	if err != nil || !helpers.HasAnyRoleOnClub(c, clubId) {
		return c.Status(http.StatusForbidden).JSON(res.ErrForbidden)
	}

	var club models.Club
	err = database.DB.Where(&models.Club{ID: clubId}).First(&club).Error
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(res.ErrDB.Err(err))
	}

	return c.JSON(club)
}

// GetTeams lists all teams of the corresponding club.
//
//	@Summary	List teams
//	@Tags		Club
//	@Produce	json
//	@Param		club_id	path		string	true	"Club identifier"
//	@Success	200		{object}	[]models.Team
//	@Failure	401		{object}	res.Error
//	@Failure	403		{object}	res.Error	"Not enough permissions"
//	@Security	ApiKeyAuth
//	@Router		/club/{club_id}/teams [get]
func GetTeams(c *fiber.Ctx) error {
	clubId, err := uuid.Parse(c.Params(ParamClubUuid))
	if err != nil || !helpers.HasRoleOnClub(c, clubId, helpers.RoleCoach, helpers.RoleHeadCoach) {
		return c.Status(http.StatusForbidden).JSON(res.ErrForbidden)
	}

	var teams []models.Team
	err = database.DB.Where(&models.Team{ClubID: clubId}).Find(&teams).Error
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(res.ErrDB.Err(err))
	}

	return c.JSON(teams)
}

// CreateTeam creates a new team for the corresponding club.
//
//	@Summary	Create a new team
//	@Tags		Club
//	@Consumes	json
//	@Produce	json
//	@Param		club_id	path		string							true	"Club identifier"
//	@Param		request	body		clubHandler.CreateTeam.ReqBody	true	"Team data"
//	@Success	200		{object}	res.UuidCreated
//	@Failure	401		{object}	res.Error
//	@Failure	403		{object}	res.Error	"Not enough permissions"
//	@Security	ApiKeyAuth
//	@Router		/club/{club_id}/teams [post]
func CreateTeam(c *fiber.Ctx) error {
	type ReqBody struct {
		Name string `json:"name" validate:"required" format:"uuid"`
	}
	body := new(ReqBody)
	if err := c.BodyParser(body); err != nil {
		return c.Status(http.StatusBadRequest).JSON(res.ErrBadRequest.Err(err))
	}

	clubId, err := uuid.Parse(c.Params(ParamClubUuid))
	if err != nil || !helpers.HasRoleOnClub(c, clubId, helpers.RoleHeadCoach) {
		return c.Status(http.StatusForbidden).JSON(res.ErrForbidden)
	}

	team := models.Team{ClubID: clubId, Name: body.Name}
	err = database.DB.Create(&team).Error
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(res.ErrDB.Err(err))
	}

	return c.JSON(res.UuidCreated{Id: team.ID})
}
