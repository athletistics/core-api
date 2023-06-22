// Brief   :
// Author  : Nicolas Baeriswyl <nicolasfabrice.baeriswyl@edu.hefr.ch>
// Created : 21.06.2023

package disciplineHandler

import (
	"errors"
	"github.com/gofiber/fiber/v2"
	"gitlab.forge.hefr.ch/tb-athletistics/backend/database"
	"gitlab.forge.hefr.ch/tb-athletistics/backend/helpers"
	"gitlab.forge.hefr.ch/tb-athletistics/backend/internal/models"
	res "gitlab.forge.hefr.ch/tb-athletistics/backend/internal/responses"
	"gorm.io/gorm"
	"net/http"
)

func GetDisciplines(c *fiber.Ctx) error {
	if !helpers.HasAnyRoleOnParentClub(c) {
		return c.Status(http.StatusForbidden).JSON(res.ErrForbidden)
	}

	db := database.DB
	user := helpers.GetUser(c)
	var disciplines []models.DisciplineShort
	err := db.Model(&models.Discipline{}).
		Where("club_id IS NULL and sport_id = (?)", db.Table("clubs").Select("sport_id").Where("id = ?", user.ClubID)).
		Or(&models.Discipline{ClubID: &user.ClubID}).Find(&disciplines).Error

	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(res.ErrDB.Err(err))
	}

	for i := 0; i < len(disciplines); i++ {
		isOfficial := disciplines[i].ClubId == nil
		disciplines[i].IsOfficial = &isOfficial
	}

	return c.JSON(disciplines)
}

func CreateDisciplines(c *fiber.Ctx) error {
	if !helpers.HasRoleOnClub(c, helpers.GetUser(c).ClubID, helpers.RoleCoach) {
		return c.Status(http.StatusForbidden).JSON(res.ErrForbidden)
	}

	// read input body
	reqBody := new(models.DisciplineInput)
	if err := c.BodyParser(reqBody); err != nil {
		return c.Status(http.StatusBadRequest).JSON(res.ErrBadRequest.Err(err))
	}

	// get sport id of parent club
	db := database.DB
	club := models.Club{ID: helpers.GetUser(c).ClubID}
	if err := db.First(&club).Error; errors.Is(err, gorm.ErrRecordNotFound) {
		return c.Status(http.StatusNotFound).JSON(res.ErrNotFound)
	} else if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(res.ErrDB.Err(err))
	}

	// write to database
	discipline := models.Discipline{
		SportID:      club.SportID,
		ClubID:       &club.ID,
		CustomName:   &reqBody.CustomName,
		IsDescending: reqBody.IsDescending,
	}
	if err := db.Create(&discipline).Error; err != nil {
		return c.Status(http.StatusInternalServerError).JSON(res.ErrDB.Err(err))
	}

	// return id of new discipline
	return c.JSON(res.UuidCreated{Id: discipline.ID})
}
