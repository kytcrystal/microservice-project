package esbookings

import (
	_ "github.com/lib/pq"
)

type Booking struct {
	ID          string `db:"id" json:"id"`
	ApartmentID string `db:"apartment_id" json:"apartment_id"`
	UserID      string `db:"user_id" json:"user_id"`
	StartDate   string `db:"start_date" json:"start_date"`
	EndDate     string `db:"end_date" json:"end_date"`
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
