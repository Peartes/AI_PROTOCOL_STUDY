package main

import (
	"fmt"
	"os"
	"strings"
)

var buddy string
var Language string

func init() {
	if len(os.Args) != 3 {
		fmt.Println("Usage: go run main.go <username_to_chat_with> <your_language>")
		os.Exit(1)
	}
	buddy = strings.Join([]string{"/tmp/", os.Args[1], ".sock"}, "")
	Language = os.Args[2]
}

func main() {
	err := RunApp(buddy)
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
}
