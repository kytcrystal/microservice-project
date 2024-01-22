package esbookings

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

func (a *Application) rollbackHandler(r *http.Request) (any, error) {
	log.Println("[rollbackHandler] received request: ", r.Method, r.URL.Path)
	switch r.Method {
	case http.MethodPost:
		return a.rollback(r)
	default:
		return nil, fmt.Errorf("method not allowed")
	}
}

func (a *Application) rollback(r *http.Request) (any, error) {
	var body struct {
		BookingID string
	}
	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		return nil, err
	}
	err = a.service.Rollback(context.Background(), body.BookingID)
	if err != nil {
		return nil, err
	}
	return "OK", nil
}

func (a *Application) bookingsHandler(r *http.Request) (any, error) {
	log.Println("[bookingsHandler] received request: ", r.Method, r.URL.Path)
	switch r.Method {
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

	err = a.service.AddBooking(context.Background(),
		booking.ApartmentID,
		booking.UserID,
		booking.ID,
		booking.StartDate,
		booking.EndDate,
	)
	if err != nil {
		return nil, err
	}
	err = a.publisher.SendMessage(MQ_BOOKING_CREATED_EXCHANGE, booking)
	if err != nil {
		return nil, err
	}
	return booking, nil
}

func (a *Application) deleteBooking(r *http.Request) (any, error) {
	var body struct {
		ID          string `json:"id"`
		ApartmentID string `json:"apartment_id"`
	}
	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		return nil, err
	}
	if body.ApartmentID == "" {
		return nil, fmt.Errorf("ApartmentID is required: %+v", body)
	}

	err = a.service.CancelBooking(context.Background(),
		body.ApartmentID,
		body.ID,
	)
	if err != nil {
		return nil, err
	}
	err = a.publisher.SendMessage(MQ_BOOKING_CANCELLED_EXCHANGE, body)
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

	err = a.service.UpdateBooking(context.Background(),
		booking.ApartmentID,
		booking.ID,
		booking.StartDate,
		booking.EndDate,
	)
	if err != nil {
		return nil, err
	}

	err = a.publisher.SendMessage("booking_updated", booking)
	if err != nil {
		return nil, err
	}
	return booking, nil
}
