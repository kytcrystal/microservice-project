package eventsourcing

import (
	"slices"

	_ "github.com/lib/pq"
	"github.com/pkg/errors"
)

// Event is a domain event marker.
type Event interface {
	isEvent()
}

type ApartmentCreatedEvent struct {
	ApartmentID   string `json:"apartment_id"`
	ApartmentName string `json:"name"`
}

type BookingCreatedEvent struct {
	BookingID   string `json:"id"`
	ApartmentID string `json:"apartment_id"`
	UserID      string `json:"user_id"`
	StartDate   string `json:"start_date"`
	EndDate     string `json:"end_date"`
}

type BookingUpdatedEvent struct {
	BookingID string `json:"id"`
	StartDate string `json:"start_date"`
	EndDate   string `json:"end_date"`
}

type BookingCancelledEvent struct {
	BookingID string `json:"id"`
}

func (e ApartmentCreatedEvent) isEvent() {}
func (e BookingCreatedEvent) isEvent()   {}
func (e BookingCancelledEvent) isEvent() {}
func (e BookingUpdatedEvent) isEvent()   {}

type ApartmentEntity struct {
	ApartmentID string
	Name        string
	Bookings    []BookingEntity

	changes []Event
	version int
}

type BookingEntity struct {
	BookingID   string
	ApartmentID string
	UserID      string
	StartDate   string
	EndDate     string
	Cancelled   bool
}

// NewFromEvents is a helper method that creates a new ApartmentEntity
// from a series of events.
func NewFromEvents(events []Event) *ApartmentEntity {
	a := &ApartmentEntity{}

	for _, event := range events {
		a.On(event, false)
	}
	return a
}

func (a *ApartmentEntity) CreateBooking(bookingID string, apartmentID string, userID string, startDate string, endDate string) error {

	if (startDate > endDate) {
		return errors.New("booking date range is invalid")
	}

	if err := a.CheckAvailable(bookingID, apartmentID, startDate, endDate); err != nil {
		return err
	}

	booking := BookingCreatedEvent{
		BookingID:   bookingID,
		ApartmentID: apartmentID,
		UserID:      userID,
		StartDate:   startDate,
		EndDate:     endDate,
	}
	a.raise(booking)
	return nil
}

func (a *ApartmentEntity) CheckAvailable(bookingID string, apartmentID string, startDate string, endDate string) error {

	for _, booking := range a.Bookings {
		if !booking.Cancelled {
			if startDate >= booking.StartDate && endDate <= booking.EndDate {
				return errors.New("apartment not available for booking")
			}
			if startDate >= booking.StartDate && startDate <= booking.EndDate {
				return errors.New("apartment not available for booking")
			}
			if endDate >= booking.StartDate && endDate <= booking.EndDate {
				return errors.New("apartment not available for booking")
			}
			if startDate <= booking.StartDate && endDate >= booking.EndDate {
				return errors.New("apartment not available for booking")
			}
		}
	}
	return nil
}

func (a *ApartmentEntity) Update(bookingId string, startDate string, endDate string) error {
	// check if apartment is available if yes emit the event otherwise raise error

	// if b.Cancelled {
	// 	return fmt.Errorf("booking %v is already cancelled", b.ID)
	// }

	a.raise(BookingUpdatedEvent{
		BookingID: bookingId,
		StartDate: startDate,
		EndDate:   endDate,
	})
	return nil
}

func (a *ApartmentEntity) Cancel(bookingId string) error {
	// if b.Cancelled {
	// 	return fmt.Errorf("booking %v is already cancelled", b.ID)
	// }

	a.raise(BookingCancelledEvent{
		BookingID: bookingId,
	})
	return nil
}

// On handles ApartmentEntity events on the Booking aggregate.
func (a *ApartmentEntity) On(event Event, new bool) {
	switch e := event.(type) {
	case ApartmentCreatedEvent:
		a.ApartmentID = e.ApartmentID
		a.Name = e.ApartmentName

	case BookingCreatedEvent:

		var b BookingEntity
		b.BookingID = e.BookingID
		b.ApartmentID = e.ApartmentID
		b.UserID = e.UserID
		b.StartDate = e.StartDate
		b.EndDate = e.EndDate
		b.Cancelled = false

		a.Bookings = append(a.Bookings, b)

	case BookingUpdatedEvent:
		index := slices.IndexFunc(a.Bookings, func(be BookingEntity) bool {
			return be.BookingID == e.BookingID
		})
		a.Bookings[index].StartDate = e.StartDate
		a.Bookings[index].EndDate = e.EndDate

	case BookingCancelledEvent:
		// b.Cancelled = true
		// remove or set cancelled flag
	}
	if !new {
		a.version++
	}
}

// Events returns the uncommitted events from the Booking aggregate.
func (a ApartmentEntity) Events() []Event {
	return a.changes
}

// Version returns the last version of the aggregate before changes.
func (a ApartmentEntity) Version() int {
	return a.version
}

func (a *ApartmentEntity) raise(event Event) {
	a.changes = append(a.changes, event)
	a.On(event, true)
}
