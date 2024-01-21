package bookings

import (
	"errors"
	"fmt"
	"log"

	"github.com/jmoiron/sqlx"

	_ "github.com/lib/pq"
)

type Booking struct {
	ID          string `db:"id" json:"id"`
	ApartmentID string `db:"apartment_id" json:"apartment_id"`
	UserID      string `db:"user_id" json:"user_id"`
	StartDate   string `db:"start_date" json:"start_date"`
	EndDate     string `db:"end_date" json:"end_date"`
}

var bookingDB *sqlx.DB = ConnectToBookingDatabase()

var bookingSchema = `
DROP TABLE IF EXISTS bookings;

CREATE TABLE IF NOT EXISTS bookings (
	id uuid primary key,
    apartment_id uuid,
    user_id text,
	start_date text,
	end_date text
);`

func CreateBooking(booking Booking) (*Booking, error) {
	log.Printf("[booking-CreateBooking] New booking received %+v\n", booking)

	if booking.StartDate >= booking.EndDate {
		return nil, errors.New("booking dates are not valid")
	}

	exists, err := CheckApartmentExists(booking)
	if err != nil {
		return nil, errors.New("error checking apartments")
	}
	if !exists {
		return nil, errors.New("apartment does not exist")

	}

	available, err := CheckApartmentAvailable(booking)
	if err != nil {
		return nil, err
	}
	if !available {
		return nil, errors.New("apartment is not available")
	}

	bookingCreated := SaveBooking(booking)
	return &bookingCreated, nil

}

func SaveBooking(booking Booking) Booking {
	_, err := bookingDB.NamedExec("INSERT INTO bookings (id, apartment_id, user_id, start_date, end_date) VALUES (:id, :apartment_id, :user_id, :start_date, :end_date)", &booking)
	if err != nil {
		log.Fatalln(err)
	}
	log.Printf("[booking-SaveBooking] booking saved: %v\n", booking)
	return booking
}

func CancelBooking(bookingId string) error {
	_, err := bookingDB.Exec("DELETE FROM bookings WHERE id = $1", bookingId)
	if err != nil {
		return err
	}
	log.Printf("[booking-CancelBooking] Deleted booking with id: %v\n", bookingId)
	return nil
}

func ChangeBooking(booking Booking) (*Booking, error) {
	err := CancelBooking(booking.ID)
	if err != nil {
		return nil, err
	}
	newBooking, err := CreateBooking(booking)
	return newBooking, err
}

func ListAllBookings() ([]Booking, error) {
	booking := Booking{}
	var bookingList []Booking

	rows, err := bookingDB.Queryx("SELECT * FROM bookings")
	if err != nil {
		log.Fatal(fmt.Errorf("failed to read bookings from db: %w", err))
	}

	for rows.Next() {
		err := rows.StructScan(&booking)
		if err != nil {
			log.Fatalln(err)
		}
		bookingList = append(bookingList, booking)
	}
	return bookingList, nil
}

func CheckApartmentExists(booking Booking) (bool, error) {
	apartmentList := ListAllApartments()

	for _, apt := range apartmentList {
		if apt.Id == booking.ApartmentID {
			return true, nil
		}
	}
	return false, nil
}

func CheckApartmentAvailable(newBooking Booking) (bool, error) {
	rows, _ := bookingDB.Queryx("SELECT * FROM bookings WHERE apartment_id = $1", newBooking.ApartmentID)
	existingBooking := Booking{}
	for rows.Next() {
		err := rows.StructScan(&existingBooking)
		if err != nil {
			log.Fatalln(err)
		}
		if !DatesAvailable(newBooking, existingBooking) {
			return false, errors.New("apartment not available")
		}
	}
	return true, nil
}

func DatesAvailable(newBooking Booking, existingBooking Booking) bool {
	if newBooking.EndDate <= existingBooking.StartDate {
		return true
	}
	if newBooking.StartDate >= existingBooking.EndDate {
		return true
	}
	return false
}
