// Brief   :
// Author  : Nicolas Baeriswyl <nicolasfabrice.baeriswyl@edu.hefr.ch>
// Created : 20.06.2023

package models

import "gorm.io/gorm"

type Gender struct {
	gorm.Model  `json:"-"`
	ID          int `gorm:"primaryKey;unique;autoIncrement"`
	ZitadelName string
	I18nName    string
}
