package main

import (
	"log"
)

func checkLogin(username, password string) (int, bool) {
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

	valid := checkPasswordHash(password, hash)

	return id, valid
}

func dbGetUsername(id int) string {
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

	return username
}