package connection

import (
	"bufio"
	"fmt"
	"net"
	"os"

	ai "github.com/peartes/scrimba/pollyglot/ai"
)

func Client(path, targetLanguage string) error {
	// Connect to the server
	conn, err := net.Dial("unix", path)
	if err != nil {
		return fmt.Errorf("error connecting to server: %v ", err)
	}
	defer conn.Close()
	fmt.Println("Connected to server. Type your messages:")

	// Handle incoming messages from the server
	go receiveMessages(conn, targetLanguage)

	// Read user input and send to server
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		text := scanner.Text()
		fmt.Fprintln(conn, text) // Send message to server
	}
	return nil
}

func receiveMessages(conn net.Conn, targetLanguage string) {
	reader := bufio.NewReader(conn)
	for {
		message, err := reader.ReadString('\n')
		if err != nil {
			fmt.Printf("Disconnected from server because of: %v\n", err)
			os.Exit(0)
		}
		res, err := ai.Translate(message, targetLanguage)
		if err != nil {
			fmt.Println("Error translating message:", err)
			continue
		}
		fmt.Printf("Server: %s\n", res)
	}
}
