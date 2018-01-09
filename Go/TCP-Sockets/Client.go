package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"strconv"
	"time"
)

const (
	Nanos = iota
)

type Message struct {
	Type  int
	BodyS string
}

func main() {
	fmt.Println("Started client!")

	conn, err := net.Dial("tcp", ":3737")
	if err != nil {
		log.Fatal("Connection error", err)
	}

	sending := "Current UNIX Epoch (Nanos): " + strconv.Itoa(int(time.Now().UnixNano()))
	message := &Message{Nanos, sending} // Make a struct to send

	json.NewEncoder(conn).Encode(message) // Encode data and send to the server
	conn.Close()

	fmt.Println("Message sent!")
}
