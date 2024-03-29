package main

import (
	"apartments"
	"log"
)

func main() {
	application, err := apartments.CreateApp()
	if err != nil {
		log.Fatal(err)
	}

	if err := application.StartApp(); err != nil {
		log.Fatal(err)
	}
}
