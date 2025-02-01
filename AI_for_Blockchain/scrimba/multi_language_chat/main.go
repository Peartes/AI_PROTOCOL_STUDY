package main

import (
	"fmt"
	"os"
	"strings"
)

var buddy string

func init() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run main.go <username_to_chat_with>")
		os.Exit(1)
	}
	buddy = strings.Join([]string{"/tmp/", os.Args[1], ".sock"}, "")
}

func main() {
	err := RunApp(buddy)
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
}
