// Brief   :
// Author  : Nicolas Baeriswyl <nicolasfabrice.baeriswyl@edu.hefr.ch>
// Created : 16.06.2024

package models

import (
	"github.com/athletistics/core-api/database/types"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

type ResultTraining struct {
	gorm.Model   `json:"-"`
	ID           uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey;unique"`
	UserID       string
	User         User
	DisciplineID uuid.UUID `gorm:"type:uuid"`
	Discipline   Discipline
	Date         types.Date
	Value        decimal.Decimal `gorm:"type:decimal(8,3)"`
}

type ResultTrainingInput struct {
	DisciplineID uuid.UUID       `json:"disciplineId" validate:"required"`
	Date         types.Date      `json:"date" validate:"required"`
	Value        decimal.Decimal `json:"value" validate:"required"`
}

func (r *ResultTrainingInput) NewDBEntry(userId string) *ResultTraining {
	return &ResultTraining{
		UserID:       userId,
		DisciplineID: r.DisciplineID,
		Date:         r.Date,
		Value:        r.Value,
	}
}

func (r *ResultTrainingInput) UpdateDBEntry(dbEntry *ResultTraining) {
	dbEntry.DisciplineID = r.DisciplineID
	dbEntry.Date = r.Date
	dbEntry.Value = r.Value
}
