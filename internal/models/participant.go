// Brief   :
// Author  : Nicolas Baeriswyl <nicolasfabrice.baeriswyl@edu.hefr.ch>
// Created : 18.06.2023

package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

type Participant struct {
	UserID    string    `gorm:"primaryKey"`
	TeamID    uuid.UUID `gorm:"type:uuid;primaryKey"`
	IsCoach   bool
	CreatedAt time.Time
	DeletedAt gorm.DeletedAt
}

type ParticipantInOut struct {
	UserID  string `json:"userId"`
	IsCoach bool   `json:"isCoach"`
}
