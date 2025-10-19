package main

import (
	"SH-password-manager/db"
	"SH-password-manager/enc"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/google/uuid"
)

type HomeHandler struct{}

func (h *HomeHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(200)
	w.Write([]byte("OK!\n"))
}

/*
   --- LOGIN ---
*/

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
	PemString       string `json:"pem_string"`
}

type LoginHandler struct{}

func (l *LoginHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(405)
		w.Write([]byte("Method Not Allowed"))
		return
	}

	if r.Header.Get("Content-Type") != "application/json" {
		w.WriteHeader(400)
		w.Write([]byte("Bad Request"))
		return
	}

	var request LoginRequest

	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		w.WriteHeader(400)
		w.Write([]byte("Bad Request"))
		return
	}

	userInformation := db.GetUserWithEmail(request.Email)
	if userInformation == nil {
		w.WriteHeader(404)
		w.Write([]byte("User not found"))
		return
	}

	if userInformation.AuthToken != "" {
		w.WriteHeader(418)
		w.Write([]byte("Already signed in"))
		return
	}

	if encryptPassword(request.Password) != userInformation.Password {
		w.WriteHeader(401)
		w.Write([]byte("Unauthorized"))
		return
	}

	pemString, err := enc.PEMFileToString("publicKey")
	if err != nil {
		log.Printf("Error: %v", err)	
		w.WriteHeader(500)
		w.Write([]byte("Internal server error"))
		return
	}

	newAuthToken := fmt.Sprintf("%s%s", uuid.New().String(), uuid.New().String())
	response, err := json.Marshal(&LoginResponse{
		ResponseCode:    200,
		ResponseMessage: "OK",
		AuthToken:       newAuthToken,
		Name:            userInformation.Name,
		Surname:         userInformation.Surname,
		Email:           userInformation.Email,
		PemString:       pemString,
	})
	if err != nil {
		log.Printf("Error: %v", err)	
		w.WriteHeader(500)
		w.Write([]byte("Internal server error"))
		return
	}

	db.SetAuthToken(request.Email, newAuthToken)

	w.WriteHeader(200)
	w.Write(response)
}

/*
   --- VALIDATE TOKEN ---
*/

type ValidateTokenHandler struct{}

type ValidateTokenRequest struct {
	AuthToken string `json:"auth_token"`
}

type ValidateTokenResponse struct {
	ResponseCode    int    `json:"response_code"`
	ResponseMessage string `json:"response_message"`
	Name            string `json:"name"`
	Surname         string `json:"surname"`
	Email           string `json:"email"`
	PemString       string `json:"pem_string"`
}

func (v *ValidateTokenHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(405)
		w.Write([]byte("Method Not Allowed"))
		return
	}

	if r.Header.Get("Content-Type") != "application/json" {
		w.WriteHeader(400)
		w.Write([]byte("Bad Request"))
		return
	}

	var request ValidateTokenRequest

	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		w.WriteHeader(400)
		w.Write([]byte("Bad Request"))
		return
	}

	userInformation := db.GetUserWithAuthToken(request.AuthToken)
	if userInformation == nil {
		w.WriteHeader(401)
		w.Write([]byte("Unauthorized"))
		return
	}

	pemString, err := enc.PEMFileToString("publicKey")
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte("Internal Server Error"))
		return
	}

	response, err := json.Marshal(&ValidateTokenResponse{
		ResponseCode:    200,
		ResponseMessage: "OK",
		Name:            userInformation.Name,
		Surname:         userInformation.Surname,
		Email:           userInformation.Email,
		PemString:       pemString,
	})
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte("Internal server error"))
		return
	}
	w.WriteHeader(200)
	w.Write(response)
}

/*
   --- SIGN OUT ---
*/

type SignOutHandler struct{}

type SignOutRequest struct {
	Email string `json:"email"`
}

type SignOutResponse struct {
	ResponseCode    int    `json:"response_code"`
	ResponseMessage string `json:"response_message"`
}

func (s *SignOutHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		w.WriteHeader(405)
		w.Write([]byte("Method Not Allowed"))
		return
	}

	if r.Header.Get("Content-Type") != "application/json" {
		w.WriteHeader(400)
		w.Write([]byte("Bad Request"))
		return
	}

	var request SignOutRequest

	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		w.WriteHeader(400)
		w.Write([]byte("Bad Request"))
		return
	}

	response, err := json.Marshal(&SignOutResponse{
		ResponseCode:    200,
		ResponseMessage: "OK",
	})
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte("Internal server error"))
		return
	}

	user := db.GetUserWithEmail(request.Email)
	if user == nil {
		w.WriteHeader(404)
		w.Write([]byte("Not found"))
		return
	}

	if user.AuthToken == "" {
		w.WriteHeader(304)
		w.Write([]byte("Not modified"))
		return
	}

	db.RemoveAuthToken(request.Email)

	w.WriteHeader(200)
	w.Write(response)
}

/*
--- CREATE NEW USER ---
*/
type CreateNewUserHandler struct{}

type CreateNewUserRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Name     string `json:"name"`
	Surname  string `json:"surname"`
}

type CreateNewUserResponse struct {
	ResponseCode    int    `json:"response_code"`
	ResponseMessage string `json:"response_message"`
	AuthToken       string `json:"auth_token"`
	Name            string `json:"name"`
	Surname         string `json:"surname"`
	PemString       string `json:"pem_string"`
}

func (c *CreateNewUserHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(405)
		w.Write([]byte("Method Not Allowed"))
		return
	}

	if r.Header.Get("Content-Type") != "application/json" {
		w.WriteHeader(400)
		w.Write([]byte("Bad Request"))
		return
	}

	var request CreateNewUserRequest

	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		w.WriteHeader(400)
		w.Write([]byte("Bad Request"))
		return
	}

	existingUser := db.GetUserWithEmail(request.Email)
	if existingUser != nil {
		w.WriteHeader(418)
		w.Write([]byte("User already exists"))
		return
	}

	pemString, err := enc.PEMFileToString("publicKey")
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte("Internal Server Error"))
		return
	}

	encPwd := encryptPassword(request.Password)
	db.CreateNewUser(request.Email, encPwd, request.Name, request.Surname)
	newAuthToken := fmt.Sprintf("%s%s", uuid.New().String(), uuid.New().String())
	response, err := json.Marshal(&CreateNewUserResponse{
		ResponseCode:    200,
		ResponseMessage: "OK",
		AuthToken:       newAuthToken,
		Name:            request.Name,
		Surname:         request.Surname,
		PemString:       pemString,
	})
	db.SetAuthToken(request.Email, newAuthToken)
	w.WriteHeader(200)
	w.Write(response)
}
