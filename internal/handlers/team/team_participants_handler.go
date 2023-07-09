// Brief   :
// Author  : Nicolas Baeriswyl <nicolasfabrice.baeriswyl@edu.hefr.ch>
// Created : 20.06.2023

package teamHandler

import (
	"errors"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"gitlab.forge.hefr.ch/tb-athletistics/backend/database"
	"gitlab.forge.hefr.ch/tb-athletistics/backend/helpers"
	"gitlab.forge.hefr.ch/tb-athletistics/backend/internal/models"
	res "gitlab.forge.hefr.ch/tb-athletistics/backend/internal/responses"
	"gorm.io/gorm"
	"net/http"
	"strconv"
)

// GetParticipants lists all participants of the specified team.
//
//	@Summary	List team members
//	@Tags		Team
//	@Produce	json
//	@Param		team_id	path		string	true	"Team identifier"
//	@Success	200		{object}	teamHandler.GetParticipants.response
//	@Failure	401		{object}	res.Error
//	@Failure	403		{object}	res.Error	"Not enough permissions"
//	@Security	ApiKeyAuth
//	@Router		/teams/{team_id}/participants [get]
func GetParticipants(c *fiber.Ctx) error {
	type response struct {
		Athletes []models.UserShort `json:"athletes"`
		Coaches  []models.UserShort `json:"coaches"`
	}

	ok, team, err := canAccessTeam(c, []helpers.Role{helpers.RoleCoach, helpers.RoleHeadCoach})
	if !ok {
		return err
	}

	db := database.DB
	var athletes []models.UserShort
	err = db.Model(&models.User{}).
		Joins("INNER JOIN participants p ON p.user_id = users.id AND p.team_id = ? AND p.is_coach = FALSE", team.ID).
		Find(&athletes).Error
	if err != nil {
		helpers.LogError(err.Error(), helpers.GetUserId(c))
		return c.Status(http.StatusInternalServerError).JSON(res.ErrDB.Err(err))
	}

	var coaches []models.UserShort
	err = db.Model(&models.User{}).
		Joins("INNER JOIN participants p ON p.user_id = users.id AND p.team_id = ? AND p.is_coach = TRUE", team.ID).
		Find(&coaches).Error
	if err != nil {
		helpers.LogError(err.Error(), helpers.GetUserId(c))
		return c.Status(http.StatusInternalServerError).JSON(res.ErrDB.Err(err))
	}

	return c.JSON(response{athletes, coaches})
}

// AddParticipants adds a list of participants to a team.
//
//	@Summary	Add participants to a team
//	@Tags		Team
//	@Consumes	json
//	@Param		team_id	path	string						true	"Team identifier"
//	@Param		request	body	[]models.ParticipantInOut	true	"Users to add"
//	@Success	200
//	@Failure	400	{object}	res.Error	"Bad user IDs, or user already added to team"
//	@Failure	401	{object}	res.Error
//	@Failure	403	{object}	res.Error	"Not enough permissions"
//	@Security	ApiKeyAuth
//	@Router		/teams/{team_id}/participants [post]
func AddParticipants(c *fiber.Ctx) error {
	db := database.DB
	ok, team, err := canAccessTeam(c, []helpers.Role{helpers.RoleHeadCoach})
	if !ok {
		return err
	}

	// extract all IDs, and check if they are in a valid format
	participantsReq, userIds, ok, err := getUserIdsFromBody(c)
	if !ok {
		return err
	}

	// get requested users
	var users []models.User
	err = db.Find(&users, userIds).Error
	if err != nil {
		helpers.LogError(err.Error(), helpers.GetUserId(c))
		return c.Status(http.StatusInternalServerError).JSON(res.ErrDB.Err(err))
	}

	// check if all users are member of the head coach's club
	for i := 0; i < len(users); i++ {
		if users[i].ClubID != team.ClubID {
			helpers.LogSecurity(fmt.Sprintf("tried modifying team members %v", team.ID), helpers.GetUserId(c))
			return c.Status(http.StatusForbidden).JSON(res.ErrForbidden.Err(err))
		}
	}

	// create necessary participants objects
	participants := make([]models.Participant, len(users))
	for i, p := range participantsReq {
		participants[i] = models.Participant{UserID: p.UserID, TeamID: team.ID, IsCoach: p.IsCoach}
	}

	// write objects to database
	err = db.Create(&participants).Error
	if errors.Is(err, gorm.ErrDuplicatedKey) {
		helpers.LogInfo(fmt.Sprintf("bad request: %v", err), helpers.GetUserId(c))
		return c.Status(http.StatusBadRequest).JSON(res.ErrParticipantsExist)
	} else if err != nil {
		helpers.LogError(err.Error(), helpers.GetUserId(c))
		return c.Status(http.StatusInternalServerError).JSON(res.ErrDB.Err(err))
	}

	return c.SendStatus(http.StatusOK)
}

