// Brief   :
// Author  : Nicolas Baeriswyl <nicolasfabrice.baeriswyl@edu.hefr.ch>
// Created : 18.06.2023

package models

import "gorm.io/gorm"

type Country struct {
	gorm.Model  `json:"-"`
	ID          int    `json:"id" gorm:"primaryKey;unique;autoIncrement"`
	Code        string `json:"code" gorm:"type:varchar(3)"`
	CodeIso3166 string `json:"codeIso3166" gorm:"type:varchar(2)"`
	Name        string `json:"name"`
}
