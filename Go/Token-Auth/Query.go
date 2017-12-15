package main

import (
	"log"
	"time"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
)

func checkLogin(username, password string) (int, bool) {
	db, err := sql.Open(dbType, dbConnString) // Connect to DB
	if err != nil {
		log.Fatal(err)
	}

	rows, err := db.Query("select id, password from users where username = ?", username) // Query DB for id and password from username
	if err != nil {
		log.Fatal(err)
	}

	defer rows.Close()

	var (
		hash string
		id int
	)

	for rows.Next() {
		err := rows.Scan(&id, &hash) // Scan data from query
		if err != nil {
			log.Fatal(err)
		}
	}

	defer db.Close() // Close connection to DB

	valid := checkPasswordHash(password, hash)

	return id, valid
}

func checkToken(tokenID string) bool {
	if (token[tokenID].userid != 0) { // Check if token exists
		if (token[tokenID].expires > int(time.Now().Unix())) { // Check if token has expired
			return true // Token is valid
		}
	}

	return false // Token is invalid
}

func dbGetUsername(id int) string {
	db, err := sql.Open(dbType, dbConnString) // Connect to DB
	if err != nil {
		log.Fatal(err)
	}

	rows, err := db.Query("select username from users where id = ?", id) // Query DB for username from id
	if err != nil {
		log.Fatal(err)
	}

	defer rows.Close()

	var (
		username string
	)

	for rows.Next() {
		err := rows.Scan(&username) // Scan data from query
		if err != nil {
			log.Fatal(err)
		}
	}

	defer db.Close() // Close connection to DB

	return username
}