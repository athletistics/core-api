// Brief   : Interface towards the Zitadel instance for managing users
// Author  : Nicolas Baeriswyl <nicolasfabrice.baeriswyl@edu.hefr.ch>
// Created : 20.06.2023

package userHandler

import (
	"context"
	"fmt"
	"github.com/athletistics/core-api/constants"
	"github.com/athletistics/core-api/helpers"
	"github.com/athletistics/core-api/internal/models"
	res "github.com/athletistics/core-api/internal/responses"
	pb "github.com/zitadel/zitadel-go/v2/pkg/client/zitadel/management"
	"github.com/zitadel/zitadel-go/v2/pkg/client/zitadel/user/v2beta"
	pbv2 "github.com/zitadel/zitadel-go/v2/pkg/client/zitadel/user/v2beta"
	"math/rand"
	"net/http"
	"strings"
)

const UserStateNoLogin = "USER_STATE_NO_LOGIN"

func createUser(user *models.User) (string, error) {
	username, err := genUsername(user)
	if err != nil {
		return "", err
	}

	// Create the user
	gender := getGenderFromString(user.Gender.ZitadelName)
	createReq := pbv2.AddHumanUserRequest{
		Username: &username,
		Profile: &pbv2.SetHumanProfile{
			GivenName:         user.FirstName,
			FamilyName:        user.LastName,
			PreferredLanguage: &user.Language.Short,
			Gender:            &gender,
		},
		Email: &pbv2.SetHumanEmail{
			Email:        user.PrimaryMail,
			Verification: &pbv2.SetHumanEmail_IsVerified{IsVerified: false},
		},
		PasswordType: &pbv2.AddHumanUserRequest_Password{
			Password: &pbv2.Password{ChangeRequired: true},
		},
	}

	createRes, err := helpers.UserClient.AddHumanUser(context.Background(), &createReq)

	if err != nil {
		return "", err
	}

	return createRes.UserId, nil
}

func createUserLogin(user *models.User) (int, *res.Error) {
	// check if user already exists
	_, err := GetUserState(user.ID)
	if err != nil {
		return http.StatusBadRequest, res.ErrBadRequest
	}

	username, err := genUsername(user)
	if err != nil {
		return http.StatusInternalServerError, res.ErrZitadel.Err(err)
	}

	// Create the user
	gender := getGenderFromString(user.Gender.ZitadelName)
	createReq := pbv2.AddHumanUserRequest{
		UserId:   &user.ID,
		Username: &username,
		Profile: &pbv2.SetHumanProfile{
			GivenName:         user.FirstName,
			FamilyName:        user.LastName,
			PreferredLanguage: &user.Language.Short,
			Gender:            &gender,
		},
		Email: &pbv2.SetHumanEmail{
			Email:        user.PrimaryMail,
			Verification: &pbv2.SetHumanEmail_IsVerified{IsVerified: false},
		},
		PasswordType: &pbv2.AddHumanUserRequest_Password{
			Password: &pbv2.Password{ChangeRequired: true},
		},
	}

	_, err = helpers.UserClient.AddHumanUser(context.Background(), &createReq)

	if err != nil {
		return http.StatusInternalServerError, res.ErrZitadel.Err(err)
	}

	return http.StatusOK, nil
}

func GetUserState(id string) (string, error) {
	req := pb.GetUserByIDRequest{Id: id}
	response, err := helpers.MgmtClient.GetUserByID(context.Background(), &req)

	if err != nil {
		return "", err
	}

	return response.User.State.String(), nil
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

func genUsername(user *models.User) (string, error) {
	userNameOk := false
	var userName string

	cleanFirstName := strings.Join(strings.Fields(user.FirstName), "")
	cleanLastName := strings.Join(strings.Fields(user.LastName), "")

	for !userNameOk {
		userName = fmt.Sprintf("%v%v_%v", cleanFirstName, cleanLastName, getRandStr(constants.UserSuffixLength))
		req := pb.IsUserUniqueRequest{UserName: userName}

		response, err := helpers.MgmtClient.IsUserUnique(context.Background(), &req)
		if err != nil {
			return "", err
		}

		userNameOk = response.IsUnique
	}

	return userName, nil
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
