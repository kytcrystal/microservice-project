package main

import (
	"gateway/gateway"
	"log"
)

func main() {
	var app = gateway.NewApplication()

	if err := app.Run(); err != nil {
		log.Fatalf("Application terminated with error: %v", err)
	}
}
