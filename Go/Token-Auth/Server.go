package main

import (
	"log"
	"time"
	"math/rand"
	"net/http"
	"github.com/gorilla/mux"
	"golang.org/x/crypto/bcrypt"
)

func checkPasswordHash(password, hash string) bool {
	match := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return match == nil
}

func main() {
	rand.Seed(time.Now().UnixNano()) // Make a seed for RNG
	token = make(map[string]Token) // Make a map for all tokens
	stopTokenCleaner := make(chan bool) // Make a quit channel to close cleaner
	go tokenCleaner(stopTokenCleaner) // Start token garbage collector

	r := mux.NewRouter()

	r.HandleFunc("/", indexPageHandler)
	r.HandleFunc("/login", loginHandler)
	r.HandleFunc("/internal", internalPageHandler)
	r.HandleFunc("/logout", logoutHandler)

	log.Printf("Server started listening.")
	http.ListenAndServe(":3737", r)
	stopTokenCleaner <- true
}