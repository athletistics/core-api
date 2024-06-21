//go:build testing

// Brief   : A test helper file that should not be contained in the build.
// Author  : Nicolas Baeriswyl <nicolasfabrice.baeriswyl@edu.hefr.ch>
// Created : 16.06.2024

package helpers

import (
	"github.com/athletistics/core-api/internal/models"
	"github.com/google/uuid"
)

// CreateContext creates a user requesting context that MUST ONLY BE USED IN TESTING FUNCTIONS!
// THIS FUNCTION IS NOT AVAILABLE IN PRODUCTION!
func CreateContext(u models.User, authByClub map[uuid.UUID]map[Role]bool, state TwoFAState) *Context {
	return &Context{user: &u, authByClub: authByClub, twoFA: state}
}

func (c *Context) SetAuthMap(authMap map[uuid.UUID]map[Role]bool) {
	c.authByClub = authMap
}
