package main

import (
	"encoding/base64"
	"log"
	"net/http"
	"strings"
)

func main() {
	myHandler := &MyHandler{}
	http.Handle("/a", &AuthHandler{myHandler})
	log.Fatal(http.ListenAndServe("localhost:8080", nil))
}

type MyHandler struct{}

func (h *MyHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("authorized"))
}

type AuthHandler struct {
	handler http.Handler
}

func (a *AuthHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	auth, ok := r.Header["Authorization"]
	if !ok {
		w.Header().Add("WWW-Authenticate", "basic realm=\"please authenticate\"")
		w.WriteHeader(http.StatusUnauthorized)
		log.Printf("Unauthorized access to %s", r.URL)
		return
	}

	encoded := strings.Split(auth[0], " ")
	decoded, _ := base64.StdEncoding.DecodeString(encoded[1])
	parts := strings.Split(string(decoded), ":")
	if parts[0] == "admin" && parts[1] == "a" {
		a.handler.ServeHTTP(w, r)
	} else {
		w.WriteHeader(http.StatusUnauthorized)
		log.Printf("Unauthorized access to %s", r.URL)
	}
}
