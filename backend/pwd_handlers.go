package main

import (
	"SH-password-manager/db"
	"encoding/json"
	"net/http"
	"strings"
)

// 4673992a-f65c-4da8-99f0-fa630f54ed28ec0e1431-1674-4c32-a606-efdd160862c7

type GetPasswordHostsHandler struct{}

func (h *GetPasswordHostsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		MethodNotAllowed(w)
		return
	}

	tokenHeader := r.Header.Get("Authorization")
	if tokenHeader == "" {
		BadRequest(w)
		return
	}

	bearer := strings.Split(tokenHeader, " ")[0]
	if bearer != "Bearer" {
		BadRequest(w)
		return
	}
	token := strings.Split(tokenHeader, " ")[1]

	userInformation := db.GetUserWithAuthToken(token)
	if userInformation == nil {
		Unauthorized(w)
		return
	}

	names := db.GetHostNames(userInformation.Id)

	bytes, err := json.Marshal(names)
	if err != nil {
		InternalServerError(w)
		return
	}

	w.WriteHeader(200)
	w.Write(bytes)
}

type UploadNewPasswordHandler struct{}

type UploadNewPasswordRequest struct {
	HostName string `json:"host_name"`
	Password string `json:"password"`
}

func (h *UploadNewPasswordHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		MethodNotAllowed(w)
		return
	}


	var request UploadNewPasswordRequest

	if r.Header.Get("Content-Type") != "application/json" {
		BadRequest(w)
		return
	}

	tokenHeader := r.Header.Get("Authorization")
	if tokenHeader == "" {
		BadRequest(w)
		return
	}

	bearer := strings.Split(tokenHeader, " ")[0]
	if bearer != "Bearer" {
		BadRequest(w)
		return
	}
	token := strings.Split(tokenHeader, " ")[1]

	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		BadRequest(w)
		return
	}
	userInformation := db.GetUserWithAuthToken(token)
	if userInformation == nil {
		Unauthorized(w)
		return
	}

	err = db.AddNewPassord(userInformation.Id, request.Password, request.HostName)
	if err != nil {
		InternalServerError(w)
		return
	}

	w.WriteHeader(200)
	w.Write([]byte("OK"))
}
