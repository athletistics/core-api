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
	gorm.Model    `json:"-"`
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

// OUTPUT FORMATS used for sending data to the client

// EventShort is used for listing events
type EventShort struct {
	ID            uuid.UUID `json:"id"`
	EventTypeID   int       `json:"eventTypeId"`
	Name          string    `json:"name"`
	Place         string    `json:"place"`
	Date          time.Time `json:"date"`
	IsHomologated bool      `json:"isHomologated"`
}

// INPUT FORMATS used for decoding bodies

type EventInput struct {
	EventTypeID   int       `json:"eventTypeId" validate:"required"`
	Name          string    `json:"name" validate:"required"`
	Place         string    `json:"place" validate:"required"`
	Date          time.Time `json:"date" validate:"required"`
	IsHomologated bool      `json:"isHomologated" validate:"required"`
}
