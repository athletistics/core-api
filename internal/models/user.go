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
	PrimaryPhone  string    `json:"primaryPhone" validate:"required"`
	LanguageID    int       `json:"languageId" validate:"required"`
	Language      Language  `json:"-"`
	GenderID      int       `json:"genderId" validate:"required"`
	Gender        Gender    `json:"-"`
	ClubID        uuid.UUID `gorm:"type:uuid" json:"clubId"`
	Club          Club      `json:"club"`
	Oasi          *string   `json:"oasi"`
	BirthDate     time.Time `json:"birthDate" validate:"required"`
	Address       string    `json:"address" validate:"required"`
	PlaceID       int       `json:"placeId" validate:"required"`
	Place         Place     `json:"place"`
	NationalityID int       `json:"nationalityId" validate:"required"`
	Nationality   Country   `json:"nationality"`
	Teams         []Team    `gorm:"many2many:participants;" json:"-"`
}

type UserShort struct {
	gorm.Model `json:"-"`
	ID         string    `gorm:"primaryKey;unique" json:"id"`
	FirstName  string    `json:"firstName" validate:"required"`
	LastName   string    `json:"lastName" validate:"required"`
	BirthDate  time.Time `json:"birthDate" validate:"required"`
}

type UserEmergency struct {
	gorm.Model   `json:"-"`
	ID           string `gorm:"primaryKey;unique" json:"id"`
	FirstName    string `json:"firstName" validate:"required"`
	LastName     string `json:"lastName" validate:"required"`
	PrimaryPhone string `json:"primaryPhone" validate:"required"`
}

type UserRole struct {
	Club string `json:"club"`
	Role string `json:"role"`
}
type UserResponse struct {
	AuthRoles []UserRole `json:"authRoles"`
	Is2FaOk   bool       `json:"is2FaOk"`
	Has2Fa    bool       `json:"has2Fa"`
}
