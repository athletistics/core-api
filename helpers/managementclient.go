// Brief   : Connection towards the Zitadel Management API
// Author  : Nicolas Baeriswyl <nicolasfabrice.baeriswyl@edu.hefr.ch>
// Created : 20.06.2023

package helpers

import (
	"github.com/zitadel/oidc/pkg/oidc"
	"github.com/zitadel/zitadel-go/v2/pkg/client/management"
	"github.com/zitadel/zitadel-go/v2/pkg/client/middleware"
	"github.com/zitadel/zitadel-go/v2/pkg/client/user/v2beta"
	"github.com/zitadel/zitadel-go/v2/pkg/client/zitadel"
	"log"
	"os"
)

var MgmtClient *management.Client
var UserClient *v2beta.Client

func InitZitadelClients() {
	initMgmtClient()
	initUserClient()
}

func CloseZitadelClients() {
	closeMgmtClient()
	closeUserClient()
}

func initMgmtClient() {
	client, err := management.NewClient(
		os.Getenv("ZITADEL_INSTANCE"),
		os.Getenv("ZITADEL_INSTANCE_GRPC"),
		[]string{oidc.ScopeOpenID, zitadel.ScopeZitadelAPI()},
		zitadel.WithJWTProfileTokenSource(middleware.JWTProfileFromFileData([]byte(os.Getenv("ZITADEL_MGMT_API_KEY")))),
	)

	if err != nil {
		log.Fatalln("Unable to connect to zitadel", err)
	}

	MgmtClient = client
}

func closeMgmtClient() {
	defer func() {
		err := MgmtClient.Connection.Close()
		if err != nil {
			log.Println("Unable to close management grpc connection", err)
		}
	}()
}

func initUserClient() {
	client, err := v2beta.NewClient(
		os.Getenv("ZITADEL_INSTANCE"),
		os.Getenv("ZITADEL_INSTANCE_GRPC"),
		[]string{oidc.ScopeOpenID, zitadel.ScopeZitadelAPI()},
		zitadel.WithJWTProfileTokenSource(middleware.JWTProfileFromFileData([]byte(os.Getenv("ZITADEL_MGMT_API_KEY")))))

	if err != nil {
		log.Fatalln("Unable to connect to zitadel", err)
	}

	UserClient = client
}

func closeUserClient() {
	defer func() {
		err := UserClient.Connection.Close()
		if err != nil {
			log.Println("Unable to close user grpc connection", err)
		}
	}()
}
