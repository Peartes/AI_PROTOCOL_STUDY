package main

import "fmt"

func main() {
	err := RunApp()
	if err != nil {
		fmt.Println(err)
	}
}
