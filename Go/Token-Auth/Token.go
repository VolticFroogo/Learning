package main

import (
	"log"
	"math/rand"
	"time"
)

// Define tokens

const tokenLifetime = 7200 // Tokens are valid for 2 hours

type Token struct {
	userid, expires int
}

var token map[string]Token // Initialize a map for all tokens

// Check if token is valid

func checkToken(tokenID string) bool {
	if token[tokenID].userid != 0 { // Check if token exists
		if token[tokenID].expires > int(time.Now().Unix()) { // Check if token has expired
			return true // Token is valid
		} else { // Token exists but has expired
			delete(token, tokenID) // Delete expired token
		}
	}

	return false // Token is invalid
}

// Generate a token

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
const (
	letterIdxBits = 6                    // 6 bits to represent a letter index
	letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
	letterIdxMax  = 63 / letterIdxBits   // # of letter indices fitting in 63 bits
)

var src = rand.NewSource(time.Now().UnixNano())

func GenerateToken() string {
	b := make([]byte, 32)
	// A src.Int63() generates 63 random bits, enough for letterIdxMax characters!
	for i, cache, remain := 31, src.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = src.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(letterBytes) {
			b[i] = letterBytes[idx]
			i--
		}
		cache >>= letterIdxBits
		remain--
	}

	return string(b)
}

// Clean up expired tokens

func tokenCleaner(quit chan bool) {
	tokenCleanerTicker := time.NewTicker(1 * time.Minute) // Tick every minute
	defer log.Println("Token cleaner has stopped.")
	for {
		select {
		case <-tokenCleanerTicker.C: // Tick tock
			if len(token) > 0 { // Check if there are any tokens
				var tokenID string              // Make a string to store a tokenID
				for tokenRange := range token { // Make a range of all tokens
					tokenID = tokenRange // Set a tokenID from range
					checkToken(tokenID)  // Check if token is valid if not it's deleted
				}
			}
		case <-quit: // Quit channel has been called
			tokenCleanerTicker.Stop() // Stop ticker
			return                    // End tokenCleaner function
		}
	}
}
