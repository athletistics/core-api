// Brief   : User management handlers
// Author  : Nicolas Baeriswyl <nicolasfabrice.baeriswyl@edu.hefr.ch>
// Created : 20.06.2023

package userHandler

import (
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
	ParamUserId    = "id"
	QueryCode      = "code"
	QueryUser      = "userId"
	QuerySeparated = "separated"
)

// GetUsers lists all users owned by the authenticated user's club, and shared with the club.
//
//	@Summary	List users
//	@Tags		User info
//	@Produce	json
//
//	@Param		separated	query		boolean								false	"Flag defining if users are separated as 'owned by club' and 'shared' or not"
//
//	@Success	200			{object}	[]models.User						"owned and shared users combined ('separated' flag omitted or false)"
//	@Success	200			{object}	userHandler.GetUsers.usersSeparated	"owned and shared users separated ('separated' flag true)"
//	@Failure	401			{object}	res.Error
//	@Failure	403			{object}	res.Error	"Not enough permissions"
//	@Security	ApiKeyAuth
//	@Router		/users [get]
func GetUsers(c *fiber.Ctx) error {
	type usersSeparated struct {
		Owned  []models.User `json:"owned"`
		Shared []models.User `json:"shared"`
	}

	db := database.DB
	user := helpers.GetUser(c)
	separated := c.QueryBool(QuerySeparated)

	if !helpers.HasRoleOnClub(&user, user.ClubID, helpers.RoleCoach, helpers.RoleHeadCoach) {
		helpers.LogSecurity(fmt.Sprintf("tried accessing users of club %v", user.ClubID), helpers.GetUserId(c))
		return c.Status(http.StatusForbidden).JSON(res.ErrForbidden)
	}

	if separated {
		var usersOwned []models.User
		err := db.Where(&models.User{ClubID: user.ClubID}).Find(&usersOwned).Error
		if err != nil {
			helpers.LogError(err.Error(), helpers.GetUserId(c))
			return c.Status(http.StatusInternalServerError).JSON(res.ErrDB.Err(err))
		}

		// find all users with authorisation on the club, but not owned by the club
		var usersShared []models.User
		err = db.InnerJoins(`INNER JOIN authorisations a ON a.user_id = users.id AND a.club_id = ? `, user.ClubID).
			Where(`users.club_id <> ?`, user.ClubID).
			Group("users.id"). // group to get distinct users
			Find(&usersShared).Error
		if err != nil {
			helpers.LogError(err.Error(), helpers.GetUserId(c))
			return c.Status(http.StatusInternalServerError).JSON(res.ErrDB.Err(err))
		}

		return c.JSON(usersSeparated{usersOwned, usersShared})
	} else {
		var users []models.User
		err := db.InnerJoins(`LEFT JOIN authorisations a ON a.user_id = users.id AND a.club_id = ? `, user.ClubID).
			Where(`users.club_id = ? OR a.club_id = ?`, user.ClubID, user.ClubID).
			Group("users.id"). // group to get distinct users
			Find(&users).Error
		if err != nil {
			helpers.LogError(err.Error(), helpers.GetUserId(c))
			return c.Status(http.StatusInternalServerError).JSON(res.ErrDB.Err(err))
		}

		return c.JSON(users)
	}
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
	err := db.Model(&models.User{}).
		InnerJoins(`LEFT JOIN authorisations a ON a.user_id = users.id AND a.club_id = ? `, user.ClubID).
		Where(`users.club_id = ? OR a.club_id = ?`, user.ClubID, user.ClubID).
		Group("users.id").
		Find(&users).Error
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
//	@Param		userId	path		string	true	"User ID"
//	@Success	200		{object}	models.User
//	@Failure	401		{object}	res.Error
//	@Failure	403		{object}	res.Error	"Not enough permissions"
//	@Security	ApiKeyAuth
//	@Router		/users/{userId} [get]
func GetUser(c *fiber.Ctx) error {
	// check that requesting user is coach/head coach
	user := helpers.GetUser(c)

	// find the user
	userId := c.Params(ParamUserId)
	targetUser, httpStatus, err := getUser(user.ID, userId)
	if httpStatus != http.StatusOK {
		return c.Status(httpStatus).JSON(err)
	}

	if !helpers.HasRoleOnOneOfTheClubs(&user, targetUser, helpers.RoleCoach, helpers.RoleHeadCoach) {
		helpers.LogSecurity(fmt.Sprintf("tried reading user %v", user.ID), helpers.GetUserId(c))
		return c.Status(http.StatusForbidden).JSON(res.ErrForbidden)
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

// UpdateUser updates a user to a new state that is being passed in the body. Values that are not passed will be set to
// empty values in the database. Thus, the entire user must be sent at any time.
//
//	@Summary	Update a user to a new state that is being passed in the body. Values that are not passed will be set to
//	            empty values in the database. Thus, the entire user must be sent at any time.
//
//	@Tags		User share
//	@Produce	json
//	@Param		userId	path		string		true	"User ID"
//	@Param		request	body		models.User	true	"Updated user data"
//	@Success	200		{object}	models.UserShare
//	@Failure	401		{object}	res.Error
//	@Failure	403		{object}	res.Error	"Not enough permissions"
//	@Security	ApiKeyAuth
//	@Router		/users/{userId} [put]
func UpdateUser(c *fiber.Ctx) error {
	b := new(models.User)
	parseErr := c.BodyParser(b)
	userId := c.Params(ParamUserId)
	user := helpers.GetUser(c)

	return helpers.ProcessJsonResponse(c, nil, updateUser(&user, userId, b, parseErr))
}

// GetShareCode issues a single-use code allowing some other club to add the user to his club. The code is valid for 24 hours.
//
//	@Summary	Get a single-use code allowing some other club to add the user to his club. The code is valid for 24 hours.
//	@Tags		User share
//	@Produce	json
//	@Param		userId	query		string	true	"User identifier"
//	@Success	200		{string}	string
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

// CreateLogin asks Zitadel to initialize a login for a given user.
//
//	@Summary	Ask Zitadel to initialize a login for a given user.
//	@Tags		User auth
//	@Consumes	json
//	@Param		userId	path		string		true	"User ID"
//	@Success	200		{object}
//	@Failure	400		{object}	res.Error	"The user already has a login"
//	@Failure	401		{object}	res.Error
//	@Failure	403		{object}	res.Error	"Forbidden"
//	@Security	ApiKeyAuth
//	@Router		/users/{userId}/create-login [post]
func CreateLogin(c *fiber.Ctx) error {
	userId := c.Params(ParamUserId)
	user := helpers.GetUser(c)

	// search the user
	targetUser, httpStatus, err := getUser(user.ID, userId)
	if httpStatus != http.StatusOK {
		return c.Status(httpStatus).JSON(err)
	}

	if !helpers.HasRoleOnOneOfTheClubs(&user, targetUser, helpers.RoleHeadCoach) {
		helpers.LogSecurity(fmt.Sprintf("tried creating a login for user %v", user.ID), helpers.GetUserId(c))
		return c.Status(http.StatusForbidden).JSON(res.ErrForbidden)
	}

	httpStatus, err = createUserLogin(targetUser)

	if httpStatus != http.StatusOK {
		return c.Status(httpStatus).JSON(err)
	}

	return c.SendStatus(http.StatusOK)
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
