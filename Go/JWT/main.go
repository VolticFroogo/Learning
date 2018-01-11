package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"time"

	"github.com/codegangsta/negroni"
	"github.com/dgrijalva/jwt-go"
)

var secretKey []byte

// Structs

type UserCredentials struct {
	Name     string `json:"name"`
	Password string `json:"password"`
}

type JwtClaims struct {
	Name string `json:"name"`
	jwt.StandardClaims
}

type Token struct {
	Token string `json:"token"`
}

type Response struct {
	Data string `json:"data"`
}

// Helper

func JsonResponse(response interface{}, w http.ResponseWriter) {
	json, err := json.Marshal(response)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintln(w, "Error generating JSON response.")
		fmt.Println("Error generating JSON response:", err)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write(json)
}

// Handlers

func loginHandler(w http.ResponseWriter, r *http.Request) {
	var user UserCredentials

	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Println("Error decoding JSON:", err)
		fmt.Fprintln(w, "Error decoding JSON.")
		return
	}

	if user.Name != "Froogo" || user.Password != "Super-Secret-Password" {
		w.WriteHeader(http.StatusForbidden)
		fmt.Fprintln(w, "Invalid credentials.")
		return
	}

	claims := JwtClaims{
		user.Name,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(30 * time.Second).Unix(),
		},
	}

	rawToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	token, err := rawToken.SignedString(secretKey)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintln(w, "Error generating JWT.")
		fmt.Println("Error generating JWT:", err)
		return
	}

	response := Token{token}
	JsonResponse(response, w)
}

func protectedHandler(w http.ResponseWriter, r *http.Request) {
	response := Response{"Protected resource!"}
	JsonResponse(response, w)
}

// Middleware

func validateTokenMiddleware(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	var d Token
	err := json.NewDecoder(r.Body).Decode(&d)
	if err != nil {
		fmt.Printf("JSON decoding error:", err)
		fmt.Fprintln(w, "JSON deocding error.")
		w.WriteHeader(http.StatusInternalServerError)
	}

	token, err := jwt.Parse(d.Token, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v.", token.Header["alg"])
		}

		return secretKey, nil
	})

	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		fmt.Fprintln(w, "Invalid token!")
		return
	}

	if token.Valid {
		next(w, r)
	} else {
		w.WriteHeader(http.StatusUnauthorized)
		fmt.Fprintln(w, "Invalid token!")
	}
}

// Main function

func main() {
	rand.Seed(time.Now().UnixNano())
	secretKey := make([]byte, 32)
	rand.Read(secretKey)
	fmt.Println("\nSecret key:", secretKey)

	//PUBLIC ENDPOINTS
	http.HandleFunc("/", loginHandler)

	//PROTECTED ENDPOINTS
	http.Handle("/protected", negroni.New(
		negroni.HandlerFunc(validateTokenMiddleware),
		negroni.Wrap(http.HandlerFunc(protectedHandler)),
	))

	fmt.Println("Now listening...")
	http.ListenAndServe(":3737", nil)
}
