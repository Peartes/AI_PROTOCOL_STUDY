package main

import "github.com/peartes/scrimba/embeddings/app"

func main() {
	err := app.RunApp()

	if err != nil {
		panic(err)
	}
}
