package bookings

import (
	"fmt"
	"log"

	"github.com/jmoiron/sqlx"

	_ "github.com/lib/pq"
)

type Apartment struct {
	Id             string
	Apartment_Name string
}

var apartmentDB *sqlx.DB = ConnectToBookingDatabase()

var apartmentSchema = `
DROP TABLE IF EXISTS apartments;

CREATE TABLE IF NOT EXISTS apartments (
	id uuid primary key,
    apartment_name text
);`

func SaveApartment(apartment Apartment) Apartment {
	_, err := apartmentDB.NamedExec("INSERT INTO apartments (id, apartment_name) VALUES (:id, :apartment_name)", &apartment)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Printf("Apartment added: %v\n", apartment)
	return apartment
}

func DeleteApartment(apartmentId string) {
	_, err := apartmentDB.Exec("DELETE FROM apartments WHERE id = $1", apartmentId)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Printf("Deleted apartment with id: %v\n", apartmentId)
}

func ListAllApartments() []Apartment {
	apartment := Apartment{}
	var apartmentList []Apartment

	rows, _ := apartmentDB.Queryx("SELECT * FROM apartments")

	for rows.Next() {
		err := rows.StructScan(&apartment)
		if err != nil {
			log.Fatalln(err)
		}
		apartmentList = append(apartmentList, apartment)
	}
	return apartmentList
}
