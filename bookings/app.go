package bookings

import (
	"encoding/json"
	"fmt"
	"net/http"
)

var listOfBookings []Booking

func Run() error {
	http.HandleFunc("/api/bookings", bookingsHandler)

	err := http.ListenAndServe(":3001", nil)
	return err
}

func bookingsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	switch r.Method {

	case http.MethodGet:
		fmt.Printf("got /api/bookings GET request\n")
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
		newBooking, err := CreateBooking(booking)
		if err != nil {
			fmt.Fprintf(w, "error in booking %v", err)
			return
		}
		json.NewEncoder(w).Encode(&newBooking)
	case http.MethodDelete:
		fmt.Printf("got /api/bookings DELETE request\n")
		var body struct{ BookingID string }
		err := json.NewDecoder(r.Body).Decode(&body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		CancelBooking(body.BookingID)
		json.NewEncoder(w).Encode(&body)

	case http.MethodPatch:
		fmt.Printf("got /api/bookings PATCH request\n")
		var booking Booking
		err := json.NewDecoder(r.Body).Decode(&booking)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		ChangeBooking(booking)
		json.NewEncoder(w).Encode(&booking)

	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}
