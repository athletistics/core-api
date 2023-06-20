// Brief   :
// Author  : Nicolas Baeriswyl <nicolasfabrice.baeriswyl@edu.hefr.ch>
// Created : 18.06.2023

package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Discipline struct {
	gorm.Model `json:"-"`
	ID           uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey;unique"`
	SportID      int
	Sport        Sport
	ClubID       *uuid.UUID `gorm:"type:uuid"`
	Club         *Club
	I18nName     *string
	CustomName   *string
	IsDescending bool
}
