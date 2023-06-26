// Brief   :
// Author  : Nicolas Baeriswyl <nicolasfabrice.baeriswyl@edu.hefr.ch>
// Created : 18.06.2023

package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Category struct {
	gorm.Model `json:"-"`
	ID         uuid.UUID `json:"id" gorm:"type:uuid;default:gen_random_uuid();primaryKey;unique"`
	Name       string    `json:"name"`
	Gender     string    `json:"gender" gorm:"type:varchar(1)"`
	SportID    int       `json:"-"`
	Sport      Sport     `json:"-"`
}
