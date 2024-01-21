package main

import (
	"bookings"
	"log"
)

func main() {
	app, err := bookings.NewApplication()
	if err != nil {
		log.Fatalf("Failed to initialize application: %v", err)
	}
	if err := app.Run(); err != nil {
		log.Fatalf("Application terminated with error: %v", err)
	}
}
