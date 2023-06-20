// Brief   :
// Author  : Nicolas Baeriswyl <nicolasfabrice.baeriswyl@edu.hefr.ch>
// Created : 18.06.2023

package models

import "gorm.io/gorm"

type State struct {
	gorm.Model
	ID        int `gorm:"primaryKey;unique;autoIncrement"`
	Name      string
	Short     string `gorm:"type:varchar(4)"`
	CountryID int
	Country   Country
}
