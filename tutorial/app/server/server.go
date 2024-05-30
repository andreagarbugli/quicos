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

func handleConnection(conn, metricsConn net.Conn) {
	defer func(conn net.Conn) {
		err := conn.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(conn)

	for {
		var number int

		_, err := fmt.Fscanf(conn, "%d\n", &number)
		if err != nil {
			fmt.Println("Error receiving data from client:", err)
			return
		}

		fmt.Printf("Received: %d\n", number)

		metric := Metric{
			Timestamp: time.Now().UnixNano() / int64(time.Millisecond),
			Name:      "ping",
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

		time.Sleep(1 * time.Second) // Wait for 1 second before handling the data received

		// Send stat1 metric
		number = rand.Intn(5) + 1

		metric = Metric{
			Timestamp: time.Now().UnixNano() / int64(time.Millisecond),
			Name:      "stat1",
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

		time.Sleep(1 * time.Second) // Wait for 1 second before handling the data received

		// Send the same number back to the client
		_, err = fmt.Fprintf(conn, "%d\n", number)
		if err != nil {
			fmt.Println("Error sending data to client:", err)
			return
		}

		fmt.Printf("Sent: %d\n", number)
	}
}

func main() {
	if len(os.Args) < 3 {
		fmt.Println("Usage: server <server-port> <metrics-port>")
		return
	}

	serverPort, err := strconv.Atoi(os.Args[1])
	if err != nil {
		fmt.Println("Invalid server port:", err)
		return
	}

	metricsPort, err := strconv.Atoi(os.Args[2])
	if err != nil {
		fmt.Println("Invalid metrics port:", err)
		return
	}

	port := fmt.Sprintf(":%d", serverPort)
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

	listener, err := net.Listen("tcp", port)
	if err != nil {
		fmt.Println("Error starting the server:", err)
		return
	}

	defer func(listener net.Listener) {
		err := listener.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(listener)

	fmt.Println("Server started, waiting for connections...")

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting connection:", err)
			continue
		}

		go handleConnection(conn, metricsConn)
	}
}
