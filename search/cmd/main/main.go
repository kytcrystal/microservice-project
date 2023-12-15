package main

import (
	"log"
	"search/pkg/search"
)

func main() {
	var app = search.NewApplication()

	if err := app.Run(); err != nil {
		log.Fatalf("Application terminated with error: %v", err)
	}
}
