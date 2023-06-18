// Brief   :
// Author  : Nicolas Baeriswyl <nicolasfabrice.baeriswyl@edu.hefr.ch>
// Created : 18.06.2023

package models

import (
	"gorm.io/gorm"
	"time"
)

type Phone struct {
	UserID      string `gorm:"primaryKey"`
	User        User
	PhoneNumber string `gorm:"primaryKey"`
	OwnerName   string
	IsEmergency bool

	// gorm.Model without ID
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}
