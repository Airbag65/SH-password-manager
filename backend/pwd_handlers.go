package main

import (
	"SH-password-manager/db"
	"encoding/json"
	"net/http"
)

type GetPasswordHostsHandler struct{}

type GetPasswordHostsRequest struct {
	AuthToken string `json:"auth_token"`
}

func (h *GetPasswordHostsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(405)
		w.Write([]byte("Method Not Allowed"))
		return
	}

	var request GetPasswordHostsRequest

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
	
	// TODO: Actually get the passwords

	w.WriteHeader(200)
	w.Write([]byte("OK"))
}

type UploadNewPasswordHandler struct{}

type UploadNewPasswordRequest struct {
	AuthToken string `json:"auth_token"`
	HostName  string `json:"host_name"`
	Password  string `json:"password"`
}

func (h *UploadNewPasswordHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(405)
		w.Write([]byte("Method Not Allowed"))
		return
	}

	var request UploadNewPasswordRequest

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
	
	err = db.AddNewPassord(userInformation.Id, request.Password, request.HostName)
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte("Internal Server Error"))
		return
	}
	
	w.WriteHeader(200)	
	w.Write([]byte("OK"))
}
