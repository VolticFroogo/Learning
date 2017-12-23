package main

import (
	"github.com/gorilla/securecookie"
	"net/http"
)

var hashKey = securecookie.GenerateRandomKey(64)        // Generate hash key for cookie encoding
var blockKey = securecookie.GenerateRandomKey(32)       // Generate block key for cookie encoding
var cookieHandler = securecookie.New(hashKey, blockKey) // Make cookie handler from hash key and block key

func getSession(cookieName string, request *http.Request) (cookieData string) {
	if cookie, err := request.Cookie("session"); err == nil { // Get session data
		cookieValue := make(map[string]string) // Make a map for cookie value
		if err = cookieHandler.Decode("session", cookie.Value, &cookieValue); err == nil {
			cookieData = cookieValue[cookieName] // Read data from cookie
		}
	}
	return
}

func setSession(cookieName string, cookieData string, response http.ResponseWriter) {
	value := map[string]string{
		cookieName: cookieData, // Write cookie in map
	}
	if encoded, err := cookieHandler.Encode("session", value); err == nil {
		cookie := &http.Cookie{
			Name:  "session",
			Value: encoded,
			Path:  "/",
		}
		http.SetCookie(response, cookie) // Set encoded cookie
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
