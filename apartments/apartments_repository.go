package apartments

import (
	"fmt"
	"log"

	"github.com/jmoiron/sqlx"

	"github.com/google/uuid"
	_ "github.com/lib/pq"
)

type Apartment struct {
	Id             string `db:"id" json:"id"`
	ApartmentName string `db:"apartment_name" json:"apartment_name"`
	Address        string `db:"address" json:"address"`
	NoiseLevel    string `db:"noise_level" json:"noise_level"`
	Floor          string `db:"floor" json:"floor"`
}

var db *sqlx.DB = ConnectToDatabase()

var schema = `
DROP TABLE IF EXISTS apartments;

CREATE TABLE IF NOT EXISTS apartments (
	id uuid primary key DEFAULT gen_random_uuid(),
    apartment_name text,
    address text,
	noise_level text,
	floor text
);`

func SaveApartment(apartment Apartment) Apartment {
	apartment.Id = uuid.NewString()
	_, err := db.NamedExec("INSERT INTO apartments (id, apartment_name, address, noise_level, floor) VALUES (:id, :apartment_name, :address, :noise_level, :floor)", &apartment)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Printf("Apartment added: %v\n", apartment)
	return apartment
}

func DeleteApartment(apartmentId string) {
	_, err := db.Exec("DELETE FROM apartments WHERE id = $1", apartmentId)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Printf("Deleted apartment with id: %v\n", apartmentId)
}

func ListAllApartments() []Apartment {
	apartment := Apartment{}
	var apartmentList []Apartment

	rows, _ := db.Queryx("SELECT * FROM apartments")

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
	connectionString := fmt.Sprintf(
		"user=MicroserviceApp dbname=ApartmentDB sslmode=disable password=MicroserviceApp host=%s port=%s",
		POSTGRES_HOST,
		POSTGRES_PORT,
	)

	db, err := sqlx.Connect("postgres", connectionString)
	
	if err != nil {
		log.Fatalln("[apartments_repository] Failed to connect to database", err)
	}

	db.MustExec(schema)

	tx := db.MustBegin()
	tx.MustExec("INSERT INTO apartments (apartment_name, address, noise_level, floor) VALUES ($1, $2, $3, $4)", "Always Green", "Bolzano", "2", "3")
	tx.MustExec("INSERT INTO apartments (apartment_name, address, noise_level, floor) VALUES ($1, $2, $3, $4)", "Rarely Yellow", "Bolzano", "4", "3")
	tx.NamedExec("INSERT INTO apartments (apartment_name, address, noise_level, floor) VALUES (:apartment_name, :address, :noise_level, :floor)", &Apartment{"0", "Sometimes Pink", "Merano", "1", "5"})
	tx.Commit()

	if err := db.Ping(); err != nil {
		log.Fatal(err)
	} else {
		log.Println("Successfully Connected")
	}
	return db
}
