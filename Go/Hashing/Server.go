package main

import (
	"log"
	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func main() {
	var (
		saved = "My-Super-Secret-Password" // Password saved in DB
		input = "My-Super-Secret-Password" // Password just inputted
	)

	hash, err := HashPassword(saved) // Hash saved password (already hashed in real DB)
	if (err != nil) {
		log.Fatalf("Error hashing password: %v.", err)
	}

	// Print info
	log.Println("Input Pass: ", input)
	log.Println("Saved Pass: ", saved)
	log.Println("Saved Hash: ", hash)

	match := CheckPasswordHash(input, hash) // Compare hash with input
	log.Println("Match:      ", match)
}