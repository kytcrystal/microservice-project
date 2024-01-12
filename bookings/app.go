package bookings

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

var listOfBookings []Booking

func Run() error {
	var port = "3001"
	log.Println("[app:run] starting booking service at port", port)

	http.HandleFunc("/api/bookings", bookingsHandler)

	err := http.ListenAndServe(":"+port, nil)
	return err
}

func bookingsHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("[bookingsHandler] received request: ", r.Method, r.URL.Path)

	w.Header().Set("Content-Type", "application/json")

	switch r.Method {

	case http.MethodGet:
		allBookings := ListAllBookings()
		json.NewEncoder(w).Encode(&allBookings)

	case http.MethodPost:
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
		var body struct{ BookingID string }
		err := json.NewDecoder(r.Body).Decode(&body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		CancelBooking(body.BookingID)
		json.NewEncoder(w).Encode(&body)

	case http.MethodPatch:
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
