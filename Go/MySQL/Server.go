package main

/*
This is compiled with a DB.go file.
Inside the DB.go file is information about connecting to the DB.

var (
	dbType = "mysql"
	dbUsername
	dbPassword
	dbProtocol = "unix"
	dbFileLocation = "/var/run/mysqld/mysqld.sock"
	dbDatabase
	dbConnString = dbUsername + ":" + dbPassword + "@" + dbProtocol + "(" + dbFileLocation + ")/" + dbDatabase
)

Replace dbFileLocation with IP:Port for remote connections.
*/

import (
	"log"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	db, err := sql.Open(dbType, dbConnString) // Connect to DB

	if err != nil {
		log.Fatal(err)
	}

	// Prepare variables for query
	var (
		id = 1
		username, email string
	)

	rows, err := db.Query("select username, email from users where id = ?", id) // Query DB for username and email from id
	if err != nil {
		log.Fatal(err)
	}

	defer rows.Close()

	for rows.Next() {
		err := rows.Scan(&username, &email) // Scan data from query
		if err != nil {
			log.Fatal(err)
		}
		log.Printf("Queried user with ID: %v. Username: %v, Email: %v.", id, username, email) // Print query results
	}

	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}

	defer db.Close() // Close connection to DB
}