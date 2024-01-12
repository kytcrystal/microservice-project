package bookings

import (
	"errors"
	"fmt"
	"log"

	"github.com/jmoiron/sqlx"

	"github.com/google/uuid"
	_ "github.com/lib/pq"
)

type Booking struct {
	ID          string `db:"id"`
	ApartmentID string `db:"apartment_id"` 
	UserID      string `db:"user_id"`
	StartDate   string `db:"start_date"`
	EndDate     string `db:"end_date"`
}

var bookingDB *sqlx.DB = ConnectToBookingDatabase()

var bookingSchema = `
DROP TABLE IF EXISTS bookings;

CREATE TABLE IF NOT EXISTS bookings (
	id uuid primary key DEFAULT gen_random_uuid(),
    apartment_id uuid,
    user_id text,
	start_date text,
	end_date text
);`

func CreateBooking(booking Booking) (*Booking, error) {
	log.Printf("[booking-CreateBooking] New booking received %+v\n", booking)
	valid := VerifyBookingDates(booking.StartDate, booking.EndDate)
	if !valid {
		return nil, errors.New("booking dates are not valid")
	}
	exists, err := CheckApartmentExists(booking)
	if err != nil {
		return nil, errors.New("error checking apartments")
	}
	if exists {
		available, err := CheckApartmentAvailable(booking)
		if err != nil {
			return nil, err
		}
		if available {
			booking := SaveBooking(booking)
			return &booking, nil
		}
		return nil, errors.New("apartment is not available")
	}
	return nil, errors.New("apartment does not exist")
}

func SaveBooking(booking Booking) Booking {
	booking.ID = uuid.NewString()
	_, err := bookingDB.NamedExec("INSERT INTO bookings (id, apartment_id, user_id, start_date, end_date) VALUES (:id, :apartment_id, :user_id, :start_date, :end_date)", &booking)
	if err != nil {
		log.Fatalln(err)
	}
	log.Printf("[booking-CancelBooking] booking saved: %v\n", booking)
	return booking
}

func CancelBooking(bookingId string) {
	_, err := bookingDB.Exec("DELETE FROM bookings WHERE id = $1", bookingId)
	if err != nil {
		log.Fatalln(err)
	}
	log.Printf("[booking-CancelBooking] Deleted booking with id: %v\n", bookingId)
}

func ChangeBooking(booking Booking) (*Booking, error) {
	CancelBooking(booking.ID)
	newBooking, err := CreateBooking(booking)
	return newBooking, err
}

func ListAllBookings() []Booking {
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
	return bookingList
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

func VerifyBookingDates(startDate string, endDate string) bool {
	if startDate >= endDate {
		return false
	}
	return true
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
