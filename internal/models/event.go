// Brief   :
// Author  : Nicolas Baeriswyl <nicolasfabrice.baeriswyl@edu.hefr.ch>
// Created : 18.06.2023

package models

import (
	"github.com/google/uuid"
	"gitlab.forge.hefr.ch/tb-athletistics/backend/database/types"
	"gorm.io/gorm"
)

type Event struct {
	gorm.Model    `json:"-"`
	ID            uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey;unique"`
	EventTypeID   int
	EventType     EventType
	Name          string
	Place         string
	ClubID        *uuid.UUID `gorm:"type:uuid"`
	Club          *Club
	Date          types.Date
	IsHomologated bool
	IsOutdoor     bool `gorm:"default:true"`
}

// OUTPUT FORMATS used for sending data to the client

// EventShort is used for listing events
type EventShort struct {
	ID            uuid.UUID  `json:"id"`
	EventTypeID   int        `json:"eventTypeId"`
	Name          string     `json:"name"`
	Place         string     `json:"place"`
	Date          types.Date `json:"date"`
	IsHomologated bool       `json:"isHomologated"`
	IsOutdoor     bool       `json:"isOutdoor"`
}

// INPUT FORMATS used for decoding bodies

type EventInput struct {
	EventTypeID   int        `json:"eventTypeId" validate:"required"`
	Name          string     `json:"name" validate:"required"`
	Place         string     `json:"place" validate:"required"`
	Date          types.Date `json:"date" validate:"required"`
	IsHomologated bool       `json:"isHomologated" validate:"required"`
	IsOutdoor     bool       `json:"isOutdoor" validate:"required"`
}
