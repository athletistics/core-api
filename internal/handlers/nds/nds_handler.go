// Brief   : Preparing all endpoints for the NDS management
// Author  : Nicolas Baeriswyl <nicolasfabrice.baeriswyl@edu.hefr.ch>
// Created : 06.07.2023

package ndsHandler

import (
	"github.com/athletistics/core-api/helpers"
	"github.com/athletistics/core-api/internal/models"
	"github.com/gofiber/fiber/v2"
)

const (
	ParamActivityId = "activityId"
	ParamInstanceId = "instanceId"
)

const (
	QueryTeamId = "teamId"
)

const DefaultActivityTypeId = 1

// GetTypes lists all available NDS activity types
//
//	@Summary	List all available NDS activity types
//	@Tags		NDS Activity
//	@Produce	json
//	@Success	200	{object}	[]models.NdsActivityType
//	@Security	ApiKeyAuth
//	@Router		/nds/types [get]
func GetTypes(c *fiber.Ctx) error {
	data, err := selectTypes()
	return helpers.ProcessJsonResponse(c, data, err)
}

// GetActivities lists all activities for a specific team.
//
//	@Summary	List all activities for a specific team.
//	@Tags		NDS Activity
//	@Produce	json
//	@Param		teamId	query		string	true	"Team id (UUID format)"
//	@Success	200		{object}	[]models.NdsActivity
//	@Failure	400		{object}	res.Error
//	@Failure	401		{object}	res.Error
//	@Failure	403		{object}	res.Error	"Not enough permissions"
//	@Security	ApiKeyAuth
//	@Router		/nds/activities [get]
func GetActivities(c *fiber.Ctx) error {
	// decode requested team
	teamId := c.Query(QueryTeamId)
	user := helpers.GetUser(c)

	data, err := selectActivities(&user, teamId)
	return helpers.ProcessJsonResponse(c, data, err)
}

// GetActivity reads a specific activity.
//
//	@Summary	Read a specific activity.
//	@Tags		NDS Activity
//	@Produce	json
//	@Param		activityId	path		string	true	"Activity ID (UUID format)"
//	@Success	200			{object}	models.NdsActivity
//	@Failure	400			{object}	res.Error
//	@Failure	401			{object}	res.Error
//	@Failure	403			{object}	res.Error	"Not enough permissions"
//	@Security	ApiKeyAuth
//	@Router		/nds/activities/{activityId} [get]
func GetActivity(c *fiber.Ctx) error {
	// decode requested activity
	activityId := c.Params(ParamActivityId)
	user := helpers.GetUser(c)

	data, err := selectActivity(&user, activityId)
	return helpers.ProcessJsonResponse(c, data, err)
}

// CreateActivity creates a new activity and instantiates the activity once a week
// on the defined weekday between start and end date.
//
//	@Summary	Create a new activity and instantiate the activity once a week on the defined weekday between start and end date.
//	@Tags		NDS Activity
//	@Consume	json
//	@Produce	json
//	@Param		request	body		models.NdsActivityInput	true	"Activity data"
//	@Success	200		{object}	res.UuidCreated
//	@Failure	400		{object}	res.Error
//	@Failure	401		{object}	res.Error
//	@Failure	403		{object}	res.Error	"Not enough permissions"
//	@Security	ApiKeyAuth
//	@Router		/nds/activities [post]
func CreateActivity(c *fiber.Ctx) error {
	// parse request body
	b := new(models.NdsActivityInput)
	parseErr := c.BodyParser(b)
	user := helpers.GetUser(c)

	data, err := createActivity(&user, b, parseErr)
	return helpers.ProcessJsonResponse(c, data, err)
}

// GetInstances lists all instances of a specific activity.
//
//	@Summary	Lists all instances of a specific activity.
//	@Tags		NDS Instance
//	@Produce	json
//	@Param		activityId	path		string	true	"Activity ID (UUID format)"
//	@Success	200			{object}	[]models.NdsActivityInstance
//	@Failure	400			{object}	res.Error
//	@Failure	401			{object}	res.Error
//	@Failure	403			{object}	res.Error	"Not enough permissions"
//	@Security	ApiKeyAuth
//	@Router		/nds/activities/{activityId}/instances [get]
func GetInstances(c *fiber.Ctx) error {
	// parse parameters
	activityId := c.Params(ParamActivityId)
	user := helpers.GetUser(c)

	data, err := selectInstances(&user, activityId)
	return helpers.ProcessJsonResponse(c, data, err)
}

