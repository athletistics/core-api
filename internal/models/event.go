// Brief   :
// Author  : Nicolas Baeriswyl <nicolasfabrice.baeriswyl@edu.hefr.ch>
// Created : 18.06.2023

package models

import (
	"fmt"
	"github.com/athletistics/core-api/database/types"
	"github.com/google/uuid"
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

func (e *Event) GetOwningClub(db *gorm.DB) (uuid.UUID, error) {
	if err := db.First(e).Error; err != nil {
		return uuid.Nil, err
	}

	if e.ClubID == nil {
		return uuid.Nil, fmt.Errorf(`event "%s" has no owning club`, e.ID)
	}

	return *e.ClubID, nil
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
