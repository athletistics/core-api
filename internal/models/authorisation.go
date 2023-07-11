// Brief   :
// Author  : Nicolas Baeriswyl <nicolasfabrice.baeriswyl@edu.hefr.ch>
// Created : 18.06.2023

package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

type Authorisation struct {
	AccessRoleID int        `json:"accessRoleId" gorm:"primaryKey"`
	AccessRole   AccessRole `json:"-"`
	UserID       string     `json:"userId" gorm:"primaryKey"`
	User         User       `json:"-"`
	ClubID       uuid.UUID  `json:"clubId" gorm:"type:uuid;primaryKey"`
	Club         Club       `json:"-"`

	// gorm.Model without ID
	CreatedAt time.Time      `json:"-"`
	UpdatedAt time.Time      `json:"-"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}
