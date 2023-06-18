// Brief   :
// Author  : Nicolas Baeriswyl <nicolasfabrice.baeriswyl@edu.hefr.ch>
// Created : 18.06.2023

package models

import (
	"gorm.io/gorm"
	"time"
)

type Email struct {
	UserID    string `gorm:"primaryKey"`
	User      User
	Email     string `gorm:"primaryKey"`
	IsPrimary bool

	// gorm.Model without ID
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}
