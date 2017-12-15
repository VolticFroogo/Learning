package main

import (
    "time"
    "math/rand"
)

// Define tokens

const tokenLifetime = 7200 // Tokens are valid for 2 hours

type Token struct {
    userid, expires int
}

var token map[string]Token // Initialize a map for all tokens

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