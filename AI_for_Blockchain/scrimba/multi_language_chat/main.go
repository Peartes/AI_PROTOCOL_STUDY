package main

import (
	"fmt"
	"os"
	"strings"
)

var client string
var targetLanguage string

func init() {
	if len(os.Args) != 3 {
		fmt.Println("Usage: go run main.go <username_to_chat_with> <your_language>")
		os.Exit(1)
	}
	client = strings.Join([]string{"/tmp/", os.Args[1], ".sock"}, "")
	targetLanguage = os.Args[2]
}

func main() {
	err := RunApp(client, targetLanguage)
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
}