// GetInstance reads details of a specific instance.
//
//	@Summary	Read details of a specific instance.
//	@Tags		NDS Instance
//	@Produce	json
//	@Param		activityId	path		string	true	"Activity ID (UUID format)"
//	@Param		instanceId	path		string	true	"Instance ID (UUID format)"
//	@Success	200			{object}	models.NdsActivityInstance
//	@Failure	400			{object}	res.Error
//	@Failure	401			{object}	res.Error
//	@Failure	403			{object}	res.Error	"Not enough permissions"
//	@Security	ApiKeyAuth
//	@Router		/nds/activities/{activityId}/instances/{instanceId} [get]
func GetInstance(c *fiber.Ctx) error {
	// parse request
	activityId := c.Params(ParamActivityId)
	instanceId := c.Params(ParamInstanceId)
	user := helpers.GetUser(c)

	data, err := selectInstance(&user, activityId, instanceId)
	return helpers.ProcessJsonResponse(c, data, err)
}

// CreateInstance creates a new instance for a specific activity.
//
//	@Summary	Create a new instance for a specific activity.
//	@Tags		NDS Instance
//	@Consume	json
//	@Produce	json
//	@Param		activityId	path		string							true	"Activity ID (UUID format)"
//	@Param		request		body		models.NdsActivityInstanceInput	true	"Instance data"
//	@Success	200			{object}	res.IdCreated
//	@Failure	400			{object}	res.Error
//	@Failure	401			{object}	res.Error
//	@Failure	403			{object}	res.Error	"Not enough permissions"
//	@Security	ApiKeyAuth
//	@Router		/nds/activities/{activityId}/instances [post]
func CreateInstance(c *fiber.Ctx) error {
	// parse request
	b := new(models.NdsActivityInstanceInput)
	parseErr := c.BodyParser(b)
	activityId := c.Params(ParamActivityId)
	user := helpers.GetUser(c)

	data, err := createInstance(&user, activityId, b, parseErr)
	return helpers.ProcessJsonResponse(c, data, err)
}

// DeleteInstance deletes an instance along with attendances for this instance.
//
//	@Summary	Delete an instance along with attendances for this instance.
//	@Tags		NDS Instance
//	@Produce	json
//	@Param		activityId	path	string	true	"Activity ID (UUID format)"
//	@Param		instanceId	path	string	true	"Instance ID (UUID format)"
//	@Success	200
//	@Failure	400	{object}	res.Error
//	@Failure	401	{object}	res.Error
//	@Failure	403	{object}	res.Error	"Not enough permissions"
//	@Security	ApiKeyAuth
//	@Router		/nds/activities/{activityId}/instances/{instanceId} [delete]
func DeleteInstance(c *fiber.Ctx) error {
	// parse request
	activityId := c.Params(ParamActivityId)
	instanceId := c.Params(ParamInstanceId)
	user := helpers.GetUser(c)

	err := deleteInstance(&user, activityId, instanceId)
	return helpers.ProcessJsonResponse(c, nil, err)
}

// GetParticipants lists all attendees for a specific activity instance.
//
//	@Summary	List all attendees for a specific activity instance.
//	@Tags		NDS Instance
//	@Produce	json
//	@Param		activityId	path		string	true	"Activity ID (UUID format)"
//	@Param		instanceId	path		string	true	"Instance ID (UUID format)"
//	@Success	200			{object}	[]models.NdsAttendance
//	@Failure	400			{object}	res.Error
//	@Failure	401			{object}	res.Error
//	@Failure	403			{object}	res.Error	"Not enough permissions"
//	@Security	ApiKeyAuth
//	@Router		/nds/activities/{activityId}/instances/{instanceId}/attendances [get]
func GetParticipants(c *fiber.Ctx) error {
	// parse request
	activityId := c.Params(ParamActivityId)
	instanceId := c.Params(ParamInstanceId)
	user := helpers.GetUser(c)

	data, err := selectAttendances(&user, activityId, instanceId)
	return helpers.ProcessJsonResponse(c, data, err)
}

