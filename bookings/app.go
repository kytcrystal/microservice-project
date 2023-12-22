package bookings

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
)

type Booking struct {
	BookingID   string
	ApartmentID string
	UserID      string
	StartDate   string
	EndDate     string
}

type Apartment struct {
	Id             string
	Apartment_Name string
	Address        string
	Noise_level    string
	Floor          string
}

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
		listOfBookings = CancelBooking(body.BookingID)
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

func CreateBooking(booking Booking) (*Booking, error) {
	fmt.Printf("New booking received %v", booking)
	exists, err := CheckApartmentExists(booking)
	if err != nil {
		return nil, err
	}
	if exists {
		listOfBookings = append(listOfBookings, booking)
		return &booking, nil
	}
	return nil, errors.New("apartment does not exist")
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

func ChangeBooking(booking Booking) ([]Booking, error) {
	listOfBookings = CancelBooking(booking.BookingID)
	newBooking, err := CreateBooking(booking)
	if err != nil {
		return listOfBookings, err
	}
	return append(listOfBookings, *newBooking), nil
}

func CheckApartmentExists(booking Booking) (bool, error) {
	response, err := http.Get("http://localhost:3000/api/apartments")
	if err != nil {
		return false, fmt.Errorf("fail to connect: %w", err)
	}
	apartments, err := io.ReadAll(response.Body)
	if err != nil {
		return false, fmt.Errorf("fail to read body: %w", err)
	}
	var apartmentList []Apartment
	err = json.Unmarshal(apartments, &apartmentList)
	if err != nil {
		return false, fmt.Errorf("fail to unmarshal apartment list: %w", err)
	}
	for _, apt := range apartmentList {

		if apt.Id == booking.ApartmentID {
			return true, nil
		}
	}
	return false, nil
}
