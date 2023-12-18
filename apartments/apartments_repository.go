package apartments

import (
	"fmt"
	"log"

	"github.com/jmoiron/sqlx"

	"github.com/google/uuid"
	_ "github.com/lib/pq"
)

type Apartment struct {
	Id             string
	Apartment_Name string
	Address        string
	Noise_level    string
	Floor          string
}

var db *sqlx.DB = ConnectToDatabase()

var schema = `
DROP TABLE apartment;

CREATE TABLE apartment (
	id uuid primary key DEFAULT gen_random_uuid(),
    apartment_name text,
    address text,
	noise_level text,
	floor text
);`

func SaveApartment(apartment Apartment) Apartment {
	apartment.Id = uuid.NewString()
	_, err := db.NamedExec("INSERT INTO apartment (id, apartment_name, address, noise_level, floor) VALUES (:id, :apartment_name, :address, :noise_level, :floor)", &apartment)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Printf("Apartment added: %v\n", apartment)
	return apartment
}

func ListAllApartments() []Apartment {
	apartment := Apartment{}
	var apartmentList []Apartment

	rows, _ := db.Queryx("SELECT * FROM apartment")

	for rows.Next() {
		err := rows.StructScan(&apartment)
		if err != nil {
			log.Fatalln(err)
		}
		apartmentList = append(apartmentList, apartment)
	}
	return apartmentList
}

func ConnectToDatabase() *sqlx.DB {
	db, err := sqlx.Connect("postgres", "user=MicroserviceApp dbname=ApartmentDB sslmode=disable password=MicroserviceApp host=localhost")
	if err != nil {
		log.Fatalln(err)
	}

	db.MustExec(schema)

	tx := db.MustBegin()
	tx.MustExec("INSERT INTO apartment (apartment_name, address, noise_level, floor) VALUES ($1, $2, $3, $4)", "Always Green", "Bolzano", "2", "3")
	// Named queries can use structs, so if you have an existing struct (i.e. person := &Person{}) that you have populated, you can pass it in as &person
	// tx.NamedExec("INSERT INTO apartment (apartment_name, address, noise_level, floor) VALUES (:id, :apartment_name, :address, :noise_level, :floor)", &Apartment{"gen_random_uuid()", "Sometimes Pink", "Merano", "1", "5"})
	tx.Commit()

	// defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatal(err)
	} else {
		log.Println("Successfully Connected")
	}
	return db

}
