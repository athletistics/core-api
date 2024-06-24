// Brief   : Creating queries for reading results in the database
// Author  : Nicolas Baeriswyl <nicolasfabrice.baeriswyl@edu.hefr.ch>
// Created : 22.06.2023

package resultHandler

import (
	"fmt"
	"github.com/athletistics/core-api/database"
	"github.com/athletistics/core-api/helpers"
	"github.com/athletistics/core-api/internal/models"
	res "github.com/athletistics/core-api/internal/responses"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"net/http"
	"strconv"
)

type ResultQuery struct {
	UserID       *string    `query:"userId"`
	EventID      *uuid.UUID `query:"eventId"`
	EventTypeID  *int       `query:"eventTypeId"`
	CategoryId   *int       `query:"categoryId"`
	DisciplineID *uuid.UUID `query:"disciplineId"`
	BestOnly     bool       `query:"bestOnly"`
	helpers.Pagination
}

// FindResults searches for results in the database. This function does not perform any security checks.
// The calling function must make sure that the query can be executed as given.
// Only results of the requesting user's club are returned
func FindResults(query *ResultQuery, page *res.Page, outModel interface{}, callingUser *models.User) error {
	var err error
	db := database.DB

	// join results with users, limiting the results to the users that belong to the requesting user's club and users
	// that are shared with their club.
	q := db.Model(&models.Result{}).Select("results.id, results.created_at, results.updated_at, results.deleted_at, results.user_id, results.discipline_id, events.date as date, results.value, results.wind, results.rank, results.event_id, results.category_id, false is_personal_training_result").
		Joins("inner join users ON users.id = results.user_id AND users.id IN (?)",
			db.Distinct("users.id").Table("users").
				Joins(`LEFT JOIN authorisations a ON a.user_id = users.id`).
				Where(`users.club_id = ? OR a.club_id = ?`, callingUser.ClubID, callingUser.ClubID),
		).
		Joins("LEFT JOIN events ON events.id = results.event_id")
	qPers := db.Model(&models.ResultTraining{}).Select("result_trainings.*, NULL as wind, NULL as rank, NULL as event_id, NULL as  category_id, true is_personal_training_result")

	// apply filters passed through the query argument
	applyQueryFilters(query, q, qPers)

	// get best results by discipline
	if query.BestOnly {
		q.Select("results.*, disciplines.is_descending").Joins("inner join disciplines ON disciplines.id = results.discipline_id")

		descendingQ := db.Table("(?) as d", q).Select("*, ROW_NUMBER() OVER (PARTITION BY d.discipline_id ORDER BY d.value DESC) as rn")
		descendingQ = db.Raw("WITH ranked AS (?) SELECT * FROM ranked WHERE rn = 1 AND is_descending = TRUE", descendingQ)

		ascendingQ := db.Table("(?) as d", q).Select("*, ROW_NUMBER() OVER (PARTITION BY d.discipline_id ORDER BY d.value ASC) as rn")
		ascendingQ = db.Raw("WITH ranked AS (?) SELECT * FROM ranked WHERE rn = 1 AND is_descending = FALSE", ascendingQ)

		q = db.Raw("(?) UNION (?)", descendingQ, ascendingQ)
	} else {
		// Pagination
		var count int64 = 0
		var countPers int64 = 0
		q.Count(&count)
		qPers.Count(&countPers)

		page.Count = count + countPers

		if query.UserID != nil {
			q = db.Raw("(?) UNION (?)", q, qPers)
		}

		q.Scopes(helpers.Paginate(&query.Pagination, page))
	}

	q.Order("events.date desc")

	err = q.Preload("User").Preload("Discipline").Preload("Event").Preload("Category").Find(outModel).Error
	if err != nil {
		return err
	}

	return nil
}

// RunQuery starts a query with the ResultQuery parameter, and sends the Fiber response.
// This function must be called as a last statement of a Fiber handler, and the error
// should be returned without checking it.
func RunQuery(c *fiber.Ctx, q *ResultQuery) error {
	user := helpers.GetUser(c)

	var results res.Paginated[models.ResultString]
	if err := FindResults(q, &results.Page, &results.Data, &user); err != nil {
		helpers.LogError(err.Error(), helpers.GetUserId(c))
		return c.Status(http.StatusInternalServerError).JSON(res.ErrDB.Err(err))
	}

	for i := 0; i < len(results.Data); i++ {
		results.Data[i].Prepare()
	}
	if q.BestOnly {
		return c.JSON(results.Data)
	} else {
		return c.JSON(results)
	}
}

func RunCreateQuery(c *fiber.Ctx, body *models.ResultInput) error {
	user := helpers.GetUser(c)
	db := database.DB

	// make sure given UserID is numeric
	if _, err := strconv.ParseInt(body.UserId, 10, 64); err != nil {
		helpers.LogInfo(fmt.Sprintf("bad request: %v", err), helpers.GetUserId(c))
		return c.Status(http.StatusBadRequest).JSON(res.ErrBadRequest.Err(err))
	}

	// make sure user, event, and discipline can be used
	u := models.User{ID: body.UserId}
	if err := db.First(&u).Error; err != nil {
		helpers.LogSecurity(fmt.Sprintf("tried using user id %v", body.UserId), helpers.GetUserId(c))
		return c.Status(http.StatusForbidden).JSON(res.ErrForbidden)
	}

	if !helpers.HasRoleOnOneOfTheClubs(&user, &u, helpers.RoleCoach) {
		helpers.LogSecurity(fmt.Sprintf("tried using user id %v", body.UserId), helpers.GetUserId(c))
		return c.Status(http.StatusForbidden).JSON(res.ErrForbidden)
	}

	d := models.Discipline{ID: body.DisciplineID}
	if err := db.Where("club_id IS NULL OR club_id = ?", user.ClubID).First(&d).Error; err != nil {
		helpers.LogSecurity(fmt.Sprintf("tried using discipline %v", body.DisciplineID), helpers.GetUserId(c))
		return c.Status(http.StatusForbidden).JSON(res.ErrForbidden.Err(err))
	}

	e := models.Event{ID: body.EventID}
	if err := db.Where("club_id IS NULL OR club_id = ?", user.ClubID).First(&e).Error; err != nil {
		helpers.LogSecurity(fmt.Sprintf("tried accessing using event %v", body.EventID), helpers.GetUserId(c))
		return c.Status(http.StatusForbidden).JSON(res.ErrForbidden.Err(err))
	}

	// create the result
	result := models.ResultFromInput(body)
	if err := db.Create(&result).Error; err != nil {
		helpers.LogError(err.Error(), helpers.GetUserId(c))
		return c.Status(http.StatusInternalServerError).JSON(res.ErrDB.Err(err))
	}

	return c.JSON(res.UuidCreated{Id: result.ID})
}

func applyQueryFilters(query *ResultQuery, q *gorm.DB, qPers *gorm.DB) {
	if query.UserID != nil {
		q.Where(&models.Result{UserID: *query.UserID})
		qPers.Where(&models.Result{UserID: *query.UserID})
	}

	if query.EventID != nil {
		q.Where(&models.Result{EventID: *query.EventID})
	}

	if query.EventTypeID != nil {
		q.Joins("inner join events ON events.id = results.event_id AND events.event_type_id = ?", *query.EventTypeID)
	}

	if query.DisciplineID != nil {
		q.Where(&models.Result{DisciplineID: *query.DisciplineID})
		qPers.Where(&models.Result{DisciplineID: *query.DisciplineID})
	}
}
