// Brief   :
// Author  : Nicolas Baeriswyl <nicolasfabrice.baeriswyl@edu.hefr.ch>
// Created : 27.06.2023

package models

import "gorm.io/gorm"

type ExternalIdType struct {
	gorm.Model `json:"-"`
	ID         int    `json:"id" gorm:"primaryKey;unique;autoIncrement"`
	I18nName   string `json:"i18nName"`
}
