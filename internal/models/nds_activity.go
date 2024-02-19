// Brief   :
// Author  : Nicolas Baeriswyl <nicolasfabrice.baeriswyl@edu.hefr.ch>
// Created : 06.07.2023

package models

import (
	"github.com/athletistics/core-api/database/types"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

type NdsActivity struct {
	gorm.Model `json:"-"`
	ID         uuid.UUID `json:"id" gorm:"type:uuid;default:gen_random_uuid();primaryKey;unique"`

	StartDate types.Date   `json:"startDate"`
	EndDate   types.Date   `json:"endDate"`
	StartTime types.Time   `json:"startTime"`
	Duration  uint         `json:"duration"`
	Weekday   time.Weekday `json:"weekday" swaggertype:"integer"`
	Place     string       `json:"place"`
	Team      Team         `json:"team"`
	TeamID    uuid.UUID    `json:"-" gorm:"type:uuid"`
}

type NdsActivityInput struct {
	StartDate types.Date   `json:"startDate" validate:"required"`
	EndDate   types.Date   `json:"endDate" validate:"required"`
	StartTime types.Time   `json:"startTime" validate:"required"`
	Duration  uint         `json:"duration" validate:"required"`
	Weekday   time.Weekday `json:"weekday" validate:"required" swaggertype:"integer"`
	Place     string       `json:"place" validate:"required"`
	TeamID    uuid.UUID    `json:"teamId" gorm:"type:uuid" validate:"required"`
}
