// Brief   :
// Author  : Nicolas Baeriswyl <nicolasfabrice.baeriswyl@edu.hefr.ch>
// Created : 18.06.2023

package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Team struct {
	gorm.Model `json:"-"`
	ID         uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey;unique" json:"id"`
	ClubID     uuid.UUID `gorm:"type:uuid" json:"clubId"`
	Club       Club      `json:"-"`
	Name       string    `json:"name"`
	Members    []User    `gorm:"many2many:participants;" json:"-"`
}
