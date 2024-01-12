package main

import (
	"bookings"
	"log"
)

func main() {
	if err := bookings.Run(); err != nil {
		log.Fatalf("Application terminated with error: %v", err)
	}
}
