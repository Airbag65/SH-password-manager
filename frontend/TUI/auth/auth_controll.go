package auth

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

type LocalAuth struct {
	AuthToken string `json:"auth_token"`
}

var (
	client = http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}
)

func ValidTokenExists() bool {
	jsonFile, err := os.Open("auth/auth.json")
	if err != nil {
		fmt.Printf("An error occured while opening file: %v\n", err)
		return false
	}

	defer func() {
		if err = jsonFile.Close(); err != nil {
			fmt.Println("Could not close file")
		}
	}()

	fileBytes, err := io.ReadAll(jsonFile)
	if err != nil {
		fmt.Printf("An error occured while parsing file content: %v\n", err)
		return false
	}

	var localAuth LocalAuth

	json.Unmarshal(fileBytes, &localAuth)

	if localAuth.AuthToken == "" {
		return false
	}

	requestBody, err := json.Marshal(localAuth)
	if err != nil {
		fmt.Printf("An error occured while serializing request body: %v\n", err)
		return false
	}

	request, err := http.NewRequest("POST", "https://localhost:443/auth/valid", bytes.NewBuffer(requestBody))
	if err != nil {
		fmt.Printf("An error occured while constructing request: %v\n", err)
	}
	// 8239d0d6-8085-4b30-8e8a-7c6052611307191bec5a-0d95-4177-bb77-9a6ad72e40f0
	request.Header.Set("Content-Type", "application/json")

	response, err := client.Do(request)
	if err != nil {
		fmt.Printf("An error occured while sending request: %v\n", err)
	}

	if response.StatusCode != 200 {
		err = RemoveLocalAuthToken()
		if err != nil {
			fmt.Println("Could now write to file")
		}
		// bytesToWrite, err := json.Marshal(LocalAuth{})
		// if err != nil {
		// }
		// os.WriteFile("./auth/auth.json", bytesToWrite, 0644)

		return false
	}
	return true
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginResponse struct {
	ResponseCode    int    `json:"response_code"`
	ResponseMessage string `json:"response_message"`
	AuthToken       string `json:"auth_token"`
	Name            string `json:"name"`
	Surname         string `json:"surname"`
	Email           string `json:"email"`
}

func Login(email, password string) (*LoginResponse, error) {
	loginRequestBody := LoginRequest{
		Email:    email,
		Password: password,
	}

	requestBodyBytes, err := json.Marshal(loginRequestBody)
	if err != nil {
		return nil, err
	}

	request, err := http.NewRequest("POST", "https://localhost:443/auth/login", bytes.NewBuffer(requestBodyBytes))
	if err != nil {
		return nil, err
	}

	request.Header.Set("Content-Type", "application/json")

	response, err := client.Do(request)
	if err != nil {
		return nil, err
	}

	var buffer []byte
	if response.StatusCode == 200 {
		buffer, err = io.ReadAll(response.Body)
		if err != nil {
			return nil, err
		}
	} else {
		return nil, fmt.Errorf("StatusCode was not 200, but was %d", response.StatusCode)
	}

	var loginRes LoginResponse

	if err = json.Unmarshal(buffer, &loginRes); err != nil {
		return nil, err
	}

	err = AddLocalAuthToken(loginRes.AuthToken, loginRes.Name, loginRes.Surname, loginRes.Email)
	if err != nil {
		return nil, err
	}

	return &loginRes, nil
}
