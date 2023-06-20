// Brief   :
// Author  : Nicolas Baeriswyl <nicolasfabrice.baeriswyl@edu.hefr.ch>
// Created : 18.06.2023

package models

import "gorm.io/gorm"

type Sport struct {
	gorm.Model `json:"-"`
	ID         int `gorm:"primaryKey;unique;autoIncrement"`
	I18nName   string
}
