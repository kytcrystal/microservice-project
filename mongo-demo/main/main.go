package main

import (
	"fmt"
	"log"

	"github.com/kamva/mgm/v3"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	if err := run(); err != nil {
		log.Fatalf("Error while running the application: %v", err)
	}
}

func run() error {
	err := mgm.SetDefaultConfig(nil, "mgm_lab", options.Client().ApplyURI("mongodb://root:example@localhost:27017"))
	if err != nil {
		return fmt.Errorf("failed to setup mongo connection: %w", err)
	}
	book := NewBook("Pride and Prejudice", 345)

	// Make sure to pass the model by reference (to update the model's "updated_at", "created_at" and "id" fields by mgm).
	err = mgm.Coll(book).Create(book)
	if err != nil {
		return fmt.Errorf("failed to create book: %w", err)
	}

	fmt.Println("Inserted new book: ", book)

	// Find and decode the doc to a book model.
	fromDB := &Book{}
	coll := mgm.Coll(fromDB)
	_ = coll.FindByID(book.ID, fromDB)
	fmt.Println("Found in DB: ", fromDB)

	// update the pages
	fromDB.Pages = 123
	err = mgm.Coll(book).Update(fromDB)
	if err != nil {
		return fmt.Errorf("failed to update book: %w", err)
	}

	// Find and decode the doc to a book model to see if it has actually been updated
	fromDB = &Book{}
	_ = coll.FindByID(book.ID, fromDB)
	fmt.Println("Found in DB after the update: ", fromDB)

	// add a step to find all and cleanup before exiting
	var allBooks = []Book{} // carefull need to be inizialised not just delcared ?!
	err = mgm.Coll(&Book{}).SimpleFind(&allBooks, bson.M{})
	if err != nil {
		return fmt.Errorf("failed to find books: %w", err)
	}

	fmt.Println("All Books:", len(allBooks))
	for i, b := range allBooks {
		fmt.Printf("%d: name = %v pages = %d\n", i, b.Name, b.Pages)
		mgm.Coll(&b).Delete(&b)
	}

	return nil
}

type Book struct {
	// DefaultModel adds _id, created_at and updated_at fields to the Model.
	mgm.DefaultModel `bson:",inline"`
	Name             string `json:"name" bson:"name"`
	Pages            int    `json:"pages" bson:"pages"`
}

func NewBook(name string, pages int) *Book {
	return &Book{
		Name:  name,
		Pages: pages,
	}
}
