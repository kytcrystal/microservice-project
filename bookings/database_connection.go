package bookings

import (
	"log"

	"github.com/jmoiron/sqlx"
)

func ConnectToBookingDatabase() *sqlx.DB {
	db, err := sqlx.Connect("postgres", "user=MicroserviceApp dbname=BookingDB sslmode=disable password=MicroserviceApp host=localhost port=5431")
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
	tx.NamedExec("INSERT INTO apartments (id, apartment_name) VALUES (:id, :apartment_name)", &Apartment{Id: "3cc6f6be-e6ea-479a-a1e7-3fd6cab8ae3f", Apartment_Name: "Rarely Orange"})
	tx.NamedExec("INSERT INTO apartments (id, apartment_name) VALUES (:id, :apartment_name)", &Apartment{Id: "d7675c3b-b97e-45a3-87a8-80b46b4d1162", Apartment_Name: "Often Blue"})
	tx.Commit()
}

func refreshBookingTable(db *sqlx.DB) {
	tx := db.MustBegin()
	tx.MustExec("INSERT INTO bookings (id, apartment_id, user_id, start_date, end_date) VALUES ($1, $2, $3, $4, $5)", "1e0cfb4e-0a11-48c8-a1f5-e82f5587818c", "d7675c3b-b97e-45a3-87a8-80b46b4d1162", "M47730", "2024-02-01", "2024-02-20")
	tx.MustExec("INSERT INTO bookings (id, apartment_id, user_id, start_date, end_date) VALUES ($1, $2, $3, $4, $5)", "c956166e-0fad-456e-8a74-e958500f987f", "d7675c3b-b97e-45a3-87a8-80b46b4d1162", "M47788", "2024-03-01", "2024-03-15")
	tx.Commit()
}
