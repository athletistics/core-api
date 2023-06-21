// Brief   :
// Author  : Nicolas Baeriswyl <nicolasfabrice.baeriswyl@edu.hefr.ch>
// Created : 20.06.2023

package helpers

import (
	"fmt"
	"github.com/zitadel/oidc/pkg/oidc"
	"github.com/zitadel/zitadel-go/v2/pkg/client/management"
	"github.com/zitadel/zitadel-go/v2/pkg/client/middleware"
	"github.com/zitadel/zitadel-go/v2/pkg/client/zitadel"
	"log"
	"os"
)

var MgmtClient *management.Client

func InitMgmtClient() {
	// todo replace (env variable is empty for some reason)
	fmt.Println(os.Getenv("ZITADEL_MGMT_KEY_PATH"))
	client, err := management.NewClient(
		os.Getenv("ZITADEL_INSTANCE"),
		os.Getenv("ZITADEL_INSTANCE_GRPC"),
		[]string{oidc.ScopeOpenID, zitadel.ScopeZitadelAPI()},
		zitadel.WithJWTProfileTokenSource(middleware.JWTProfileFromPath("apikey-mgmt.json")),
	)

	if err != nil {
		log.Fatalln("Unable to connect to zitadel", err)
	}

	MgmtClient = client
}

func CloseMgmtClient() {
	defer func() {
		err := MgmtClient.Connection.Close()
		if err != nil {
			log.Println("Unable to close grpc connection", err)
		}
	}()
}
