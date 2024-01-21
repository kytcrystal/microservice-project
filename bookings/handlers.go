package bookings

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

func (a *Application) bookingsHandler(r *http.Request) (any, error) {
	log.Println("[bookingsHandler] received request: ", r.Method, r.URL.Path)
	switch r.Method {
	case http.MethodGet:
		return ListAllBookings()
	case http.MethodPost:
		return a.createBooking(r)
	case http.MethodDelete:
		return a.deleteBooking(r)
	case http.MethodPatch:
		return a.updateBooking(r)
	default:
		return nil, fmt.Errorf("method not allowed")
	}
}

func (a *Application) createBooking(r *http.Request) (any, error) {
	var booking Booking
	err := json.NewDecoder(r.Body).Decode(&booking)
	if err != nil {
		return nil, err
	}
	bookingCreated, err := CreateBooking(booking)
	if err != nil {
		return nil, err
	}
	err = a.publisher.SendMessage(MQ_BOOKING_CREATED_EXCHANGE, BookingCreatedEvent{bookingCreated})
	if err != nil {
		return nil, err
	}
	return bookingCreated, nil
}

func (a *Application) deleteBooking(r *http.Request) (any, error) {
	var body struct{ ID string }
	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		return nil, err
	}
	if err := CancelBooking(body.ID); err != nil {
		return nil, err
	}
	err = a.publisher.SendMessage(MQ_BOOKING_CANCELLED_EXCHANGE, BookingCancelledEvent{body.ID})
	if err != nil {
		return nil, err
	}
	return body, nil
}

func (a *Application) updateBooking(r *http.Request) (any, error) {
	var booking Booking
	err := json.NewDecoder(r.Body).Decode(&booking)
	if err != nil {
		return nil, err
	}
	updatedBooking, err := ChangeBooking(booking)
	if err != nil {
		return nil, err
	}
	err = a.publisher.SendMessage("booking_updated", BookingUpdatedEvent{updatedBooking})
	if err != nil {
		return nil, err
	}
	return updatedBooking, nil
}
