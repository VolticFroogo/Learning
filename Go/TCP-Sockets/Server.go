package main

import (
	"encoding/json"
	"fmt"
	"net"
)

const (
	Nanos = iota
)

type Message struct {
	Type  int    `json:"type"`
	BodyS string `json:"bodys"`
}

func handleConnection(conn net.Conn) {
	message := &Message{}                  // Create a Message to store information
	json.NewDecoder(conn).Decode(&message) // Make a new JSON decoder and decode data

	switch message.Type { // What type of message is being sent?
	case Nanos: // Type is nano
		fmt.Printf("Body: %v.\n", message.BodyS) // Show the message (in this case it's a single string)
	}

	conn.Close() // Close the connection with the client
}

func main() {
	fmt.Println("Started listening!")

	ln, err := net.Listen("tcp", ":3737")
	if err != nil {
		fmt.Println("Starting server error: ", err)
	}

	for {
		conn, err := ln.Accept() // This blocks until connection or error
		fmt.Println("Accepted new client!")
		if err != nil {
			fmt.Println("Accepting client error: ", err)
		}

		go handleConnection(conn) // A goroutine handles conn so that the loop can accept other connections
	}
}
