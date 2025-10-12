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

func ValidTokenExists() bool {
	jsonFile, err := os.Open("auth/auth.json")
	if err != nil {
		fmt.Printf("An error occured while opening file: %v\n", err)
		return false
	}

	defer func(){
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

	client := http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}
	request.Header.Set("Content-Type", "application/json")

	response, err := client.Do(request)
	if err != nil {
		fmt.Printf("An error occured while sending request: %v\n", err)
	}

	if response.StatusCode != 200 {
		bytesToWrite, err := json.Marshal(LocalAuth{})
		if err != nil {
			fmt.Println("Could now write to file")
		}
		os.WriteFile("./auth/auth.json", bytesToWrite, 0644)

		return false
	}
	return true
}
