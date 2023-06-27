// Brief   :
// Author  : Nicolas Baeriswyl <nicolasfabrice.baeriswyl@edu.hefr.ch>
// Created : 20.06.2023

package helpers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gitlab.forge.hefr.ch/tb-athletistics/backend/database"
	"gitlab.forge.hefr.ch/tb-athletistics/backend/internal/models"
)

type Role int64

const (
	RoleAthlete Role = iota
	RoleCoach
	RoleHeadCoach
)

func (r Role) String() string {
	switch r {
	case RoleAthlete:
		return "athlete"
	case RoleCoach:
		return "coach"
	case RoleHeadCoach:
		return "head_coach"
	}
	return "unknown"
}

// HasRoleOnClub checks if the requesting user has at least one of the specified roles on the specified club.
func HasRoleOnClub(c *fiber.Ctx, clubId uuid.UUID, roles ...Role) bool {
	if len(roles) == 0 {
		return false
	}

	authorisations, err := GetAuthorisations(GetUserId(c), clubId)
	if err != nil {
		return false
	}

	for i := 0; i < len(authorisations); i++ {
		for j := 0; j < len(roles); j++ {
			if authorisations[i].AccessRole.I18nName == roles[j].String() {
				return true
			}
		}
	}

	return false
}

func HasRoleOnParentClub(c *fiber.Ctx, roles ...Role) bool {
	return HasRoleOnClub(c, GetUser(c).ClubID, roles...)
}

// HasAnyRoleOnParentClub returns true if the user has any role on his parent club.
// Short for HasAnyRoleOnClub(c, helperse.GetUser(c).ClubID)
func HasAnyRoleOnParentClub(c *fiber.Ctx) bool {
	return HasAnyRoleOnClub(c, GetUser(c).ClubID)
}

// HasAnyRoleOnClub checks if the requesting user has any role on the specificed club.
func HasAnyRoleOnClub(c *fiber.Ctx, clubId uuid.UUID) bool {
	authorisations, err := GetAuthorisations(GetUserId(c), clubId)
	if err != nil {
		return false
	}

	return len(authorisations) > 0
}

// IsAthlete checks if the requesting user has the role "athlete" on any club.
func IsAthlete(c *fiber.Ctx) bool {
	db := database.DB

	var authorisations []models.Authorisation
	err := db.Where(&models.Authorisation{UserID: GetUserId(c)}).Where(&models.AccessRole{I18nName: RoleAthlete.String()}).Find(&authorisations).Error
	if err != nil {
		return false
	}

	return len(authorisations) > 0
}

// GetAuthorisations lists all authorisations of the requesting user on a specific club.
func GetAuthorisations(userId string, clubId uuid.UUID) ([]models.Authorisation, error) {
	db := database.DB

	var authorisations []models.Authorisation
	err := db.Where(&models.Authorisation{UserID: userId, ClubID: clubId}).Preload("AccessRole").Find(&authorisations).Error
	return authorisations, err
}
