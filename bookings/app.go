package bookings

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Booking struct {
	BookingID   string
	ApartmentID string
	UserID      string
	StartDate   string
	EndDate     string
}

var listOfBookings []Booking

func Run() error {
	http.HandleFunc("/api/bookings", bookingsHandler)

	err := http.ListenAndServe(":3001", nil)
	return err
}

func bookingsHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		fmt.Printf("got /api/bookings GET request\n")
		w.Header().Set("Content-Type", "application/json")
		allBookings := ListAllBookings()
		json.NewEncoder(w).Encode(&allBookings)
	case http.MethodPost:
		fmt.Printf("got /api/bookings POST request\n")
		var booking Booking
		err := json.NewDecoder(r.Body).Decode(&booking)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		booking = CreateBooking(booking)
		json.NewEncoder(w).Encode(&booking)
	case http.MethodDelete:
		fmt.Printf("got /api/bookings DELETE request\n")
		var body struct{ BookingID string }
		err := json.NewDecoder(r.Body).Decode(&body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		fmt.Println("here")
		w.Header().Set("Content-Type", "application/json")
		listOfBookings = CancelBooking(body.BookingID)
		json.NewEncoder(w).Encode(&body)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func CreateBooking(booking Booking) Booking {
	fmt.Printf("New booking received %v", booking)
	listOfBookings = append(listOfBookings, booking)
	return booking
}

func ListAllBookings() []Booking {
	return listOfBookings
}

func CancelBooking(bookingID string) []Booking {
	for i, b := range listOfBookings {
		if b.BookingID == bookingID {
			return remove(listOfBookings, i)
		}
	}
	return listOfBookings
}

func remove(listOfBookings []Booking, i int) []Booking {
	return append(listOfBookings[:i], listOfBookings[i+1:]...)
}
