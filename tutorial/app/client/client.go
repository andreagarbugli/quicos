package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net"
	"os"
	"strconv"
	"time"
)

type Metric struct {
	Timestamp int64   `json:"timestamp"`
	Name      string  `json:"name"`
	Value     float64 `json:"value"`
}

func main() {
	if len(os.Args) < 4 {
		fmt.Println("Usage: client <server-address> <server-port> <metrics-port>")
		return
	}

	serverAddr := os.Args[1]
	serverPort, err := strconv.Atoi(os.Args[2])
	if err != nil {
		fmt.Println("Invalid server port:", err)
		return
	}

	metricsPort, err := strconv.Atoi(os.Args[3])
	if err != nil {
		fmt.Println("Invalid metrics port:", err)
		return
	}

	server := fmt.Sprintf("%s:%d", serverAddr, serverPort)
	metrics := fmt.Sprintf(":%d", metricsPort)

	metricsListener, err := net.Listen("tcp", metrics)
	if err != nil {
		fmt.Println("Error listening on metrics port:", err)
		return
	}

	fmt.Println("waiting for metrics-client connection...")

	metricsConn, err := metricsListener.Accept()
	if err != nil {
		fmt.Println("Error accepting metrics connection:", err)
		return
	}

	conn, err := net.Dial("tcp", server)
	if err != nil {
		fmt.Println("Error connecting to server:", err)
		return
	}

	defer func(conn net.Conn) {
		err := conn.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(conn)

	rand.Seed(time.Now().UnixNano())

	for {
		// Send stat2 metric
		number := rand.Intn(4) + 6

		metric := Metric{
			Timestamp: time.Now().UnixNano() / int64(time.Millisecond),
			Name:      "stat2",
			Value:     float64(number),
		}

		metricJSON, err := json.Marshal(metric)
		if err != nil {
			fmt.Println("Error marshaling metric:", err)
			return
		}

		_, err = fmt.Fprintf(metricsConn, "%s\n", metricJSON)
		if err != nil {
			fmt.Println("Error sending metric to metrics port:", err)
			return
		}

		fmt.Printf("Sent metric: %s\n", metricJSON)

		time.Sleep(1 * time.Second) // Wait for 1 second before sending the next number

		number = rand.Intn(100) // Generate a random number between 0 and 99

		metric = Metric{
			Timestamp: time.Now().UnixNano() / int64(time.Millisecond),
			Name:      "ping",
			Value:     float64(number),
		}

		metricJSON, err = json.Marshal(metric)
		if err != nil {
			fmt.Println("Error marshaling metric:", err)
			return
		}

		_, err = fmt.Fprintf(metricsConn, "%s\n", metricJSON)
		if err != nil {
			fmt.Println("Error sending metric to metrics port:", err)
			return
		}

		fmt.Printf("Sent metric: %s\n", metricJSON)

		_, err = fmt.Fprintf(conn, "%d\n", number)
		if err != nil {
			fmt.Println("Error sending data to server:", err)
			return
		}

		fmt.Printf("Sent: %d\n", number)

		// Wait for server response
		var response int
		_, err = fmt.Fscanf(conn, "%d\n", &response)
		if err != nil {
			fmt.Println("Error receiving data from server:", err)
			return
		}

		fmt.Printf("Received: %d\n", response)

		time.Sleep(1 * time.Second) // Wait for 1 second before sending the next number
	}
}
