// Brief   :
// Author  : Nicolas Baeriswyl <nicolasfabrice.baeriswyl@edu.hefr.ch>
// Created : 18.06.2023

package models

import "gorm.io/gorm"

const (
	AccessRoleIdAthlete   = 1
	AccessRoleIdCoach     = 2
	AccessRoleIdHeadCoach = 3
)

type AccessRole struct {
	gorm.Model `json:"-"`
	ID         int    `json:"id" gorm:"primaryKey;unique;autoIncrement"`
	I18nName   string `json:"i18nName"`
}
