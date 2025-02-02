package connection

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
)

func Client(path string) error {
	// Connect to the server
	conn, err := net.Dial("unix", path)
	if err != nil {
		return fmt.Errorf("error connecting to server: %v ", err)
	}
	defer conn.Close()
	fmt.Println("Connected to server. Type your messages:")

	// Handle incoming messages from the server
	go receiveMessages(conn)

	// Read user input and send to server
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		text := scanner.Text()
		fmt.Fprintln(conn, text) // Send message to server
	}
	return nil
}

func receiveMessages(conn net.Conn) {
	reader := bufio.NewReader(conn)
	for {
		message, err := reader.ReadString('\n')
		if err != nil {
			fmt.Printf("Disconnected from server because of: %v\n", err)
			os.Exit(0)
		}
		if strings.HasPrefix(message, "/language") {
			Language = strings.TrimSpace(strings.TrimPrefix(message, "/language"))
			continue
		}
		fmt.Printf("Server: %s", message)
	}
}
