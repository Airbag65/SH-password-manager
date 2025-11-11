package main

import "net/http"

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