// RemoveParticipants removes a set of users from a specific team.
//
//	@Summary	Remove participants
//	@Tags		Team
//	@Consumes	json
//	@Produce	json
//	@Param		team_id	path	string		true	"Team identifier"
//	@Param		request	body	[]string	true	"User IDs"
//	@Success	200
//	@Failure	400	{object}	res.Error	"Bad user IDs, or user already added to team"
//	@Failure	401	{object}	res.Error
//	@Failure	403	{object}	res.Error	"Not enough permissions"
//	@Security	ApiKeyAuth
//	@Router		/teams/{team_id}/participants [delete]
func RemoveParticipants(c *fiber.Ctx) error {
	ok, team, err := canAccessTeam(c, []helpers.Role{helpers.RoleHeadCoach})
	if !ok {
		return err
	}

	// extract all IDs, and check if they are valid UUIDs
	_, userIds, ok, err := getUserIdsFromBody(c)
	if !ok {
		return err
	}

	participants := make([]models.Participant, len(userIds))
	for i := 0; i < len(userIds); i++ {
		participants[i] = models.Participant{TeamID: team.ID, UserID: userIds[i]}
	}

	err = database.DB.Unscoped().Delete(&participants).Error
	if err != nil {
		helpers.LogError(err.Error(), helpers.GetUserId(c))
		return c.Status(http.StatusInternalServerError).JSON(res.ErrDB.Err(err))
	}

	return c.SendStatus(http.StatusOK)
}

// getUserIdsFromBody reads a request body and extracts user IDs. The body must
// be an array of strings, and user IDs must be numerical positive integers.
//
// In case of any error, the function directly creates the response to the HTTP request.
// Therefore, if the second result is "false", return the third parameter and exit the
// calling function.
func getUserIdsFromBody(c *fiber.Ctx) ([]models.ParticipantInOut, []string, bool, error) {
	var err error

	// parse the body
	userIds := new([]models.ParticipantInOut)
	if err = c.BodyParser(userIds); err != nil {
		helpers.LogInfo(fmt.Sprintf("bad request: %v", err), helpers.GetUserId(c))
		return nil, nil, false, c.Status(http.StatusBadRequest).JSON(res.ErrBadRequest.Err(err))
	}

	// verify that IDs are numerical (base 10, 64 bit integer)
	userIdsStr := make([]string, len(*userIds))
	for i, id := range *userIds {
		if _, err = strconv.ParseInt(id.UserID, 10, 64); err != nil {
			helpers.LogInfo(fmt.Sprintf("bad request: %v", err), helpers.GetUserId(c))
			return nil, nil, false, c.Status(http.StatusBadRequest).JSON(res.ErrBadRequest.Err(err))
		}
		userIdsStr[i] = id.UserID
	}

	// no IDs, call is considered to be finished as nothing can be done.
	if len(*userIds) == 0 {
		return nil, nil, false, c.SendStatus(http.StatusOK)
	}

	return *userIds, userIdsStr, true, nil
}
