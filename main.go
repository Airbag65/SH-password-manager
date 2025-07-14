package main

import (
	"net/http"

	"github.com/rs/cors"
)


func main()  {
    server := http.NewServeMux()

    server.Handle("/", &homeHandler{})
    server.Handle("/login", &loginHandler{})

    handler := cors.Default().Handler(server)
    http.ListenAndServe(":8080", handler)
}
