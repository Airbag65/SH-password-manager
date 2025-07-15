package main

import (
	"SH-password-manager/db"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/google/uuid"
)

type HomeHandler struct{}

func (h *HomeHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(200)
	w.Write([]byte("OK!\n"))
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
}

type LoginHandler struct{}

func (l *LoginHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(405)
		w.Write([]byte("Method Not Allowed\n"))
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
		w.Write([]byte("Bad Request\n"))
		return
	}

	userInformation := db.GetUserWithEmail(request.Email)
	if userInformation == nil {
		w.WriteHeader(404)
		w.Write([]byte("User not found"))
		return
	}

	if encryptPassword(request.Password) != userInformation.Password {
		w.WriteHeader(401)
		w.Write([]byte("Unauthorized"))
		return
	}

	newAuthToken := fmt.Sprintf("%s%s",uuid.New().String(), uuid.New().String())
	response, err := json.Marshal(&LoginResponse{
		ResponseCode:    200,
		ResponseMessage: "OK",
		AuthToken:       newAuthToken,
		Name:            userInformation.Name,
		Surname:         userInformation.Surname,
	})
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte("Internal server error"))
		return
	}

	db.SetAuthToken(request.Email, newAuthToken)

	w.WriteHeader(200)
	w.Write(response)
}

type ValidateTokenHandler struct{}

type ValidateTokenRequest struct {
	AuthToken string `json:"auth_token"`
}

type ValidateTokenResponse struct {
	ResponseCode    int    `json:"response_code"`
	ResponseMessage string `json:"response_message"`
	Name            string `json:"name"`
	Surname         string `json:"surname"`
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
		w.Write([]byte("Bad Request\n"))
		return
	}

	userInformation := db.GetUserWithAuthToken(request.AuthToken)
	if userInformation == nil {
		w.WriteHeader(401)
		w.Write([]byte("Unauthorized"))
		return
	}

	response, err := json.Marshal(&ValidateTokenResponse{
		ResponseCode:    200,
		ResponseMessage: "OK",
		Name:            userInformation.Name,
		Surname:         userInformation.Surname,
	})
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte("Internal server error"))
		return
	}
	w.WriteHeader(200)
	w.Write(response)
}

type SignOutHandler struct{}


func (s *SignOutHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		w.WriteHeader(405)
		w.Write([]byte("Method Not Allowed"))
		return
	}

    // TODO: FINISH THIS STUFF
}
