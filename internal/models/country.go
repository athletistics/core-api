// Brief   :
// Author  : Nicolas Baeriswyl <nicolasfabrice.baeriswyl@edu.hefr.ch>
// Created : 18.06.2023

package models

import "gorm.io/gorm"

type Country struct {
	gorm.Model
	ID   int    `gorm:"primaryKey;unique;autoIncrement"`
	Code string `gorm:"type:varchar(3)"`
	Name string
}
