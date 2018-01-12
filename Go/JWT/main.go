package main

import (
	"crypto/rsa"
	"encoding/json"
	"fmt"
	"io/ioutil"
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

const (
	privKeyPath = "keys/app.rsa"     // `$ openssl genrsa -out app.rsa 2048`
	pubKeyPath  = "keys/app.rsa.pub" // `$ openssl rsa -in app.rsa -pubout > app.rsa.pub`
)

var (
	rsaVerifyKey *rsa.PublicKey
	rsaSignKey   *rsa.PrivateKey
)

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

	// if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
	if token.Valid {
		next(w, r)
	} else {
		w.WriteHeader(http.StatusUnauthorized)
		fmt.Fprintln(w, "Invalid token!")
	}
}

func initJWT() error {
	signBytes, err := ioutil.ReadFile(privKeyPath)
	if err != nil {
		return err
	}

	rsaSignKey, err = jwt.ParseRSAPrivateKeyFromPEM(signBytes)
	if err != nil {
		return err
	}

	verifyBytes, err := ioutil.ReadFile(pubKeyPath)
	if err != nil {
		return err
	}

	rsaVerifyKey, err = jwt.ParseRSAPublicKeyFromPEM(verifyBytes)
	if err != nil {
		return err
	}

	return nil
}

// Main function

func main() {
	rand.Seed(time.Now().UnixNano())
	secretKey := make([]byte, 32)
	rand.Read(secretKey)
	initJWT()
	fmt.Println("\nSecret key:", secretKey)

	//PUBLIC ENDPOINTS
	http.HandleFunc("/", loginHandler)

	//PROTECTED ENDPOINTS
	http.Handle("/protected", negroni.New(
		negroni.HandlerFunc(validateTokenMiddleware),
		negroni.Wrap(http.HandlerFunc(protectedHandler)),
	))

	/*

		Generating an RSA JWT is the exact same as with HMAC.
		The only exception is you sign with sign and verify with verify.
		Rather than the usual sign/sign with HMAC.
		An example is below for reference.

	*/

	claims := JwtClaims{
		"Froogo",
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(30 * time.Second).Unix(),
		},
	}

	rawToken := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)

	tokenString, err := rawToken.SignedString(rsaSignKey)
	if err != nil {
		fmt.Println("Error generating JWT:", err)
		return
	}
	fmt.Println("\nNew JWT:", tokenString)

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v.", token.Header["alg"])
		}

		return rsaVerifyKey, nil
	})
	if err != nil {
		fmt.Println("Invalid token!")
		return
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		fmt.Printf("\nUser: %v has a valid token!\n", claims["name"])
	} else {
		fmt.Println("Invalid token!")
	}

	fmt.Println("\nNow listening...")
	http.ListenAndServe(":3737", nil)
}
