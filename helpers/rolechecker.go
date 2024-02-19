// Brief   : Checking roles of a user
// Author  : Nicolas Baeriswyl <nicolasfabrice.baeriswyl@edu.hefr.ch>
// Created : 20.06.2023

package helpers

import (
	"github.com/athletistics/core-api/database"
	"github.com/athletistics/core-api/internal/models"
	"github.com/google/uuid"
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
func HasRoleOnClub(user *models.User, clubId uuid.UUID, roles ...Role) bool {
	if len(roles) == 0 || user == nil {
		return false
	}

	authorisations, err := GetAuthorisations(user.ID, clubId)
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

func HasRoleOnParentClub(user *models.User, roles ...Role) bool {
	return HasRoleOnClub(user, user.ClubID, roles...)
}

// HasAnyRoleOnParentClub returns true if the user has any role on his parent club.
// Short for HasAnyRoleOnClub(c, helperse.GetUser(c).ClubID)
func HasAnyRoleOnParentClub(user *models.User) bool {
	return HasAnyRoleOnClub(user, user.ClubID)
}

// HasAnyRoleOnClub checks if the requesting user has any role on the specificed club.
func HasAnyRoleOnClub(user *models.User, clubId uuid.UUID) bool {
	if user == nil {
		return false
	}

	authorisations, err := GetAuthorisations(user.ID, clubId)
	if err != nil {
		return false
	}

	return len(authorisations) > 0
}

// HasRoleOnTeam is similar to HasRoleOnClub, but checks the roles for the club of a specific team.
func HasRoleOnTeam(user *models.User, teamId uuid.UUID, roles ...Role) bool {
	db := database.DB

	team := models.Team{ID: teamId}
	if err := db.First(&team).Error; err != nil {
		return false
	}
	return HasRoleOnClub(user, team.ClubID, roles...)
}

// HasRoleOnActivity is similar to HasRoleOnClub, but checks the roles for the club of a specific team.
func HasRoleOnActivity(user *models.User, activityId uuid.UUID, roles ...Role) bool {
	db := database.DB

	team := models.Team{}
	err := db.
		Joins("inner join nds_activities a ON a.team_id = teams.id AND a.id = ?", activityId).
		First(&team).Error

	if err != nil {
		return false
	}
	return HasRoleOnClub(user, team.ClubID, roles...)
}

// HasRoleOnOneOfTheClubs checks for overlapping in the source and target user club authorisations. The function
// returns true if either the source user has one of the required roles on the target user's parent club, or one of
// the target user's authorised clubs (any access role).
func HasRoleOnOneOfTheClubs(user *models.User, targetUser *models.User, roles ...Role) bool {
	db := database.DB

	var authorisations []models.Authorisation
	if err := db.Where(&models.Authorisation{UserID: user.ID}).Preload("AccessRole").Find(&authorisations).Error; err != nil {
		return false
	}

	var targetAuthorisations []models.Authorisation
	if err := db.Where(&models.Authorisation{UserID: targetUser.ID}).Preload("AccessRole").Find(&targetAuthorisations).Error; err != nil {
		return false
	}

	userClubIds := getDistinctClubs(authorisations, roles...)
	targetUserClubIds := getDistinctClubs(targetAuthorisations)

	for _, id := range userClubIds {
		if id == targetUser.ClubID {
			return true
		}

		for _, targetId := range targetUserClubIds {
			if targetId == id {
				return true
			}
		}
	}

	return false
}

// IsAthlete checks if the requesting user has the role "athlete" on any club.
func IsAthlete(user *models.User) bool {
	db := database.DB

	if user == nil {
		return false
	}

	var authorisations []models.Authorisation
	err := db.Where(&models.Authorisation{UserID: user.ID}).Where(&models.AccessRole{I18nName: RoleAthlete.String()}).Find(&authorisations).Error
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

// getDistinctClubs reads a list of authorisations and returns all clubs where at least one of the given roles
// are available. If no role is given, all clubs will be returned instead.
func getDistinctClubs(authorisations []models.Authorisation, roles ...Role) []uuid.UUID {
	clubs := make(map[uuid.UUID]struct{})

	if len(roles) == 0 {
		for _, auth := range authorisations {
			clubs[auth.ClubID] = struct{}{}
		}
	} else {
		for _, auth := range authorisations {
			for _, role := range roles {
				if auth.AccessRole.I18nName == role.String() {
					clubs[auth.ClubID] = struct{}{}
				}
			}
		}
	}

	var ids []uuid.UUID
	for k := range clubs {
		ids = append(ids, k)
	}

	return ids
}
