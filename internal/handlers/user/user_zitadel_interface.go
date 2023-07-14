// Brief   : Interface towards the Zitadel instance for managing users
// Author  : Nicolas Baeriswyl <nicolasfabrice.baeriswyl@edu.hefr.ch>
// Created : 20.06.2023

package userHandler

import (
	"context"
	"fmt"
	pb "github.com/zitadel/zitadel-go/v2/pkg/client/zitadel/management"
	"github.com/zitadel/zitadel-go/v2/pkg/client/zitadel/user"
	"gitlab.forge.hefr.ch/tb-athletistics/backend/constants"
	"gitlab.forge.hefr.ch/tb-athletistics/backend/helpers"
	"gitlab.forge.hefr.ch/tb-athletistics/backend/internal/models"
	"math/rand"
	"strings"
)

const UserStateNoLogin = "USER_STATE_NO_LOGIN"

func createUser(u *models.User) (string, error) {
	userNameOk := false
	var userName string

	cleanFirstName := strings.Join(strings.Fields(u.FirstName), "")
	cleanLastName := strings.Join(strings.Fields(u.LastName), "")

	for !userNameOk {
		userName = fmt.Sprintf("%v%v_%v", cleanFirstName, cleanLastName, getRandStr(constants.UserSuffixLength))
		req := pb.IsUserUniqueRequest{UserName: userName}

		res, err := helpers.MgmtClient.IsUserUnique(context.Background(), &req)
		if err != nil {
			return "", err
		}

		userNameOk = res.IsUnique
	}

	// Create the user
	createReq := pb.ImportHumanUserRequest{
		UserName: userName,
		Profile: &pb.ImportHumanUserRequest_Profile{
			FirstName:         u.FirstName,
			LastName:          u.LastName,
			PreferredLanguage: u.Language.Short,
			Gender:            getGenderFromString(u.Gender.ZitadelName),
		},
		Email: &pb.ImportHumanUserRequest_Email{
			Email:           u.PrimaryMail,
			IsEmailVerified: false,
		},
		PasswordChangeRequired: true,
	}

	createRes, err := helpers.MgmtClient.ImportHumanUser(context.Background(), &createReq)

	if err != nil {
		return "", err
	}

	return createRes.UserId, nil
}

func GetUserState(id string) (string, error) {
	req := pb.GetUserByIDRequest{Id: id}
	res, err := helpers.MgmtClient.GetUserByID(context.Background(), &req)

	if err != nil {
		return "", err
	}

	return res.User.State.String(), nil
}

func getGenderFromString(gender string) user.Gender {
	switch gender {
	case user.Gender_GENDER_UNSPECIFIED.String():
		return user.Gender_GENDER_UNSPECIFIED
	case user.Gender_GENDER_FEMALE.String():
		return user.Gender_GENDER_FEMALE
	case user.Gender_GENDER_MALE.String():
		return user.Gender_GENDER_MALE
	default:
		return user.Gender_GENDER_DIVERSE
	}
}

// source: https://golangdocs.com/generate-random-string-in-golang
func getRandStr(n int) string {
	var letters = []rune("ABCDEFGHIJKLMNPQRSTUVWXYZ123456789")

	s := make([]rune, n)
	for i := range s {
		s[i] = letters[rand.Intn(len(letters))]
	}
	return string(s)
}
