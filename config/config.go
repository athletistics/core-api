// Brief   : A set of API configurations
// Author  : Nicolas Baeriswyl <nicolasfabrice.baeriswyl@edu.hefr.ch>
// Created : 18.06.2023

package config

import (
	"log"
	"os"
	"strconv"
)

const (
	AppName                 = "athletistics"
	ApiVersion              = 1
	defaultTotpSecondsLimit = 100
)

var (
	TotpSecondsLimit = getTotpSecondsLimit()
)

func getTotpSecondsLimit() int {
	l := os.Getenv("TOTP_SECONDS_LIMIT")
	limit, err := strconv.Atoi(l)
	if err != nil || limit < 0 {
		log.Println("unable to load env variables", err)
		return defaultTotpSecondsLimit
	}
	return limit
}
