package main

import "os"

func RunApp(path string) error {
	fs, _ := os.Stat(path)
	if fs != nil {
		return Client(path)
	} else {
		return Server(path)
	}
}
