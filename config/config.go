// Brief   : A set of API configurations
// Author  : Nicolas Baeriswyl <nicolasfabrice.baeriswyl@edu.hefr.ch>
// Created : 18.06.2023

package config

import (
	"fmt"
	"os"
)

const (
	AppName          = "athletistics"
	ApiVersion       = 1
	TotpSecondsLimit = 1000
)

var (
	Issuer = getIssuer()
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
