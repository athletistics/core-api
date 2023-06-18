// Brief   :
// Author  : Nicolas Baeriswyl <nicolasfabrice.baeriswyl@edu.hefr.ch>
// Created : 18.06.2023

package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Category struct {
	gorm.Model
	Id     uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey;unique"`
	Name   string
	Gender string `gorm:"type:varchar(1)"`
}
