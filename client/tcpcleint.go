package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
)

func listenToServer(conn net.Conn) {
	reader := bufio.NewReader(conn)
	for {

		message, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("\n Disconnected from broker server.", err)
			return
		}
		fmt.Printf("\n New Message Received: %s", message)
		fmt.Print("Enter command > ")
	}
}

func main() {
	conn, err := net.Dial("tcp", "localhost:8080")
	if err != nil {
		fmt.Println("Error connecting to broker:", err)
		return
	}
	defer conn.Close()
	fmt.Println("Connected to Broker on localhost:8080")

	go listenToServer(conn)

	reader := bufio.NewReader(os.Stdin)
	fmt.Println("\nFormats:\n - subscriber [topic]\n - publisher [topic] [message]")
	fmt.Println("--------------------------------------------------")

	for {
		fmt.Print("Enter command > ")
		input, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Error reading input:", err)
			break
		}

		input = strings.TrimSpace(input)
		if input == "exit" {
			break
		}

		if input == "" {
			continue
		}

		_, err = conn.Write([]byte(input + "\n"))
		if err != nil {
			fmt.Println("Error sending to broker:", err)
			break
		}
	}
}
