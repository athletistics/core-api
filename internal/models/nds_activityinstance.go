// Brief   :
// Author  : Nicolas Baeriswyl <nicolasfabrice.baeriswyl@edu.hefr.ch>
// Created : 06.07.2023

package models

import (
	"fmt"
	"github.com/google/uuid"
	"gitlab.forge.hefr.ch/tb-athletistics/backend/database/types"
	"gorm.io/gorm"
	"strconv"
)

type NdsActivityInstance struct {
	gorm.Model        `json:"-"`
	ID                uuid.UUID       `json:"id" gorm:"type:uuid;default:gen_random_uuid();primaryKey;unique"`
	Date              types.Date      `json:"date"`
	StartTime         *types.Time     `json:"startTime"`
	Duration          *uint           `json:"duration"`
	Place             *string         `json:"place"`
	NdsActivity       NdsActivity     `json:"-"`
	NdsActivityID     uuid.UUID       `json:"-" gorm:"type:uuid"`
	NdsActivityType   NdsActivityType `json:"ndsActivityType"`
	NdsActivityTypeID int             `json:"-"`

	Users []User `json:"-" gorm:"many2many:nds_attendances;"`
}

// INPUT FORMATS used for decoding bodies

type NdsActivityInstanceInput struct {
	Date              types.Date  `json:"date"  validate:"required"`
	StartTime         *types.Time `json:"startTime"`
	Duration          *uint       `json:"duration"`
	Place             *string     `json:"place"`
	NdsActivityID     uuid.UUID   `json:"ndsActivityId" gorm:"type:uuid"  validate:"required"`
	NdsActivityTypeID int         `json:"ndsActivityTypeId" validate:"required"`
}

// CSV utilities

const NdsActivityInstanceCsvHeader = "AKTIVITAETSTYP;DATUM;ZEIT;DAUER;ORT;FOKUS"
const defaultFocus = "Physis"

func (i *NdsActivityInstance) AsCsv() string {
	t := ""
	if i.StartTime != nil {
		t = i.StartTime.ForNds()
	}

	d := ""
	if i.Duration != nil {
		d = strconv.Itoa(int(*i.Duration))
	}

	p := ""
	if i.Place != nil {
		p = *i.Place
	}

	return fmt.Sprintf("%v;%v;%v;%v;%v;%v", i.NdsActivityType.NdsName, i.Date.ForNds(), t, d, p, defaultFocus)
}
