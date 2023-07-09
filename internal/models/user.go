// Brief   :
// Author  : Nicolas Baeriswyl <nicolasfabrice.baeriswyl@edu.hefr.ch>
// Created : 18.06.2023

package models

import (
	"fmt"
	"github.com/google/uuid"
	"gitlab.forge.hefr.ch/tb-athletistics/backend/database/types"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model    `json:"-"`
	ID            string     `gorm:"primaryKey;unique" json:"id"`
	FirstName     string     `json:"firstName" validate:"required"`
	LastName      string     `json:"lastName" validate:"required"`
	PrimaryMail   string     `json:"primaryMail" validate:"email"`
	PrimaryPhone  string     `json:"primaryPhone" validate:"required"`
	LanguageID    int        `json:"languageId" validate:"required"`
	Language      Language   `json:"-"`
	GenderID      int        `json:"genderId" validate:"required"`
	Gender        Gender     `json:"-"`
	ClubID        uuid.UUID  `gorm:"type:uuid" json:"clubId"`
	Club          Club       `json:"club"`
	Oasi          *string    `json:"oasi"`
	BirthDate     types.Date `json:"birthDate" validate:"required"`
	Address       string     `json:"address"`
	HouseNumber   string     `json:"houseNumber"`
	PostCode      string     `json:"postCode" validate:"required"`
	City          string     `json:"city" validate:"required"`
	CountryID     int        `json:"countryId" validate:"required"`
	Country       Country    `json:"country"`
	NationalityID int        `json:"nationalityId" validate:"required"`
	Nationality   Country    `json:"nationality"`
	Teams         []Team     `gorm:"many2many:participants;" json:"-"`
	CreateLogin   bool       `json:"createLogin" gorm:"-" validate:"required"`
}

type UserShort struct {
	gorm.Model `json:"-"`
	ID         string     `gorm:"primaryKey;unique" json:"id"`
	FirstName  string     `json:"firstName" validate:"required"`
	LastName   string     `json:"lastName" validate:"required"`
	BirthDate  types.Date `json:"birthDate" validate:"required"`
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

// CSV utilities

const NdsUserCsvHeader = "PERSONENNUMMER;NAME;VORNAME;GEBURTSDATUM;GESCHLECHT;AHV_NR;PEID;NATIONALITAET;MUTTERSPRACHE;STRASSE;HAUSNUMMER;PLZ;ORT;LAND"

func (u *User) AsCsv(ndsUserIds map[string]string) string {
	id, ok := ndsUserIds[u.ID]
	if !ok {
		id = ""
	}
	oasi := ""
	if u.Oasi != nil {
		oasi = *u.Oasi
	}

	return fmt.Sprintf("%v;%v;%v;%v;%v;%v;%v;%v;%v;%v;%v;%v;%v;%v",
		id, u.LastName, u.FirstName, u.BirthDate.T().Format(types.NdsDateFormat),
		u.Gender.NdsName, oasi, "", u.Nationality.CodeIso3166, u.Language.Short,
		u.Address, "", u.PostCode, u.City, u.Country.CodeIso3166)
}
