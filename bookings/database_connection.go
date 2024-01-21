package bookings

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/jmoiron/sqlx"
)

func ConnectToBookingDatabase() *sqlx.DB {
	connectionString := fmt.Sprintf(
		"user=MicroserviceApp dbname=BookingDB sslmode=disable password=MicroserviceApp host=%s port=%s",
		POSTGRES_HOST,
		POSTGRES_PORT,
	)

	db, err := sqlx.Connect("postgres", connectionString)
	if err != nil {
		log.Fatalln("[booking:connect_to_booking_database] Failed to connect to database", err)
	}

	log.Println("[booking:connect_to_booking_database] starting to set up database schema")

	db.MustExec(apartmentSchema)
	refreshApartmentTable(db)

	db.MustExec(bookingSchema)
	refreshBookingTable(db)

	log.Println("[booking:connect_to_booking_database] database schema set up without errors")

	if err := db.Ping(); err != nil {
		log.Fatal(err)
	} else {
		log.Println("Successfully Connected")
	}
	return db

}

func refreshApartmentTable(db *sqlx.DB) {
	tx := db.MustBegin()

	response, err := http.Get(APARTMENT_URL + "/api/apartments")
	if err != nil {
		log.Fatalf("fail to connect: %w", err)
	}
	var apartmentList []Apartment
	if err = json.NewDecoder(response.Body).Decode(&apartmentList); err != nil {
		log.Fatalf("fail to unmarshal apartment list: %w", err)
	}

	for _, apt := range apartmentList {
		tx.NamedExec("INSERT INTO apartments (id, apartment_name) VALUES (:id, :apartment_name)", &Apartment{Id: apt.Id, Apartment_Name: apt.Apartment_Name})
	}
	tx.Commit()
	log.Println("Initialized database with the following apartments", apartmentList)
}

func refreshBookingTable(db *sqlx.DB) {
	tx := db.MustBegin()
	tx.MustExec("INSERT INTO bookings (id, apartment_id, user_id, start_date, end_date) VALUES ($1, $2, $3, $4, $5)", "1e0cfb4e-0a11-48c8-a1f5-e82f5587818c", "d7675c3b-b97e-45a3-87a8-80b46b4d1162", "M47730", "2024-02-01", "2024-02-20")
	tx.MustExec("INSERT INTO bookings (id, apartment_id, user_id, start_date, end_date) VALUES ($1, $2, $3, $4, $5)", "c956166e-0fad-456e-8a74-e958500f987f", "d7675c3b-b97e-45a3-87a8-80b46b4d1162", "M47788", "2024-03-01", "2024-03-15")
	tx.Commit()
}
