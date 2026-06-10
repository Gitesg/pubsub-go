package main

import (
	pubsub "ai-leraning/broker"
	"fmt"
	"net"
	"strings"
)

var ps = pubsub.NewPubSub()

func handleConnection(conn net.Conn) {

	buf := make([]byte, 1024)
	n, err := conn.Read(buf)

	if err != nil {
		fmt.Println("Error reading from connection:", err)
		return
	}
	x := string(buf[:n])
	x = strings.TrimSpace(x)
	parts := strings.SplitN(x, " ", 3)
	if len(parts) < 3 {
		fmt.Println("Invalid message format")
		return
	}
	role := parts[0]
	topic := parts[1]
	msg := parts[2]

	if role == "publisher" {
		defer conn.Close()
		code, err := ps.Publisher(msg, topic)
		if err != nil {
			fmt.Printf("Error publishing message: %v (code: %d)\n", err, code)
		} else {
			fmt.Printf("Message published successfully (code: %d)\n", code)
		}

	}
	if role == "subscriber" {
		ch := ps.Subscribe(topic, conn)
		for msg := range ch {
			_, err := conn.Write([]byte(msg))
			if err != nil {
				conn.Close()
				fmt.Println("Error writing to connection:", err)
				return
			}
		}
	}

}

func main() {
	ln, err := net.Listen("tcp", "localhost:8080")
	fmt.Println("Server is listening on localhost:8080")

	if err != nil {
		panic(err)
	}
	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Println("Error accepting connection:", err)
			continue
		}

		go handleConnection(conn)
	}
}
