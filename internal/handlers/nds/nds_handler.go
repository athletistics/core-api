// Brief   :
// Author  : Nicolas Baeriswyl <nicolasfabrice.baeriswyl@edu.hefr.ch>
// Created : 06.07.2023

package ndsHandler

import (
	"github.com/gofiber/fiber/v2"
	"gitlab.forge.hefr.ch/tb-athletistics/backend/helpers"
	"gitlab.forge.hefr.ch/tb-athletistics/backend/internal/models"
)

const (
	ParamActivityId = "activityId"
	ParamInstanceId = "instanceId"
)

const (
	QueryTeamId = "teamId"
)

const DefaultActivityTypeId = 1

func GetTypes(c *fiber.Ctx) error {
	data, err := selectTypes()
	return helpers.ProcessResponse(c, data, err)
}

func GetActivities(c *fiber.Ctx) error {
	// decode requested team
	teamId := c.Query(QueryTeamId)
	user := helpers.GetUser(c)

	data, err := selectActivities(&user, teamId)
	return helpers.ProcessResponse(c, data, err)
}

func GetActivity(c *fiber.Ctx) error {
	// decode requested activity
	activityId := c.Params(ParamActivityId)
	user := helpers.GetUser(c)

	data, err := selectActivity(&user, activityId)
	return helpers.ProcessResponse(c, data, err)
}

func CreateActivity(c *fiber.Ctx) error {
	// parse request body
	b := new(models.NdsActivityInput)
	parseErr := c.BodyParser(b)
	user := helpers.GetUser(c)

	data, err := createActivity(&user, b, parseErr)
	return helpers.ProcessResponse(c, data, err)
}

func GetInstances(c *fiber.Ctx) error {
	// parse parameters
	activityId := c.Params(ParamActivityId)
	user := helpers.GetUser(c)

	data, err := selectInstances(&user, activityId)
	return helpers.ProcessResponse(c, data, err)
}

func GetInstance(c *fiber.Ctx) error {
	// parse request
	activityId := c.Params(ParamActivityId)
	instanceId := c.Params(ParamInstanceId)
	user := helpers.GetUser(c)

	data, err := selectInstance(&user, activityId, instanceId)
	return helpers.ProcessResponse(c, data, err)
}

func CreateInstance(c *fiber.Ctx) error {
	// parse request
	b := new(models.NdsActivityInstanceInput)
	parseErr := c.BodyParser(b)
	activityId := c.Params(ParamActivityId)
	user := helpers.GetUser(c)

	data, err := createInstance(&user, activityId, b, parseErr)
	return helpers.ProcessResponse(c, data, err)
}

func DeleteInstance(c *fiber.Ctx) error {
	// parse request
	activityId := c.Params(ParamActivityId)
	instanceId := c.Params(ParamInstanceId)
	user := helpers.GetUser(c)

	err := deleteInstance(&user, activityId, instanceId)
	return helpers.ProcessResponse(c, nil, err)
}

func GetParticipants(c *fiber.Ctx) error {
	// parse request
	activityId := c.Params(ParamActivityId)
	instanceId := c.Params(ParamInstanceId)
	user := helpers.GetUser(c)

	data, err := selectAttendances(&user, activityId, instanceId)
	return helpers.ProcessResponse(c, data, err)
}

func SetParticipants(c *fiber.Ctx) error {
	// parse request
	b := new([]models.NdsAttendanceInput)
	parseErr := c.BodyParser(b)

	activityId := c.Params(ParamActivityId)
	instanceId := c.Params(ParamInstanceId)
	user := helpers.GetUser(c)

	err := setParticipants(&user, activityId, instanceId, b, parseErr)
	return helpers.ProcessResponse(c, nil, err)
}
