package main

import (
	"os"

	"github.com/peartes/scrimba/pollyglot/connection"
)

func RunApp(path, targetLanguage string) error {
	fs, _ := os.Stat(path)
	if fs != nil {
		return connection.Client(path, targetLanguage)
	} else {
		return connection.Server(path, targetLanguage)
	}
}
