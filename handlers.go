package main

import (
	"SH-password-manager/db"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"net/http"

	// "github.com/google/uuid"
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
}

type LoginHandler struct{}

func (l *LoginHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(405)
		w.Write([]byte("Method Not Allowed\n"))
		return
	}

	var request LoginRequest

	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		w.WriteHeader(400)
		w.Write([]byte("Bad Request\n"))
		return
	}
    
    userInformation := db.GetUser(request.Email)
    if userInformation == nil {
        w.WriteHeader(404)
        w.Write([]byte("User not found"))
        return
    }

	sha256Encoder := sha256.New()
	sha256Encoder.Write([]byte(request.Password))
	_ = fmt.Sprintf("%x", sha256Encoder.Sum(nil))
	// fmt.Printf("%x\n", h.Sum(nil))
    response, err := json.Marshal(&LoginResponse{
        ResponseCode: 200,
        ResponseMessage: "OK",
        AuthToken: userInformation.AuthToken,
    })
    if err != nil {
        w.WriteHeader(500)
        w.Write([]byte("Internal server error"))
        return
    }

	w.WriteHeader(200)
	w.Write(response)
}
