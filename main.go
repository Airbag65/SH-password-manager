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
	f, err := os.OpenFile("logs/file.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		fmt.Printf("error opening file: %v", err)
		return
	}
	defer f.Close()
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
	// if len(os.Args) > 1 && os.Args[1] == "migrate" {
	// }
	// if len(os.Args) > 1 && os.Args[1] == "experiment" {
	// }

	server := http.NewServeMux()

	server.Handle("/", &HomeHandler{})
	server.Handle("/login", &LoginHandler{})
	server.Handle("/validateToken", &ValidateTokenHandler{})
	server.Handle("/signOut", &SignOutHandler{})
	server.Handle("/createUser", &CreateNewUserHandler{})

	handler := cors.Default().Handler(server)
	http.ListenAndServeTLS(":443", "cert.pem", "key.pem" handler)
}
