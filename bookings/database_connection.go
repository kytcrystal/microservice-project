package bookings

import (
	"log"

	"github.com/jmoiron/sqlx"
)


func ConnectToBookingDatabase() *sqlx.DB {
	db, err := sqlx.Connect("postgres", "user=MicroserviceApp dbname=BookingDB sslmode=disable password=MicroserviceApp host=localhost")
	if err != nil {
		log.Fatalln("[database_connection] Failed to connect to database", err)
	}

	db.MustExec(apartmentSchema)
	refreshApartmentTable(db)

	db.MustExec(bookingSchema)
	refreshBookingTable(db)

	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatal(err)
	} else {
		log.Println("Successfully Connected")
	}
	return db

}

func refreshApartmentTable(db *sqlx.DB) {
	tx := db.MustBegin()
	tx.MustExec("INSERT INTO apartment (apartment_name, address, noise_level, floor) VALUES ($1, $2, $3, $4)", "Always Green", "Bolzano", "2", "3")
	tx.MustExec("INSERT INTO apartment (apartment_name, address, noise_level, floor) VALUES ($1, $2, $3, $4)", "Rarely Yellow", "Bolzano", "4", "3")
	tx.Commit()
}

func refreshBookingTable(db *sqlx.DB) {
	tx := db.MustBegin()
	tx.MustExec("INSERT INTO booking (id, apartment_id, user_id, start_date, end_date) VALUES ($1, $2, $3, $4, $5)", "1e0cfb4e-0a11-48c8-a1f5-e82f5587818c", "d7675c3b-b97e-45a3-87a8-80b46b4d1162", "M47730", "2024-02-01", "2024-02-20")
	tx.MustExec("INSERT INTO booking (id, apartment_id, user_id, start_date, end_date) VALUES ($1, $2, $3, $4, $5)", "5kshcub4e-0a11-48c8-a1f5-e82f5587818c", "d7675c3b-b97e-45a3-87a8-80b46b4d1162", "M47788", "2024-03-01", "2024-03-15")
	tx.Commit()
}
