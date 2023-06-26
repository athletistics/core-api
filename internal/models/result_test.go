// Brief   :
// Author  : Nicolas Baeriswyl <nicolasfabrice.baeriswyl@edu.hefr.ch>
// Created : 26.06.2023

package models

import (
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"testing"
)

func newUUID() *uuid.UUID {
	u, _ := uuid.NewUUID()
	return &u
}

func TestResultFromInput(t *testing.T) {
	r := ResultInput{
		UserId:       "1",
		CategoryID:   newUUID(),
		DisciplineID: *newUUID(),
		EventID:      *newUUID(),
	}

	rfi := ResultFromInput(&r)

	assert.Equal(t, r.CategoryID, rfi.CategoryID)
	assert.Equal(t, r.DisciplineID, rfi.DisciplineID)
	assert.Equal(t, r.EventID, rfi.EventID)
}

func TestResultString_Prepare_Discipline(t *testing.T) {
	i18n := "A"
	custom := "B"

	r := ResultString{}
	r.Prepare()
	assert.Empty(t, r.DisciplineName)

	r = ResultString{Discipline: Discipline{I18nName: &i18n}}
	r.Prepare()
	assert.Equal(t, i18n, r.DisciplineName)

	r = ResultString{Discipline: Discipline{CustomName: &custom}}
	r.Prepare()
	assert.Equal(t, custom, r.DisciplineName)

	r = ResultString{Discipline: Discipline{I18nName: &i18n, CustomName: &custom}}
	r.Prepare()
	assert.Equal(t, i18n, r.DisciplineName)
}

func TestResultString_Prepare_Category(t *testing.T) {
	r := ResultString{
		Category: nil,
	}
	r.Prepare()
	assert.Empty(t, r.CategoryName)

	r = ResultString{
		Category: &Category{
			Name:   "A",
			Gender: "B",
		},
	}
	r.Prepare()
	assert.Equal(t, "A B", r.CategoryName)
}
