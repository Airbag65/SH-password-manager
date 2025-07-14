package main

import (
	"crypto/sha256"
	"fmt"
	"net/http"
)

type homeHandler struct{}
func (h* homeHandler) ServeHTTP(w http.ResponseWriter, r* http.Request) {
    w.WriteHeader(200)
    w.Write([]byte("OK!\n"))
}

type loginHandler struct{}
func (l* loginHandler) ServeHTTP(w http.ResponseWriter, r* http.Request) {
    h := sha256.New()
    h.Write([]byte("Testing testing...\n"))
    fmt.Printf("%x\n", h.Sum(nil))
    w.WriteHeader(200)
    w.Write([]byte("OK!\n"))
}
