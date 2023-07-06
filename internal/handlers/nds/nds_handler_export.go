// Brief   : Data preparator for the NDS system
// Author  : Nicolas Baeriswyl <nicolasfabrice.baeriswyl@edu.hefr.ch>
// Created : 06.07.2023

package ndsHandler

import (
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

func exportAttendances() {
	// read all attendances belonging to a given activity
}

func exportUsers() {
	// read distinct users and verify that the club owning the activity has permissions to read all users
}
