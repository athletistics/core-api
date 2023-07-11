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
	QueryCode   = "code"
	QueryUser   = "userId"
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

	if !helpers.HasRoleOnClub(&user, user.ClubID, helpers.RoleCoach, helpers.RoleHeadCoach) {
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

	if !helpers.HasRoleOnClub(&user, user.ClubID, helpers.RoleCoach, helpers.RoleHeadCoach) {
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
	db := database.DB
	user := helpers.GetUser(c)

	// check that requesting user is head coach
	if !helpers.HasRoleOnClub(&user, user.ClubID, helpers.RoleHeadCoach) {
		helpers.LogSecurity("tried creating user (missing required role)", helpers.GetUserId(c))
		return c.Status(http.StatusForbidden).JSON(res.ErrForbidden)
	}

	// extract values from body
	newUser := new(models.User)
	if err := c.BodyParser(newUser); err != nil {
		helpers.LogInfo(fmt.Sprintf("bad request: %v", err), helpers.GetUserId(c))
		return c.Status(http.StatusBadRequest).JSON(res.ErrBadRequest.Err(err))
	}

	// create the user on zitadel if necessary. generate an ID otherwise
	if newUser.WithLogin {
		id, err := createUser(newUser)
		newUser.ID = id
		if err != nil {
			helpers.LogError(err.Error(), helpers.GetUserId(c))
			return c.Status(http.StatusInternalServerError).JSON(res.ErrBadRequest.Err(err))
		}
	} else {
		id, err := helpers.IdGen.NextID()
		if err != nil {
			helpers.LogError(err.Error(), helpers.GetUserId(c))
			return c.Status(http.StatusInternalServerError).JSON(res.ErrBadRequest.Err(err))
		}
		newUser.ID = strconv.FormatUint(id, 10)
	}

	err := db.Transaction(func(tx *gorm.DB) error {
		// create the user locally
		if err := tx.Create(newUser).Error; err != nil {
			return err
		}

		// assign 'athlete' role on parent club
		tx.Create(&models.Authorisation{AccessRoleID: models.AccessRoleIdAthlete, UserID: newUser.ID, ClubID: newUser.ClubID})

		// commit
		return nil
	})
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

	if !helpers.HasRoleOnClub(&user, user.ClubID, helpers.RoleCoach, helpers.RoleHeadCoach) {
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

	// users with logins read the login state from Zitadel
	if targetUser.WithLogin {
		state, err := GetUserState(targetUser.ID)
		if err != nil {
			helpers.LogError(err.Error(), helpers.GetUserId(c))
			return c.Status(http.StatusInternalServerError).JSON(res.ErrZitadel.Err(err))
		}
		targetUser.UserState = state
	} else {
		targetUser.UserState = UserStateNoLogin
	}

	return c.JSON(targetUser)
}

// GetShareCode issues a single-use code allowing some other club to add the user to his club. The code is valid for 24 hours.
//
//	@Summary	Get a single-use code allowing some other club to add the user to his club. The code is valid for 24 hours.
//	@Tags		User share
//	@Produce	json
//	@Param		userId	query		string	true	"User identifier"
//	@Success	200		{object}	models.UserShare
//	@Failure	401		{object}	res.Error
//	@Failure	403		{object}	res.Error	"Not enough permissions"
//	@Security	ApiKeyAuth
//	@Router		/users/share [get]
func GetShareCode(c *fiber.Ctx) error {
	user := helpers.GetUser(c)
	targetUserId := c.Query(QueryUser)

	data, err := createShareCode(&user, targetUserId)
	return helpers.ProcessJsonResponse(c, data, err)
}

// RedeemShareCode redeems a code and adds a user as athlete to the club.
//
//	@Summary	Add a user as athlete to the club using a share code.
//	@Tags		User share
//	@Consumes	json
//	@Param		code	query		string	true	"Share code"
//	@Success	200		{object}	models.UserShort
//	@Failure	401		{object}	res.Error
//	@Failure	403		{object}	res.Error	"Not enough permissions"
//	@Failure	404		{object}	res.Error	"Invalid code"
//	@Security	ApiKeyAuth
//	@Router		/users/share [post]
func RedeemShareCode(c *fiber.Ctx) error {
	user := helpers.GetUser(c)
	code := c.Query(QueryCode)

	data, err := redeemShareCode(&user, code)
	return helpers.ProcessJsonResponse(c, data, err)
}

// GetRoles lists all role of a specific user that the requesting user has permissions to read.
//
//	@Summary	Add an authorisation role to a specific user.
//	@Tags		User auth
//	@Produce	json
//	@Param		userId	query		string	true	"Target user ID"
//	@Success	200		{object}	[]models.Authorisation
//	@Failure	401		{object}	res.Error
//	@Failure	403		{object}	res.Error	"Not enough permissions"
//	@Security	ApiKeyAuth
//	@Router		/users/{userId}/roles [get]
func GetRoles(c *fiber.Ctx) error {
	user := helpers.GetUser(c)
	targetUserId := c.Params(ParamUserId)

	data, err := getRoles(&user, targetUserId)
	return helpers.ProcessJsonResponse(c, data, err)
}

// AddRole adds an authorisation role to a specific user.
//
//	@Summary	Add an authorisation role to a specific user.
//	@Tags		User auth
//	@Consume	json
//	@Param		request	body		models.Authorisation	true	"Authorisation to add"
//	@Success	200		{string}	string
//	@Failure	401		{object}	res.Error
//	@Failure	403		{object}	res.Error	"Not enough permissions"
//	@Security	ApiKeyAuth
//	@Router		/users/{userId}/roles [post]
func AddRole(c *fiber.Ctx) error {
	user := helpers.GetUser(c)
	b := new(models.Authorisation)
	parseErr := c.BodyParser(b)
	targetUserId := c.Params(ParamUserId)

	err := addRole(&user, targetUserId, b, parseErr)
	return helpers.ProcessJsonResponse(c, nil, err)
}

// RemoveRole returns details of a specific user.
//
//	@Summary	Remove an authorisation role from a specific user.
//	@Tags		User auth
//	@Produce	json
//	@Param		request	body		models.Authorisation	true	"Authorisation to remove"
//	@Success	200		{string}	string
//	@Failure	401		{object}	res.Error
//	@Failure	403		{object}	res.Error	"Not enough permissions"
//	@Security	ApiKeyAuth
//	@Router		/users/{userId}/roles [delete]
func RemoveRole(c *fiber.Ctx) error {
	user := helpers.GetUser(c)
	b := new(models.Authorisation)
	parseErr := c.BodyParser(b)
	targetUserId := c.Params(ParamUserId)

	err := removeRole(&user, targetUserId, b, parseErr)
	return helpers.ProcessJsonResponse(c, nil, err)
}
