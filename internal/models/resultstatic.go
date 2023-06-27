// Brief   :
// Author  : Nicolas Baeriswyl <nicolasfabrice.baeriswyl@edu.hefr.ch>
// Created : 27.06.2023

package models

import (
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

type ResultStatic struct {
	gorm.Model         `json:"-"`
	ID                 int              `json:"id" gorm:"primaryKey;unique;autoIncrement"`
	Value              decimal.Decimal  `json:"value" gorm:"type:decimal(8,3)"`
	CategoryID         uuid.UUID        `json:"categoryId" gorm:"type:uuid"`
	Category           Category         `json:"-"`
	DisciplineID       uuid.UUID        `json:"disciplineId" gorm:"type:uuid"`
	Discipline         Discipline       `json:"-"`
	ResultStaticTypeID int              `json:"resultStaticTypeId" gorm:"type:uuid"`
	ResultStaticType   ResultStaticType `json:"-"`
}

type ResultStaticType struct {
	gorm.Model `json:"-"`
	ID         int    `json:"id" gorm:"primaryKey;unique;autoIncrement"`
	I18nType   string `json:"i18nType"`
}

type ResultStaticOut struct {
	Results []ResultStatic     `json:"results"`
	Types   []ResultStaticType `json:"types"`
}
