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
	gorm.Model    `json:"-"`
	ID            string    `gorm:"primaryKey;unique" json:"id"`
	FirstName     string    `json:"firstName" validate:"required"`
	LastName      string    `json:"lastName" validate:"required"`
	PrimaryMail   string    `json:"primaryMail" validate:"required,email"`
	LanguageID    int       `json:"languageId" validate:"required"`
	Language      Language  `json:"-"`
	GenderID      int       `json:"genderId" validate:"required"`
	Gender        Gender    `json:"-"`
	ClubID        uuid.UUID `gorm:"type:uuid" json:"clubId"`
	Club          Club      `json:"-"`
	Oasi          *string   `json:"oasi"`
	BirthDate     time.Time `json:"birthDate" validate:"required"`
	Address       string    `json:"address" validate:"required"`
	PlaceID       int       `json:"placeId" validate:"required"`
	Place         Place     `json:""`
	NationalityID int       `json:"nationalityId" validate:"required"`
	Nationality   Country   `json:"-"`
	Teams         []Team    `gorm:"many2many:participants;" json:"-"`
}

type ShortUser struct {
	gorm.Model `json:"-"`
	ID         string    `gorm:"primaryKey;unique" json:"id"`
	FirstName  string    `json:"firstName" validate:"required"`
	LastName   string    `json:"lastName" validate:"required"`
	BirthDate  time.Time `json:"birthDate" validate:"required"`
}

type EmergencyUser struct {
	gorm.Model `json:"-"`
	ID         string `gorm:"primaryKey;unique" json:"id"`
	FirstName  string `json:"firstName" validate:"required"`
	LastName   string `json:"lastName" validate:"required"`
}
