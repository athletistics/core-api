// Brief   :
// Author  : Nicolas Baeriswyl <nicolasfabrice.baeriswyl@edu.hefr.ch>
// Created : 18.06.2023

package models

import "gorm.io/gorm"

type Place struct {
	gorm.Model
	Id       int `gorm:"primaryKey;unique;autoIncrement"`
	PostCode string
	Place    string
	StateID  int
	State    State
}
