// Brief   :
// Author  : Nicolas Baeriswyl <nicolasfabrice.baeriswyl@edu.hefr.ch>
// Created : 18.06.2023

package models

import "gorm.io/gorm"

type State struct {
	gorm.Model `json:"-"`
	ID         int     `json:"id" gorm:"primaryKey;unique;autoIncrement"`
	Name       string  `json:"name"`
	Short      string  `json:"short" gorm:"type:varchar(4)"`
	CountryID  int     `json:"countryId"`
	Country    Country `json:"-"`
}
