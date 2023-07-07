// Brief   :
// Author  : Nicolas Baeriswyl <nicolasfabrice.baeriswyl@edu.hefr.ch>
// Created : 06.07.2023

package models

import (
	"fmt"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

type NdsAttendance struct {
	UserID                string `gorm:"primaryKey"`
	NdsActivityInstance   NdsActivityInstance
	NdsActivityInstanceID uuid.UUID `gorm:"type:uuid;primaryKey"`
	IsCoach               bool

	CreatedAt time.Time      `json:"-"`
	DeletedAt gorm.DeletedAt `json:"-"`
}

type NdsAttendanceInput struct {
	UserID  string `gorm:"primaryKey"`
	IsCoach bool
}

// CSV utilities

const NdsAttendanceCsvHeader = "PERSONENNUMMER;FUNKTION;DATUM;AKTIVITÄTSTYP;ZEIT;DAUER;ORT"
const NdsAttendanceRoleCoach = "Teilnehmer/in"
const NdsAttendanceRoleAthlete = "Leiter/in"

func (a *NdsAttendance) AsCsv(ndsUserIds map[string]string) string {
	var role string
	if a.IsCoach {
		role = NdsAttendanceRoleCoach
	} else {
		role = NdsAttendanceRoleAthlete
	}

	id, ok := ndsUserIds[a.UserID]
	if !ok {
		return ""
	}

	return fmt.Sprintf("%v,%v,%v,%v,%v,%v,%v",
		id, role,
		a.NdsActivityInstance.Date.ForNds(),
		a.NdsActivityInstance.NdsActivityType.NdsName,
		a.NdsActivityInstance.StartTime.ForNds(),
		a.NdsActivityInstance.Duration,
		a.NdsActivityInstance.Place)
}
