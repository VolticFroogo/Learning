package main

import (
	"log"
	"time"
	"net/http"
	"html/template"
)

type pageVariables struct {
	Username string
}

func indexPageHandler(w http.ResponseWriter, r *http.Request) {
	tokenID := getSession("Token", r) // Get tokenID from cookie
	if checkToken(tokenID) {
		http.Redirect(w, r, "/internal", 302) // If token exists and is valid redirect to internal
		return
	}

	t, err := template.ParseFiles("Websites/index.html") // Parse the index HTML page
	if err != nil {
		log.Print("Template parsing error: ", err)
	}

	temp := 0
	err = t.Execute(w, temp) // Execute temmplate with variables
	if err != nil {
		log.Print("Template execution error: ", err)
	}
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	username := r.FormValue("username")
	password := r.FormValue("password")
	redirectTarget := "/"
	if username != "" && password != "" {
		id, valid := checkLogin(username, password)
		if (valid) {
			tokenID := GenerateToken()
			token[tokenID] = Token {
				userid: id,
				expires: int(time.Now().Unix()) + tokenLifetime,
			}
			setSession("Token", tokenID, w)
			redirectTarget = "/internal"
		}
	}
	http.Redirect(w, r, redirectTarget, 302) // Login not valid
}

func internalPageHandler(w http.ResponseWriter, r *http.Request) {
	tokenID := getSession("Token", r) // Get tokenID from cookie
	if checkToken(tokenID) { // Check if token is valid
		t, err := template.ParseFiles("Websites/internal.html") // Parse the internal HTML page
		if err != nil {
			log.Print("Template parsing error: ", err)
		}

		internalPageVariables := pageVariables {
			Username: dbGetUsername(token[tokenID].userid), // Get username from token for template
		}

		err = t.Execute(w, internalPageVariables) // Execute temmplate with variables
		if err != nil {
			log.Print("Template execution error: ", err)
		}

		return // Page is sent, return
	}

	http.Redirect(w, r, "/", 302) // Redirect back to index: token not valid
}

func logoutHandler(w http.ResponseWriter, r *http.Request) {
	tokenID := getSession("Token", r) // Get tokenID from cookie
	delete(token, tokenID) // Revoke token
	clearSession(w) // Delete session data
	http.Redirect(w, r, "/", 302) // Redirect to index
}