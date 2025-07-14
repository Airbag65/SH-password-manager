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

	if len(os.Args) > 1 && os.Args[1] == "migrate" {
		db.Migrate()
        fmt.Println(db.GetUser("normananton03@gmail.com").ToString())
		return
	}

	server := http.NewServeMux()

	server.Handle("/", &HomeHandler{})
	server.Handle("/login", &LoginHandler{})

	handler := cors.Default().Handler(server)
	http.ListenAndServe(":8080", handler)
}
