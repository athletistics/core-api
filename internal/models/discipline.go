// Brief   :
// Author  : Nicolas Baeriswyl <nicolasfabrice.baeriswyl@edu.hefr.ch>
// Created : 18.06.2023

package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Discipline struct {
	gorm.Model   `json:"-"`
	ID           uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey;unique"`
	SportID      int
	Sport        Sport
	ClubID       *uuid.UUID `gorm:"type:uuid"`
	Club         *Club
	I18nName     *string
	CustomName   *string
	IsDescending bool
}

// OUTPUT FORMATS used for sending data to the client

// DisciplineShort is used to list disciplines
type DisciplineShort struct {
	ID           uuid.UUID  `json:"id"`
	I18nName     *string    `json:"i18nName"`
	CustomName   *string    `json:"customName"`
	ClubId       *uuid.UUID `gorm:"type:uuid" json:"-"`
	IsDescending bool       `json:"isDescending"`
	IsOfficial   *bool      `gorm:"-" json:"isOfficial"`
}

// INPUT FORMATS used for decoding bodies

// DisciplineInput is used to create or modify a discipline
type DisciplineInput struct {
	CustomName   string `json:"customName" validate:"required"`
	IsDescending bool   `json:"isDescending" validate:"required"`
}
