package auth

import (
	"encoding/json"
	"os"
)

type UserInformation struct {
	AuthToken string `json:"auth_token"`
	Name      string `json:"name"`
	Surname   string `json:"surname"`
	Email     string `json:"email"`
}

func RemoveLocalAuthToken() error {
	bytesToWrite, err := json.Marshal(UserInformation{})
	if err != nil {
		return err
	}
	os.WriteFile("./auth/auth.json", bytesToWrite, 0644)
	return nil
}

func AddLocalAuthToken(authToken, name, surname, email string) error {
	bytesToWrite, err := json.Marshal(UserInformation{
		AuthToken: authToken,
		Name: name,
		Surname: surname,
		Email: email,
	})
	if err != nil {
		return err
	}
	os.WriteFile("./auth/auth.json", bytesToWrite, 0644)
	return nil
}
