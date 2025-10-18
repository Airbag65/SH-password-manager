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
		w.WriteHeader(405)
		w.Write([]byte("Method Not Allowed"))
		return
	}

	tokenHeader := r.Header.Get("Authorization")
	if tokenHeader == "" {
		w.WriteHeader(400)
		w.Write([]byte("Bad Request"))
		return
	}

	bearer := strings.Split(tokenHeader, " ")[0]
	if bearer != "Bearer" {
		w.WriteHeader(400)
		w.Write([]byte("Bad Request"))
		return
	}
	token := strings.Split(tokenHeader, " ")[1]

	userInformation := db.GetUserWithAuthToken(token)
	if userInformation == nil {
		w.WriteHeader(401)
		w.Write([]byte("Unauthorized"))
		return
	}

	names := db.GetHostNames(userInformation.Id)

	bytes, err := json.Marshal(names)
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte("Internal Server Error"))
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
		w.WriteHeader(405)
		w.Write([]byte("Method Not Allowed"))
		return
	}

	var request UploadNewPasswordRequest

	if r.Header.Get("Content-Type") != "application/json" {
		w.WriteHeader(400)
		w.Write([]byte("Bad Request"))
		return
	}

	tokenHeader := r.Header.Get("Authorization")
	if tokenHeader == "" {
		w.WriteHeader(400)
		w.Write([]byte("Bad Request"))
		return
	}

	bearer := strings.Split(tokenHeader, " ")[0]
	if bearer != "Bearer" {
		w.WriteHeader(400)
		w.Write([]byte("Bad Request"))
		return
	}
	token := strings.Split(tokenHeader, " ")[1]

	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		w.WriteHeader(400)
		w.Write([]byte("Bad Request"))
		return
	}
	userInformation := db.GetUserWithAuthToken(token)
	if userInformation == nil {
		w.WriteHeader(401)
		w.Write([]byte("Unauthorized"))
		return
	}

	err = db.AddNewPassord(userInformation.Id, request.Password, request.HostName)
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte("Internal Server Error"))
		return
	}

	w.WriteHeader(200)
	w.Write([]byte("OK"))
}
