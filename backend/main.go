package main

import (
	"SH-password-manager/db"

	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/rs/cors"
)

func main() {
	f, err := os.OpenFile("target/log/file.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		fmt.Printf("error opening file: %v", err)
		return
	}
	defer func() {
		err = f.Close()
		if err != nil {
			fmt.Println("Could not close log-file")
		}
	}()

	log.SetOutput(f)

	if len(os.Args) > 1 {
		switch os.Args[1] {
		case "migrate":
			db.Migrate()
			fmt.Println(db.GetUserWithEmail("normananton03@gmail.com").ToString())
			return
		default:
			return
		}
	}

	server := http.NewServeMux()

	// Auth handlers
	server.Handle("/", &HomeHandler{})
	server.Handle("/login", &LoginHandler{})
	server.Handle("/validateToken", &ValidateTokenHandler{})
	server.Handle("/signOut", &SignOutHandler{})
	server.Handle("/createUser", &CreateNewUserHandler{})

	// PWD handlers
	server.Handle("/getPasswordHosts", &GetPasswordHostsHandler{})

	handler := cors.Default().Handler(server)
	fmt.Println("Server running on: https://localhost:443 ...")
	err = http.ListenAndServeTLS(":443", "cert.pem", "key.pem", handler)
	if err != nil {
		log.Fatal("Could not start server")
	}
}
