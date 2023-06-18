// Brief   :
// Author  : Nicolas Baeriswyl <nicolasfabrice.baeriswyl@edu.hefr.ch>
// Created : 18.06.2023

package models

import "gorm.io/gorm"

type AccessRole struct {
	gorm.Model
	ID       int `gorm:"primaryKey;unique;autoIncrement"`
	I18nName string
}
