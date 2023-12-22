package main

import (
	"bookings"
	"log"
)

func main() {
	err := bookings.Run()
	if err != nil {
		log.Fatalf("Application terminated with error: %v", err)
	}
}

