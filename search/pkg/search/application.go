package search

import (
	"log"
)

type Application struct {
	// in here we will be able to add "objects" will need to be shared across our application
	// for instance shared database connection, or similar things;
	// Having it here rather than it in main make it more easy to test, for instance in a unit test.
}

func NewApplication() Application {
	return Application{}
}

func (a Application) Run() error {
	log.Println("Starting Search Application")
	return nil
}
