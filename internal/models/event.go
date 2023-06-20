// Brief   :
// Author  : Nicolas Baeriswyl <nicolasfabrice.baeriswyl@edu.hefr.ch>
// Created : 18.06.2023

package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

type Event struct {
	gorm.Model
	ID            uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey;unique"`
	EventTypeID   int
	EventType     EventType
	Name          string
	Place         string
	ClubID        *uuid.UUID `gorm:"type:uuid"`
	Club          *Club
	Date          time.Time
	IsHomologated bool
}
