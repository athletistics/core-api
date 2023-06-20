// Brief   :
// Author  : Nicolas Baeriswyl <nicolasfabrice.baeriswyl@edu.hefr.ch>
// Created : 18.06.2023

package models

import (
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

type Result struct {
	gorm.Model
	ID           uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey;unique"`
	UserID       string
	User         User
	CategoryID   uuid.UUID `gorm:"type:uuid"`
	Category     *Category
	DisciplineID uuid.UUID `gorm:"type:uuid"`
	Discipline   Discipline
	EventID      uuid.UUID `gorm:"type:uuid"`
	Event        Event
	Value        decimal.Decimal  `gorm:"type:decimal(8,3)"`
	Wind         *decimal.Decimal `gorm:"type:decimal(2,3)"`
	Rank         *string          `gorm:"type:varchar(4)"`
}
