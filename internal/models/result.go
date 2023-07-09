// Brief   :
// Author  : Nicolas Baeriswyl <nicolasfabrice.baeriswyl@edu.hefr.ch>
// Created : 18.06.2023

package models

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"gitlab.forge.hefr.ch/tb-athletistics/backend/database/types"
	"gorm.io/gorm"
)

const TrainingTypeId = 4

type Result struct {
	gorm.Model   `json:"-"`
	ID           uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey;unique"`
	UserID       string
	User         User
	CategoryID   *uuid.UUID `gorm:"type:uuid"`
	Category     *Category
	DisciplineID uuid.UUID `gorm:"type:uuid"`
	Discipline   Discipline
	EventID      uuid.UUID `gorm:"type:uuid"`
	Event        Event
	Value        decimal.Decimal  `gorm:"type:decimal(8,3)"`
	Wind         *decimal.Decimal `gorm:"type:decimal(2,1)"`
	Rank         *string          `gorm:"type:varchar(4)"`
}

func ResultFromInput(input *ResultInput) *Result {
	return &Result{
		UserID:       input.UserId,
		CategoryID:   input.CategoryID,
		DisciplineID: input.DisciplineID,
		EventID:      input.EventID,
		Value:        input.Value,
		Wind:         input.Wind,
		Rank:         input.Rank,
	}
}

type ResultString struct {
	ID             uuid.UUID        `json:"id"`
	User           User             `json:"-"`
	UserID         string           `json:"-"`
	UserFirstName  string           `json:"userFirstName" gorm:"-"`
	UserLastName   string           `json:"userLastName" gorm:"-"`
	Category       *Category        `json:"-"`
	CategoryID     uuid.UUID        `gorm:"type:uuid" json:"-"`
	CategoryName   string           `json:"categoryName" gorm:"-"`
	Discipline     Discipline       `json:"-"`
	DisciplineID   uuid.UUID        `gorm:"type:uuid" json:"disciplineId"`
	DisciplineName string           `json:"disciplineName" gorm:"-"`
	Event          Event            `json:"-"`
	EventID        uuid.UUID        `gorm:"type:uuid" json:"-"`
	EventName      string           `json:"eventName" gorm:"-"`
	EventPlace     string           `json:"eventPlace" gorm:"-"`
	EventOutdoor   bool             `json:"eventOutdoor" gorm:"-"`
	IsTraining     bool             `json:"isTraining" gorm:"-"`
	Date           types.Date       `json:"date" gorm:"-"`
	Value          decimal.Decimal  `json:"value"`
	Wind           *decimal.Decimal `json:"wind"`
	Rank           *string          `json:"rank"`
}

func (r *ResultString) Prepare() {
	r.UserFirstName = r.User.FirstName
	r.UserLastName = r.User.LastName
	r.EventName = r.Event.Name
	r.EventPlace = r.Event.Place
	r.EventOutdoor = r.Event.IsOutdoor
	r.Date = r.Event.Date
	r.IsTraining = r.Event.EventTypeID == TrainingTypeId

	if r.Discipline.I18nName != nil {
		r.DisciplineName = *r.Discipline.I18nName
	} else if r.Discipline.CustomName != nil {
		r.DisciplineName = *r.Discipline.CustomName
	}

	if r.Category != nil {
		r.CategoryName = fmt.Sprintf("%v %v", r.Category.Name, r.Category.Gender)
	}
}

type ResultInput struct {
	UserId       string           `json:"userId"`
	CategoryID   *uuid.UUID       `json:"categoryId"`
	DisciplineID uuid.UUID        `json:"disciplineId" validate:"required"`
	EventID      uuid.UUID        `json:"eventId" validate:"required"`
	Value        decimal.Decimal  `json:"value" validate:"required"`
	Wind         *decimal.Decimal `json:"wind"`
	Rank         *string          `json:"rank"`
}