// SetParticipants updates the list of attendees by first deleting the existing list, and then adding
// the given list.
//
//	@Summary	Update the list of attendees by first deleting the existing list, and then adding the given list.
//	@Tags		NDS Instance
//	@Consumes	json
//	@Produce	json
//	@Param		activityId	path	string						true	"Activity ID (UUID format)"
//	@Param		instanceId	path	string						true	"Instance ID (UUID format)"
//	@Param		request		body	[]models.NdsAttendanceInput	true	"List of attendees"
//	@Success	200
//	@Failure	400	{object}	res.Error
//	@Failure	401	{object}	res.Error
//	@Failure	403	{object}	res.Error	"Not enough permissions"
//	@Security	ApiKeyAuth
//	@Router		/nds/activities/{activityId}/instances/{instanceId}/attendances [put]
func SetParticipants(c *fiber.Ctx) error {
	// parse request
	b := new([]models.NdsAttendanceInput)
	parseErr := c.BodyParser(b)

	activityId := c.Params(ParamActivityId)
	instanceId := c.Params(ParamInstanceId)
	user := helpers.GetUser(c)

	err := setAttendances(&user, activityId, instanceId, b, parseErr)
	return helpers.ProcessJsonResponse(c, nil, err)
}

// ExportInstances exports the instances in CSV format ready for import into the NDS database.
//
//	@Summary	Export the instances in CSV format ready for import into the NDS database.
//	@Tags		NDS Export
//	@Produce	plain
//	@Param		activityId	path		string	true	"Activity ID (UUID format)"
//	@Success	200			{string}	string	"OK"
//	@Failure	400			{object}	res.Error
//	@Failure	401			{object}	res.Error
//	@Failure	403			{object}	res.Error	"Not enough permissions"
//	@Security	ApiKeyAuth
//	@Router		/nds/activities/{activityId}/export/instances [get]
func ExportInstances(c *fiber.Ctx) error {
	// parse request
	activityId := c.Params(ParamActivityId)
	user := helpers.GetUser(c)

	csv, err := exportInstances(&user, activityId)
	return helpers.ProcessCsvResponse(c, csv, err)
}

// ExportAttendances exports the attendances in CSV format ready for import into the NDS database.
//
//	@Summary	Export the attendances in CSV format ready for import into the NDS database.
//	@Tags		NDS Export
//	@Produce	plain
//	@Param		activityId	path		string	true	"Activity ID (UUID format)"
//	@Success	200			{string}	string	"OK"
//	@Failure	400			{object}	res.Error
//	@Failure	401			{object}	res.Error
//	@Failure	403			{object}	res.Error	"Not enough permissions"
//	@Security	ApiKeyAuth
//	@Router		/nds/activities/{activityId}/export/attendances [get]
func ExportAttendances(c *fiber.Ctx) error {
	// parse request
	activityId := c.Params(ParamActivityId)
	user := helpers.GetUser(c)

	csv, err := exportAttendances(&user, activityId)
	return helpers.ProcessCsvResponse(c, csv, err)
}

// ExportUsers exports the users in CSV format ready for import into the NDS database.
//
//	@Summary	Export the users in CSV format ready for import into the NDS database.
//	@Tags		NDS Export
//	@Produce	plain
//	@Param		activityId	path		string	true	"Activity ID  (UUID format)"
//	@Success	200			{string}	string	"OK"
//	@Failure	400			{object}	res.Error
//	@Failure	401			{object}	res.Error
//	@Failure	403			{object}	res.Error	"Not enough permissions"
//	@Security	ApiKeyAuth
//	@Router		/nds/activities/{activityId}/export/users [get]
func ExportUsers(c *fiber.Ctx) error {
	// parse request
	activityId := c.Params(ParamActivityId)
	user := helpers.GetUser(c)

	csv, err := exportUsers(&user, activityId)
	return helpers.ProcessCsvResponse(c, csv, err)
}
