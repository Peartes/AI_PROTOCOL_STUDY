package connection

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"os/signal"
	"strings"
	"syscall"
)

var Language string // Language of the client you're chatting with

func Server(path string) error {
	// Start listening on a TCP port
	listener, err := net.Listen("unix", path)
	if err != nil {
		return Client(path)
	}
	defer listener.Close()
	fmt.Println("Server started on unix file: ", path)

	handleShutdown(path, listener)
	// Accept incoming connections
	for {
		conn, err := listener.Accept()
		if err != nil {
			continue
		}
		fmt.Println("New client connected")

		// Handle each client in a new goroutine
		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()
	reader := bufio.NewReader(conn)

	go func() { // Read user input and send to server
		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			text := scanner.Text()
			fmt.Fprintln(conn, text) // Send message to server
		}
	}()
	// Continuously read messages from the client
	for {
		message, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Client disconnected.")
			return
		}
		if strings.HasPrefix(message, "/language") {
			Language = strings.TrimSpace(strings.TrimPrefix(message, "/language"))
			continue
		}
		fmt.Printf("Client: %s", message)
	}

}

// handleShutdown removes the socket file on exit
func handleShutdown(socketPath string, listener net.Listener) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-c
		fmt.Println("\nShutting down server...")
		listener.Close()
		os.Remove(socketPath)
		os.Exit(0)
	}()
}
