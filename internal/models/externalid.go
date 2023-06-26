// Brief   :
// Author  : Nicolas Baeriswyl <nicolasfabrice.baeriswyl@edu.hefr.ch>
// Created : 27.06.2023

package models

import (
	"gorm.io/gorm"
	"time"
)

type ExternalId struct {
	ExternalID string `json:"externalId" gorm:"primaryKey"`
	InternalID string `json:"internalId" gorm:"primaryKey"`
	Short      string `json:"short"`
	I18nName   string `json:"i18nName"`

	// gorm.Model without ID
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}
