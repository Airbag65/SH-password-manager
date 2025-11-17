package main

import (
	"SH-password-manager/enc"
	"encoding/json"
	"log"
	"net/http"
	"strings"
)

// 4673992a-f65c-4da8-99f0-fa630f54ed28ec0e1431-1674-4c32-a606-efdd160862c7


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

	userInformation := s.GetUserWithAuthToken(token)
	if userInformation == nil {
		Unauthorized(w)
		return
	}

	names := s.GetHostNames(userInformation.Id)

	// bytes, err := json.Marshal(names)
	// if err != nil {
	// 	InternalServerError(w)
	// 	return
	// }
	WriteJSON(w, GetPasswordHostsResponse{Hosts: names})
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
	userInformation := s.GetUserWithAuthToken(token)
	if userInformation == nil {
		Unauthorized(w)
		return
	}

	err = s.AddNewPassord(userInformation.Id, request.Password, request.HostName)
	if err != nil {
		InternalServerError(w)
		return
	}

	w.WriteHeader(200)
	w.Write([]byte("OK"))
}


func (h *GetPasswordValueHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
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

	userInformation := s.GetUserWithAuthToken(token)
	if userInformation == nil {
		Unauthorized(w)
		return
	}

	var request GetPasswordRequest

	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		BadRequest(w)
		return
	}
	
	encPassword, err := s.GetPassword(userInformation.Id, request.HostName)
	if err != nil {
		NotFound(w)
		return
	}

	privatePemString, err := enc.PEMFileToString("privateKey")
	if err != nil {
		InternalServerError(w)
		return
	}
	
	privateKey := enc.PemStringToPrivateKey(privatePemString)

	decPassword, err := enc.Decrypt([]byte(encPassword), privateKey)
	if err != nil {
		log.Printf("Error: %v", err)
		InternalServerError(w)
		return
	}

	WriteJSON(w, GetPasswordResonse{Password: decPassword})
}
