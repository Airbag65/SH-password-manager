package auth

import (
	"encoding/json"
	"os"
)

func RemoveLocalAuthToken() error {
	bytesToWrite, err := json.Marshal(LocalAuth{})
	if err != nil {
		return err
	}
	os.WriteFile("./auth/auth.json", bytesToWrite, 0644)
	return nil
}

func AddLocalAuthToken(authToken string) error {
	bytesToWrite, err := json.Marshal(LocalAuth{
		AuthToken: authToken,
	})
	if err != nil {
		return err
	}
	os.WriteFile("./auth/auth.json", bytesToWrite, 0644)
	return nil
}
