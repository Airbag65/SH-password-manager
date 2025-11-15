package main

import (
	"encoding/json"
	"net/http"
)

func BadRequest(w http.ResponseWriter) {
	w.WriteHeader(400)
	w.Write([]byte("Bad Request"))
}

func InternalServerError(w http.ResponseWriter) {
	w.WriteHeader(500)
	w.Write([]byte("Internal Server Error"))
}

func NotFound(w http.ResponseWriter) {
	w.WriteHeader(404)
	w.Write([]byte("Not Found"))
}

func MethodNotAllowed(w http.ResponseWriter) {
	w.WriteHeader(405)
	w.Write([]byte("Method Not Allowed"))
}

func Unauthorized(w http.ResponseWriter) {
	w.WriteHeader(401)
	w.Write([]byte("Unauthorized"))
}

func WriteJSON(w http.ResponseWriter, v any) error {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(200)
	return json.NewEncoder(w).Encode(v)
}
