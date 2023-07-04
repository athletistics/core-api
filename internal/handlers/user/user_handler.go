// Brief   :
// Author  : Nicolas Baeriswyl <nicolasfabrice.baeriswyl@edu.hefr.ch>
// Created : 20.06.2023

package userHandler

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

const (
	ParamUserId = "id"
)

// GetUsers lists all users owned by the authenticated user's club.
//
//	@Summary	List users
//	@Tags		User info
//	@Produce	json
//	@Success	200	{object}	[]models.User
//	@Failure	401	{object}	res.Error
//	@Failure	403	{object}	res.Error	"Not enough permissions"
//	@Security	ApiKeyAuth
//	@Router		/users [get]
func GetUsers(c *fiber.Ctx) error {
	db := database.DB
	user := helpers.GetUser(c)

	if !helpers.HasRoleOnClub(c, user.ClubID, helpers.RoleCoach, helpers.RoleHeadCoach) {
		helpers.LogSecurity(fmt.Sprintf("tried accessing users of club %v", user.ClubID), helpers.GetUserId(c))
		return c.Status(http.StatusForbidden).JSON(res.ErrForbidden)
	}

	var users []models.User
	err := db.Where(&models.User{ClubID: user.ClubID}).Find(&users).Error
	if err != nil {
		helpers.LogError(err.Error(), helpers.GetUserId(c))
		return c.Status(http.StatusInternalServerError).JSON(res.ErrDB.Err(err))
	}

	return c.JSON(users)
}

// GetEmergencyUsers lists all users owned by the authenticated user's club, but only shows emergency information.
//
//	@Summary	List users
//	@Tags		User info
//	@Produce	json
//	@Success	200	{object}	[]models.UserEmergency
//	@Failure	401	{object}	res.Error
//	@Failure	403	{object}	res.Error	"Not enough permissions"
//	@Security	ApiKeyAuth
//	@Router		/users/emergency [get]
func GetEmergencyUsers(c *fiber.Ctx) error {
	db := database.DB
	user := helpers.GetUser(c)

	if !helpers.HasRoleOnClub(c, user.ClubID, helpers.RoleCoach, helpers.RoleHeadCoach) {
		helpers.LogSecurity(fmt.Sprintf("tried accessing users of club %v", user.ClubID), helpers.GetUserId(c))
		return c.Status(http.StatusForbidden).JSON(res.ErrForbidden)
	}

	var users []models.UserEmergency
	err := db.Model(&models.User{}).Where(&models.User{ClubID: user.ClubID}).Find(&users).Error
	if err != nil {
		helpers.LogError(err.Error(), helpers.GetUserId(c))
		return c.Status(http.StatusInternalServerError).JSON(res.ErrDB.Err(err))
	}

	return c.JSON(users)
}

// CreateUser creates a new user on Zitadel and the core database.
//
//	@Summary	Create a new user
//	@Tags		User info
//	@Consume	json
//	@Produce	json
//	@Param		request	body		models.User	true	"User data"
//	@Success	200		{object}	res.IdCreated
//	@Failure	400		{object}	res.Error
//	@Failure	401		{object}	res.Error
//	@Failure	403		{object}	res.Error	"Not enough permissions"
//	@Security	ApiKeyAuth
//	@Router		/users [post]
func CreateUser(c *fiber.Ctx) error {
	if !helpers.HasRoleOnParentClub(c, helpers.RoleHeadCoach) {
		helpers.LogSecurity("tried creating user (missing role)", helpers.GetUserId(c))
		return c.Status(http.StatusForbidden).JSON(res.ErrForbidden)
	}

	// check that requesting user is head coach
	db := database.DB
	user := helpers.GetUser(c)

	if !helpers.HasRoleOnClub(c, user.ClubID, helpers.RoleHeadCoach) {
		helpers.LogSecurity("tried creating user (missing required role)", helpers.GetUserId(c))
		return c.Status(http.StatusForbidden).JSON(res.ErrForbidden)
	}

	// extract values from body
	newUser := new(models.User)
	if err := c.BodyParser(newUser); err != nil {
		helpers.LogInfo(fmt.Sprintf("bad request: %v", err), helpers.GetUserId(c))
		return c.Status(http.StatusBadRequest).JSON(res.ErrBadRequest.Err(err))
	}

	// create the user on zitadel
	id, err := createUser(newUser)
	if err != nil {
		helpers.LogError(err.Error(), helpers.GetUserId(c))
		return c.Status(http.StatusInternalServerError).JSON(res.ErrBadRequest.Err(err))
	}

	// create the user locally
	newUser.ID = id
	newUser.ClubID = user.ClubID
	err = db.Create(newUser).Error
	if err != nil {
		helpers.LogError(err.Error(), helpers.GetUserId(c))
		return c.Status(http.StatusInternalServerError).JSON(res.ErrDB.Err(err))
	}

	return c.JSON(res.IdCreated{Id: newUser.ID})
}

// GetUser returns details of a specific user.
//
//	@Summary	Get a user's details
//	@Tags		User info
//	@Produce	json
//	@Success	200	{object}	models.User
//	@Failure	401	{object}	res.Error
//	@Failure	403	{object}	res.Error	"Not enough permissions"
//	@Security	ApiKeyAuth
//	@Router		/users/{userId} [get]
func GetUser(c *fiber.Ctx) error {
	// check that requesting user is coach/head coach
	db := database.DB
	user := helpers.GetUser(c)

	if !helpers.HasRoleOnClub(c, user.ClubID, helpers.RoleCoach, helpers.RoleHeadCoach) {
		helpers.LogSecurity(fmt.Sprintf("tried reading user %v", user.ID), helpers.GetUserId(c))
		return c.Status(http.StatusForbidden).JSON(res.ErrForbidden)
	}

	// read the requested user id, and verify that it is a numeric string
	userId := c.Params(ParamUserId)
	if _, err := strconv.ParseInt(userId, 10, 64); err != nil {
		helpers.LogInfo(fmt.Sprintf("bad request: %v", err), helpers.GetUserId(c))
		return c.Status(http.StatusBadRequest).JSON(res.ErrBadRequest.Err(err))
	}

	// search the user
	targetUser := models.User{ID: userId}
	err := db.Where(&models.User{ClubID: user.ClubID}).
		Preload("Nationality").Preload("Country").Preload("Club").
		First(&targetUser).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		helpers.LogSecurity(fmt.Sprintf("tried accessing user %v", userId), helpers.GetUserId(c))
		return c.Status(http.StatusForbidden).JSON(res.ErrForbiddenNotFound)
	} else if err != nil {
		helpers.LogError(err.Error(), helpers.GetUserId(c))
		return c.Status(http.StatusInternalServerError).JSON(res.ErrDB.Err(err))
	}

	return c.JSON(targetUser)
}
