// Author  : Nicolas Baeriswyl <nicolasfabrice.baeriswyl@edu.hefr.ch>
// Created : 08.07.2023

package types

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestDate_MarshalJSON(t *testing.T) {
	d := NewDate(2000, 02, 18)

	marshaled, _ := d.MarshalJSON()
	assert.Equal(t, `"2000-02-18"`, string(marshaled))
}

func TestDate_UnmarshalJSON(t *testing.T) {
	dString := `"2000-02-18"`
	d := Date{}

	err := d.UnmarshalJSON([]byte(dString))

	assert.Nil(t, err)
	assert.Equal(t, NewDate(2000, 2, 18).T(), d.T())
}
