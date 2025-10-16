package auth

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
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
	localData := GetSavedData()

	if localData.AuthToken == "" {
		return false
	}

	localAuth := LocalAuth{
		AuthToken: localData.AuthToken,
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
	if email == "" || password == "" {
		return nil, fmt.Errorf("No email or password provided\n")
	}
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

type SignupRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Name     string `json:"name"`
	Surname  string `json:"surname"`
}

type SignupResponse struct {
	ResponseCode    int    `json:"response_code"`
	ResponseMessage string `json:"response_message"`
	AuthToken       string `json:"auth_token"`
	Name            string `json:"name"`
	Surname         string `json:"surname"`
}

func SignUp(email, password, name, surname string) (*SignupResponse, error) {
	if email == "" || password == ""  || name == "" || surname == ""{
		return nil, fmt.Errorf("Insufficient infromation provided\n")
	}
	
	signupRequest := SignupRequest{
		Email:    email,
		Password: password,
		Name:     name,
		Surname:  surname,
	}

	reqBodyBytes, err := json.Marshal(signupRequest)
	if err != nil {
		return nil, err
	}

	request, err := http.NewRequest("POST", "https://localhost:443/auth/new", bytes.NewBuffer(reqBodyBytes))
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

	var signupResponse SignupResponse

	err = json.Unmarshal(buffer, &signupResponse)
	if err != nil {
		return nil, err
	}

	err = AddLocalAuthToken(signupResponse.AuthToken, signupResponse.Name,signupResponse.Surname, email)
	if err != nil {
		return nil, err
	}

	return &signupResponse, nil
}
