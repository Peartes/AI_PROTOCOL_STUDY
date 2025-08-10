package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"

	"6.5840/kvsrv"
)

func main() {
	clerk := kvsrv.MakeClerk(nil) // Replace nil with actual server endpoint
	listenOnKeyInput(clerk)
}

func listenOnKeyInput(ck *kvsrv.Clerk) {
	for {
		var input string
		scanner := bufio.NewScanner(os.Stdin)
		fmt.Print("Enter your command: ")

		if scanner.Scan() { // waits for user to press Enter
			input = scanner.Text() // no newline included
		}

		if err := scanner.Err(); err != nil {
			fmt.Println("Error reading:", err)
		}
		if input == "exit" {
			log.Println("Exiting kv server")
			os.Exit(0)
		} else {
			// get the command and operands
			parts := strings.Split(input, " ")
			if len(parts) < 3 || len(parts) > 4 {
				log.Println("Invalid command. Usage: <server> <command> <key> [value]")
				continue
			}
			server := parts[0]
			command := parts[1]
			if command == "append" {
				if len(parts) != 4 {
					log.Println("Invalid command. Usage: <server> append <key> <value>")
					continue
				}
				if parts[2] == "" || parts[3] == "" {
					log.Println("Key and value cannot be empty")
					continue
				}
				key := parts[2]
				value := parts[3]

				ck.AppendReplica(key, value, server)
			} else if command == "get" {
				if len(parts) != 3 {
					log.Println("Invalid command. Usage: <server> get <key>")
					continue
				}
				key := parts[2]
				ck.GetReplica(key, server)
			}
		}
	}
}
