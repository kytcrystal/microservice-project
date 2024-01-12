package bookings

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
}

var apartmentDB *sqlx.DB = ConnectToApartmentDatabase()

var apartmentSchema = `
DROP TABLE apartment;

CREATE TABLE apartment (
	id uuid primary key DEFAULT gen_random_uuid(),
    apartment_name text,
    address text,
	noise_level text,
	floor text
);`

func RefreshApartment() {

}

func SaveApartment(apartment Apartment) Apartment {
	apartment.Id = uuid.NewString()
	_, err := apartmentDB.NamedExec("INSERT INTO apartment (id, apartment_name, address, noise_level, floor) VALUES (:id, :apartment_name, :address, :noise_level, :floor)", &apartment)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Printf("Apartment added: %v\n", apartment)
	return apartment
}

func DeleteApartment(apartmentId string) {
	_, err := apartmentDB.Exec("DELETE FROM apartment WHERE id = $1", apartmentId)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Printf("Deleted apartment with id: %v\n", apartmentId)
}

func ListAllApartments() []Apartment {
	apartment := Apartment{}
	var apartmentList []Apartment

	rows, _ := apartmentDB.Queryx("SELECT * FROM apartment")

	for rows.Next() {
		err := rows.StructScan(&apartment)
		if err != nil {
			log.Fatalln(err)
		}
		apartmentList = append(apartmentList, apartment)
	}
	return apartmentList
}

func ConnectToApartmentDatabase() *sqlx.DB {
	db, err := sqlx.Connect("postgres", "user=MicroserviceApp dbname=ApartmentDB sslmode=disable password=MicroserviceApp host=localhost")
	if err != nil {
		log.Fatalln("[apartments_repository] Failed to connect to database", err)
	}

	db.MustExec(apartmentSchema)

	tx := db.MustBegin()
	tx.MustExec("INSERT INTO apartment (apartment_name, address, noise_level, floor) VALUES ($1, $2, $3, $4)", "Always Green", "Bolzano", "2", "3")
	tx.MustExec("INSERT INTO apartment (apartment_name, address, noise_level, floor) VALUES ($1, $2, $3, $4)", "Rarely Yellow", "Bolzano", "4", "3")
	tx.Commit()

	// defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatal(err)
	} else {
		log.Println("Successfully Connected")
	}
	return db

}
