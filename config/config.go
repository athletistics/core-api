// Brief   : A set of API configurations
// Author  : Nicolas Baeriswyl <nicolasfabrice.baeriswyl@edu.hefr.ch>
// Created : 18.06.2023

package config

import (
	"fmt"
	"log"
	"os"
	"strconv"
)

const (
	AppName    = "athletistics"
	ApiVersion = 1
)

var (
	Issuer           = getIssuer()
	TotpSecondsLimit = getTotpSecondsLimit()
)

func getIssuer() string {
	env := os.Getenv("ENV")

	if env == "production" {
		return AppName
	} else if env == "" {
		env = "development"
	}

	return fmt.Sprintf("%v-%v", AppName, env)
}

func getTotpSecondsLimit() int {
	l := os.Getenv("TOTP_SECONDS_LIMIT")
	limit, err := strconv.Atoi(l)
	if err != nil {
		log.Fatalln("unable to load env variables", err)
	}
	return limit
}
