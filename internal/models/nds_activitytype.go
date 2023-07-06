// Brief   :
// Author  : Nicolas Baeriswyl <nicolasfabrice.baeriswyl@edu.hefr.ch>
// Created : 06.07.2023

package models

import (
	"gorm.io/gorm"
)

type NdsActivityType struct {
	gorm.Model `json:"-"`
	ID         int    `json:"id" gorm:"primaryKey;unique;autoIncrement"`
	I18nName   string `json:"i18nName"`
	NdsName    string `json:"-"`
}
