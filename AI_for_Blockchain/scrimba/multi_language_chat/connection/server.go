package connection

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"

	ai "github.com/peartes/scrimba/pollyglot/ai"
)

func Server(path, targetLanguage string) error {
	// Start listening on a TCP port
	listener, err := net.Listen("unix", path)
	if err != nil {
		return Client(path, targetLanguage)
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
		go handleConnection(conn, targetLanguage)
	}
}

func handleConnection(conn net.Conn, targetLanguage string) {
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
		res, err := ai.Translate(message, targetLanguage)
		if err != nil {
			fmt.Println("Error translating message:", err)
			continue
		}
		fmt.Printf("Client: %s\n", res)
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
