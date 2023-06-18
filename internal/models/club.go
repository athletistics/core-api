// Brief   :
// Author  : Nicolas Baeriswyl <nicolasfabrice.baeriswyl@edu.hefr.ch>
// Created : 18.06.2023

package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Club struct {
	gorm.Model
	ID      uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey;unique"`
	SportID int
	Sport   Sport
	Name    string
}
