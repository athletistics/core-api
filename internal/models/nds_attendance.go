// Brief   :
// Author  : Nicolas Baeriswyl <nicolasfabrice.baeriswyl@edu.hefr.ch>
// Created : 06.07.2023

package models

import (
	"fmt"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"strconv"
	"time"
)

type NdsAttendance struct {
	UserID                string              `json:"userId" gorm:"primaryKey"`
	NdsActivityInstance   NdsActivityInstance `json:"-"`
	NdsActivityInstanceID uuid.UUID           `json:"ndsActivityInstanceId" gorm:"type:uuid;primaryKey"`
	IsCoach               bool                `json:"isCoach"`

	CreatedAt time.Time      `json:"-"`
	DeletedAt gorm.DeletedAt `json:"-"`
}

type NdsAttendanceInput struct {
	UserID  string `json:"userId" gorm:"primaryKey"`
	IsCoach bool   `json:"isCoach"`
}

// CSV utilities

const NdsAttendanceCsvHeader = "PERSONENNUMMER;FUNKTION;DATUM;AKTIVITAETSTYP;ZEIT;DAUER;ORT"
const NdsAttendanceRoleAthlete = "Teilnehmer/in"
const NdsAttendanceRoleCoach = "Leiter/in"

func (a *NdsAttendance) AsCsv(ndsUserIds map[string]string) string {
	var role string
	if a.IsCoach {
		role = NdsAttendanceRoleCoach
	} else {
		role = NdsAttendanceRoleAthlete
	}

	t := ""
	if a.NdsActivityInstance.StartTime != nil {
		t = a.NdsActivityInstance.StartTime.ForNds()
	}

	d := ""
	if a.NdsActivityInstance.Duration != nil {
		d = strconv.Itoa(int(*a.NdsActivityInstance.Duration))
	}

	p := ""
	if a.NdsActivityInstance.Place != nil {
		p = *a.NdsActivityInstance.Place
	}

	id, ok := ndsUserIds[a.UserID]
	if !ok {
		return ""
	}

	return fmt.Sprintf("%v;%v;%v;%v;%v;%v;%v",
		id, role,
		a.NdsActivityInstance.Date.ForNds(),
		a.NdsActivityInstance.NdsActivityType.NdsName,
		t, d, p)
}
