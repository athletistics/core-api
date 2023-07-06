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
	UserID                string    `gorm:"primaryKey"`
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

const NdsAttendanceCsvHeader = "AKTIVITAETSTYP;DATUM;ZEIT;DAUER;ORT;FOKUS"
const NdsAttendanceRoleCoach = "Teilnehmer/in"
const NdsAttendanceRoleAthlete = "Leiter/in"

func (i *NdsAttendance) AsCsv() string {
	return fmt.Sprintf("%v,%v,%v,%v,%v,%v", i.NdsActivityType.NdsName, i.StartTime.ForNds(), i.StartTime.ForNds(), i.Duration, i.Place, defaultFocus)
}
