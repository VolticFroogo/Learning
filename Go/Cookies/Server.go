package main

import (
	"fmt"
	"log"
	"strings"
	"github.com/gorilla/mux"
	"github.com/gorilla/securecookie"
	"net/http"
)

// Split address for logging

type Address struct {
	ip, port string
}

func splitAddress(address string) Address {
	// Split IP:Port
	split := strings.Split(address, ":")

	// Fix for when localhost returns as [::1]:port
	if (len(split) > 2) {
		return Address{"localhost", split[len(split) - 1]}
	} else {
		return Address{split[0], split[1]}
	}

	// Returns IP and Port
}

// Cookie handling

var cookieHandler = securecookie.New(
	securecookie.GenerateRandomKey(64),
	securecookie.GenerateRandomKey(32))

func getUserName(request *http.Request) (userName string) {
	if cookie, err := request.Cookie("session"); err == nil {
		cookieValue := make(map[string]string)
		if err = cookieHandler.Decode("session", cookie.Value, &cookieValue); err == nil {
			userName = cookieValue["name"]
		}
	}
	return userName
}

func setSession(userName string, response http.ResponseWriter) {
	value := map[string]string{
		"name": userName,
	}
	if encoded, err := cookieHandler.Encode("session", value); err == nil {
		cookie := &http.Cookie{
			Name:  "session",
			Value: encoded,
			Path:  "/",
		}
		http.SetCookie(response, cookie)
	}
}

func clearSession(response http.ResponseWriter) {
	cookie := &http.Cookie{
		Name:   "session",
		Value:  "",
		Path:   "/",
		MaxAge: -1,
	}
	http.SetCookie(response, cookie)
}

// Login handler

func loginHandler(response http.ResponseWriter, request *http.Request) {
	name := request.FormValue("name")
	redirectTarget := "/"
	if name != "" {
		setSession(name, response)
		redirectTarget = "/internal"
		log.Printf("%v successfully logged in from: %v.", name, splitAddress(request.RemoteAddr).ip)
	}
	http.Redirect(response, request, redirectTarget, 302)
}

// Logout handler

func logoutHandler(response http.ResponseWriter, request *http.Request) {
	clearSession(response)
	http.Redirect(response, request, "/", 302)
}

// Index

const indexPage = `
<h1>Login</h1>
<hr>
<form method="post" action="/login">
    <label for="name">Username: </label>
    <input type="text" id="name" name="name">
    <button type="submit">Login</button>
</form>
`

func indexPageHandler(response http.ResponseWriter, request *http.Request) {
	userName := getUserName(request)
	if userName == "" {
		fmt.Fprintf(response, indexPage)
	} else {
		http.Redirect(response, request, "/internal", 302)
	}
}

// Internal

const internalPage = `
<h1>Internal</h1>
<hr>
<small>Username: %v</small>
<form method="post" action="/logout">
    <button type="submit">Logout</button>
</form>
`

func internalPageHandler(response http.ResponseWriter, request *http.Request) {
	userName := getUserName(request)
	if userName != "" {
		fmt.Fprintf(response, internalPage, userName)
	} else {
		http.Redirect(response, request, "/", 302)
	}
}

// Request handler

var router = mux.NewRouter()

func main() {
	router.HandleFunc("/", indexPageHandler)
	router.HandleFunc("/internal", internalPageHandler)
	router.HandleFunc("/login", loginHandler).Methods("POST")
	router.HandleFunc("/logout", logoutHandler).Methods("POST")

	http.Handle("/", router)
	log.Printf("Server started listening.")
	http.ListenAndServe(":3737", nil)
}