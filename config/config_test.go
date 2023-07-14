// Author  : Nicolas Baeriswyl <nicolasfabrice.baeriswyl@edu.hefr.ch>
// Created : 26.06.2023

package config

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestGetTotpSecondsLimit_Normal(t *testing.T) {
	inputs := []int{100, 1000, 100000000000}

	for i := 0; i < len(inputs); i++ {
		_ = os.Setenv("TOTP_SECONDS_LIMIT", fmt.Sprintf("%v", inputs[i]))
		limit := getTotpSecondsLimit()
		assert.Equal(t, inputs[i], limit)
	}
}

func TestGetTotpSecondsLimit_Empty(t *testing.T) {
	_ = os.Unsetenv("TOTP_SECONDS_LIMIT")

	assert.Equal(t, defaultTotpSecondsLimit, getTotpSecondsLimit())
}

func TestGetTotpSecondsLimit_Text(t *testing.T) {
	_ = os.Setenv("TOTP_SECONDS_LIMIT", fmt.Sprintf("%v", "sometext"))

	assert.Equal(t, defaultTotpSecondsLimit, getTotpSecondsLimit())
}

func TestGetTotpSecondsLimit_Negative(t *testing.T) {
	_ = os.Setenv("TOTP_SECONDS_LIMIT", fmt.Sprintf("%v", "-1000"))

	assert.Equal(t, defaultTotpSecondsLimit, getTotpSecondsLimit())
}
