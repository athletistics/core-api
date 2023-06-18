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
	AccessRoleID int `gorm:"primaryKey"`
	AccessRole   AccessRole
	UserID       string `gorm:"primaryKey"`
	User         User
	ClubID       uuid.UUID `gorm:"type:uuid;primaryKey"`
	Club         Club

	// gorm.Model without ID
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}
