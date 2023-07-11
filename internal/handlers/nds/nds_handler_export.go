// Brief   : Data preparator for the NDS system
// Author  : Nicolas Baeriswyl <nicolasfabrice.baeriswyl@edu.hefr.ch>
// Created : 06.07.2023

package ndsHandler

import (
	"fmt"
	"github.com/google/uuid"
	"gitlab.forge.hefr.ch/tb-athletistics/backend/database"
	"gitlab.forge.hefr.ch/tb-athletistics/backend/helpers"
	"gitlab.forge.hefr.ch/tb-athletistics/backend/internal/models"
	res "gitlab.forge.hefr.ch/tb-athletistics/backend/internal/responses"
	"strings"
)

func exportInstances(user *models.User, activityId string) (string, *res.Error) {
	// read all instances of a given activity
	instances, err := selectInstances(user, activityId)
	if err != nil {
		return "", err
	}

	// construct the string
	csvBuilder := strings.Builder{}
	csvBuilder.WriteString(models.NdsActivityInstanceCsvHeader)
	for i := 0; i < len(instances); i++ {
		csvBuilder.WriteString("\n")
		csvBuilder.WriteString(instances[i].AsCsv())
	}

	return csvBuilder.String(), nil
}

func exportAttendances(user *models.User, activityId string) (string, *res.Error) {
	db := database.DB

	activityUuid, err := uuid.Parse(activityId)
	if err != nil {
		helpers.LogInfo(fmt.Sprintf("bad request: %v", err), user.ID)
		return "", res.ErrBadRequest.Err(err)
	}

	// read all instances of a given activity
	instances, errParam := selectInstances(user, activityId)
	if errParam != nil {
		helpers.LogInfo(fmt.Sprintf("bad request: %v", errParam), user.ID)
		return "", errParam
	}

	instanceIds := make([]uuid.UUID, len(instances))
	for i, obj := range instances {
		instanceIds[i] = obj.ID
	}

	// read all attendances belonging to a given activity
	var attendances []models.NdsAttendance
	err = db.Where("nds_activity_instance_id IN (?)", instanceIds).Preload("NdsActivityInstance").Preload("NdsActivityInstance.NdsActivityType").Find(&attendances).Error
	if err != nil {
		return "", res.NewErrorFromDb(err)
	}

	users, resErr := findUsersOfActivity(user, activityUuid)
	if resErr != nil {
		return "", resErr
	}
	ndsIds, resErr := findNdsIds(users)
	if resErr != nil {
		return "", resErr
	}

	// construct the string
	csvBuilder := strings.Builder{}
	csvBuilder.WriteString(models.NdsAttendanceCsvHeader)
	for i := 0; i < len(attendances); i++ {
		csvBuilder.WriteString("\n")
		csvBuilder.WriteString(attendances[i].AsCsv(ndsIds))
	}

	return csvBuilder.String(), nil
}

func exportUsers(user *models.User, activityId string) (string, *res.Error) {
	// read distinct users and verify that the club owning the activity has permissions to read all users
	activityUuid, err := uuid.Parse(activityId)
	if err != nil {
		helpers.LogInfo(fmt.Sprintf("bad request: %v", err), user.ID)
		return "", res.ErrBadRequest.Err(err)
	}

	if !helpers.HasRoleOnActivity(user, activityUuid, helpers.RoleCoach) {
		helpers.LogSecurity("tried exporting users", user.ID)
		return "", res.ErrForbidden
	}

	users, resErr := findUsersOfActivity(user, activityUuid)
	if resErr != nil {
		return "", resErr
	}
	ndsIds, resErr := findNdsIds(users)
	if resErr != nil {
		return "", resErr
	}

	// construct the string
	csvBuilder := strings.Builder{}
	csvBuilder.WriteString(models.NdsUserCsvHeader)
	for i := 0; i < len(users); i++ {
		csvBuilder.WriteString("\n")
		csvBuilder.WriteString(users[i].AsCsv(ndsIds))
	}

	return csvBuilder.String(), nil
}

func findUsersOfActivity(user *models.User, activityUuid uuid.UUID) ([]models.User, *res.Error) {
	db := database.DB

	if !helpers.HasRoleOnActivity(user, activityUuid, helpers.RoleCoach) {
		helpers.LogSecurity(fmt.Sprintf("tried finding users of activity %v", activityUuid), user.ID)
		return nil, res.ErrForbidden
	}

	// find team of activity
	team := models.Team{}
	err := db.Joins("inner join nds_activities a ON a.team_id = teams.id AND a.id = ?", activityUuid).First(&team).Error
	if err != nil {
		return nil, res.NewErrorFromDb(err)
	}

	// find distinct user IDs taking part in at least one activity
	var users []models.User
	if err := db.Model(&models.User{}).Distinct("users.id").
		Joins("INNER JOIN nds_attendances a ON a.user_id = users.id").
		Joins("INNER JOIN nds_activity_instances i ON a.nds_activity_instance_id = i.id AND i.nds_activity_id = ?", activityUuid).
		Scan(&users).Error; err != nil {
		return nil, res.NewErrorFromDb(err)
	}

	userIds := make([]string, len(users))
	for i := 0; i < len(users); i++ {
		userIds[i] = users[i].ID
	}

	// find distinct users having a role on the club of the activity
	var nbRecords int64
	query := models.Authorisation{ClubID: team.ClubID}
	db.Table("authorisations").Distinct("user_id").Where(&query).Where("user_id IN (?)", userIds).Count(&nbRecords)

	if int64(len(users)) != nbRecords {
		helpers.LogSecurity(fmt.Sprintf("tried accessing users of activity %v", activityUuid), user.ID)
		return nil, res.ErrForbidden
	}

	if err := db.Preload("Nationality").Preload("Gender").Preload("Country").Preload("Language").Find(&users, userIds).Error; err != nil {
		return nil, res.NewErrorFromDb(err)
	}

	return users, nil
}

// findNdsIds searches for external IDs of type NDS. This function has no security
// checks, access to given users must be checked first.
func findNdsIds(users []models.User) (map[string]string, *res.Error) {
	// find list of external IDs
	db := database.DB
	userIds := make([]string, len(users))
	for i := 0; i < len(users); i++ {
		userIds[i] = users[i].ID
	}

	var ids []models.ExternalId
	err := db.Where("internal_id IN (?)", userIds).Where("external_id_type_id = ?", models.ExternalIdTypeNds).Find(&ids).Error
	if err != nil {
		return nil, res.NewErrorFromDb(err)
	}

	// convert to map with key being the internal ID, value being the NDS id
	results := make(map[string]string)
	for i := 0; i < len(ids); i++ {
		results[ids[i].InternalID] = ids[i].ExternalID
	}

	return results, nil
}
