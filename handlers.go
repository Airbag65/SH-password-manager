package main

import (
	"SH-password-manager/db"
	"crypto/sha256"
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

func extractSalt(pwd string) (string, string) {
	saltBeginning := ""
	for i := len(pwd) - 2; i < len(pwd); i++ {
		saltBeginning = fmt.Sprintf("%s%c", saltBeginning, pwd[i])
	}

	saltEnd := ""
	for _, c := range pwd[:5] {
		saltEnd = fmt.Sprintf("%s%c", saltEnd, c)
	}
	return saltBeginning, saltEnd
}

func encryptPassword(origPwd string) string {
	encPassword := origPwd
	saltBeginning, saltEnd := extractSalt(encPassword)

	sha256Encoder := sha256.New()
	sha256Encoder.Write([]byte(encPassword))
	encPassword = fmt.Sprintf("%x", sha256Encoder.Sum(nil))

	encPassword = fmt.Sprintf("%s%s%s", saltBeginning, encPassword, saltEnd)

	sha256Encoder.Write([]byte(encPassword))
	encPassword = fmt.Sprintf("%x", sha256Encoder.Sum(nil))
	return encPassword
}

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

	newAuthToken := uuid.New().String()
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

	fmt.Println(db.GetUserWithEmail(request.Email).ToString())

	w.WriteHeader(200)
	w.Write(response)
}
