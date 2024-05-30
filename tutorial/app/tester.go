package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"os"
)

type Metric struct {
	Timestamp int64  `json:"timestamp"`
	Name      string `json:"name"`
	Value     float64    `json:"value"`
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: client <client-address> <client-port>")
		return
	}

    clientAddress := os.Args[1]
	clientPort := os.Args[2]
	client := fmt.Sprintf("%s:%s", clientAddress, clientPort)

	conn, err := net.Dial("tcp", client)
	if err != nil {
		fmt.Println("Error connecting to metrics port:", err)
		return
	}

	defer func() {
		err := conn.Close()
		if err != nil {
			log.Fatal(err)
		}
	}()

	for {
		buffer := make([]byte, 1024)
		n, err := conn.Read(buffer)
		if err != nil {
			fmt.Println("Error receiving data from server:", err)
			return
		}

		var metric Metric
		err = json.Unmarshal(buffer[:n], &metric)
		if err != nil {
			fmt.Println("Error decoding JSON:", err)
			return
		}

		fmt.Println("Received metric:")
		fmt.Println("Timestamp:", metric.Timestamp)
		fmt.Println("Name:", metric.Name)
		fmt.Println("Value:", metric.Value)
		fmt.Println()
	}
}
