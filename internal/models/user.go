// Brief   :
// Author  : Nicolas Baeriswyl <nicolasfabrice.baeriswyl@edu.hefr.ch>
// Created : 18.06.2023

package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

type User struct {
	gorm.Model
	Id            string    `gorm:"primaryKey;unique"`
	ClubID        uuid.UUID `gorm:"type:uuid"`
	Club          Club
	Oasi          *string
	BirthDate     time.Time
	Address       string
	PlaceID       int
	Place         Place
	NationalityID int
	Nationality   Country
	Teams         []Team `gorm:"many2many:participants;"`
}
