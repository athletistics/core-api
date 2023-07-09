// Brief   :
// Author  : Nicolas Baeriswyl <nicolasfabrice.baeriswyl@edu.hefr.ch>
// Created : 27.06.2023

package models

import (
	"gorm.io/gorm"
	"time"
)

type ExternalId struct {
	ExternalID       string         `json:"externalId" gorm:"primaryKey"`
	InternalID       string         `json:"internalId" gorm:"primaryKey"`
	ExternalIdType   ExternalIdType `json:"-"`
	ExternalIdTypeID int            `json:"externalIdTypeId"`

	// gorm.Model without ID
	CreatedAt time.Time      `json:"-"`
	UpdatedAt time.Time      `json:"-"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}

type ExternalIdInput struct {
	ExternalID       string `json:"externalId"`
	InternalID       string `json:"internalId"`
	ExternalIdTypeID int    `json:"externalIdTypeId"`
}
