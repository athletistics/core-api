// Brief   :
// Author  : Nicolas Baeriswyl <nicolasfabrice.baeriswyl@edu.hefr.ch>
// Created : 13.06.2024

package helpers

import (
	"fmt"
	"github.com/athletistics/core-api/config"
	"github.com/athletistics/core-api/constants"
	"github.com/athletistics/core-api/database"
	"github.com/athletistics/core-api/internal/models"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"net/http"
	"os"
)

// roles
const (
	RoleAny Role = iota
	RoleAthlete
	RoleCoach
	RoleHeadCoach
	RoleUnknown
)

func (r Role) String() string {
	switch r {
	case RoleAny:
		return "any"
	case RoleAthlete:
		return "athlete"
	case RoleCoach:
		return "coach"
	case RoleHeadCoach:
		return "head_coach"
	default:
		return "unknown"
	}
}

func RoleFromI18nName(name string) (Role, bool) {
	switch name {
	case "athlete":
		return RoleAthlete, true
	case "coach":
		return RoleCoach, true
	case "head_coach":
		return RoleHeadCoach, true
	default:
		return RoleUnknown, false
	}
}

const (
	AuthHeaderName = "Authorization"
)

type ClubOwned interface {
	GetOwningClub() uuid.UUID
}

type TwoFAState struct {
	exists bool
	valid  bool
}

type Context struct {
	user       *models.User
	authByClub map[uuid.UUID]map[Role]bool
	twoFA      TwoFAState
}

func UserContext(c *fiber.Ctx) *Context {
	return c.Locals(constants.LocalsContext).(*Context)
}

// LoadContext generates, based on the UserId stored in fiber's Locals, the requesting user's context. Loads the user
// model, the authorization map, and the 2FA state.
func LoadContext(c *fiber.Ctx) *Context {
	ctx := &Context{authByClub: make(map[uuid.UUID]map[Role]bool), twoFA: newTwoFAState()}
	userId := c.Locals(constants.LocalsUserId)
	if userId == nil {
		return ctx
	}

	ctx.user = getUser(userId.(string))
	if ctx.user != nil {
		ctx.twoFA = getTwoFA(c)
		ctx.authByClub = getAuthsByClub(ctx.user)
	}

	return ctx
}

func (c *Context) IsAuthenticated() bool {
	return c.user != nil
}

func (c *Context) User() *models.User {
	return c.user
}

func (c *Context) HasRoleOnClub(club uuid.UUID, roles ...Role) bool {
	if !c.IsAuthenticated() {
		return false
	}

	authorizations, ok := c.authByClub[club]

	if !ok {
		return false
	}

	if roles[0] == RoleAny && len(authorizations) > 0 {
		return true
	}

	for _, role := range roles {
		if authorizations[role] {
			return true
		}
	}

	return false
}

// HasRoleOnParentClub returns true, if the requesting user has at least one of the provided rules on his parent club.
func (c *Context) HasRoleOnParentClub(roles ...Role) bool {
	return c.HasRoleOnClub(c.user.ClubID, roles...)
}

func (c *Context) HasRoleOnClubOwned(target ClubOwned, roles ...Role) bool {
	return c.HasRoleOnClub(target.GetOwningClub(), roles...)
}

// -------------------
// 2FA state functions
// -------------------
func newTwoFAState() TwoFAState {
	return TwoFAState{false, false}
}

func (c *Context) TwoFA() TwoFAState {
	return c.twoFA
}

func (t *TwoFAState) IsValid() bool {
	return t.valid
}

func (t *TwoFAState) Exists() bool {
	return t.exists
}

// getTwoFA reaches out to the TOTP microservice and requests the current user's 2FA state by reusing the
// authorization header. In case of any errors, the 2FA is considered to not exist and not to be verified.
func getTwoFA(c *fiber.Ctx) TwoFAState {
	tfa := TwoFAState{exists: false, valid: false}

	authHeader, ok := c.GetReqHeaders()[AuthHeaderName]
	if !ok || len(authHeader) < 1 {
		return tfa
	}

	client := &http.Client{}
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%v/verify?%v=%v", os.Getenv("TOTP_ENDPOINT"), constants.KeySecondsAgo, config.TotpSecondsLimit), nil)
	if err != nil {
		return tfa
	}

	req.Header.Set(AuthHeaderName, authHeader[0])
	res, err := client.Do(req)

	if res == nil || err != nil {
		return tfa
	}

	// mark 2FA as verified if response has OK code
	tfa.valid = res.StatusCode == http.StatusOK
	tfa.exists = res.StatusCode == http.StatusExpectationFailed || res.StatusCode == http.StatusOK
	_ = res.Body.Close()

	return tfa
}

// getUser reads the user model identified by the ID stored in context's locals "LocalUserId" from the database.
// If the ID is not present, or the requested user does not exist, nil is returned instead.
func getUser(userId string) *models.User {
	user := models.User{ID: userId}
	if err := database.DB.First(&user).Error; err != nil {
		return nil
	}

	return &user
}

// getAuthsByClub reads all the access roles of a specific user, and puts them into a map of map of booleans.
// A user has a role on a specific club, if the Role entry exists on the club, and has value true.
func getAuthsByClub(u *models.User) map[uuid.UUID]map[Role]bool {
	result := make(map[uuid.UUID]map[Role]bool)

	if u == nil {
		return result
	}

	db := database.DB

	var authorisations []models.Authorisation
	err := db.Where(&models.Authorisation{UserID: u.ID}).Preload("AccessRole").Find(&authorisations).Error

	if err != nil {
		return result
	}

	for _, auth := range authorisations {
		role, ok := RoleFromI18nName(auth.AccessRole.I18nName)
		if !ok {
			continue
		}

		if _, ok := result[auth.ClubID]; !ok {
			result[auth.ClubID] = make(map[Role]bool)
		}

		result[auth.ClubID][role] = true
	}

	return result
}
