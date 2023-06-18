// Brief   :
// Author  : Nicolas Baeriswyl <nicolasfabrice.baeriswyl@edu.hefr.ch>
// Created : 18.06.2023

package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Team struct {
	gorm.Model
	Id      uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey;unique"`
	ClubID  uuid.UUID `gorm:"type:uuid"`
	Club    Club
	Name    string
	Members []User `gorm:"many2many:participants;"`
}
