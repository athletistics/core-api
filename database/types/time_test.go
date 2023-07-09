// Author  : Nicolas Baeriswyl <nicolasfabrice.baeriswyl@edu.hefr.ch>
// Created : 08.07.2023

package types

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestTime_MarshalJSON(t *testing.T) {
	d := NewTime(12, 13, 14)

	marshaled, _ := d.MarshalJSON()
	assert.Equal(t, `"12:13"`, string(marshaled))
}

func TestTime_UnmarshalJSON(t *testing.T) {
	dString := `"12:13"`
	d := Time{}

	err := d.UnmarshalJSON([]byte(dString))

	assert.Nil(t, err)
	assert.Equal(t, NewTime(12, 13, 00).T(), d.T())
}
