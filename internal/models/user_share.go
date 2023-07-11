// Brief   :
// Author  : Nicolas Baeriswyl <nicolasfabrice.baeriswyl@edu.hefr.ch>
// Created : 11.07.2023

package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserShare struct {
	gorm.Model `json:"-"`
	ID         uuid.UUID `json:"code" gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	UserID     string    `json:"-"`
	User       User      `json:"-"`
}
