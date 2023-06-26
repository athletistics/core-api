// Brief   :
// Author  : Nicolas Baeriswyl <nicolasfabrice.baeriswyl@edu.hefr.ch>
// Created : 18.06.2023

package models

import "gorm.io/gorm"

type Place struct {
	gorm.Model `json:"-"`
	ID         int `gorm:"primaryKey;unique;autoIncrement"`
	PostCode   string
	Place      string
	StateID    int
	State      State
}

type PlaceOut struct {
	ID        int    `json:"id"`
	PostCode  string `json:"postCode"`
	Place     string `json:"place"`
	StateID   int    `json:"stateId"`
	CountryID int    `json:"countryId"`
}
