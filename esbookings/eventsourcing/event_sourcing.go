package eventsourcing

import (
	"slices"

	_ "github.com/lib/pq"
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
	b := &ApartmentEntity{}

	for _, event := range events {
		b.On(event, false)
	}
	return b
}

func (b *ApartmentEntity) CreateBooking(bookingID string, apartmentID string, userID string, startDate string, endDate string) error {
	// check if apartment is available if yes emit the event otherwise raise error
	booking := BookingCreatedEvent{
		BookingID:   bookingID,
		ApartmentID: apartmentID,
		UserID:      userID,
		StartDate:   startDate,
		EndDate:     endDate,
	}
	b.raise(booking)
	return nil
}

func (b *ApartmentEntity) Update(bookingId string, startDate string, endDate string) error {
	// check if apartment is available if yes emit the event otherwise raise error

	// if b.Cancelled {
	// 	return fmt.Errorf("booking %v is already cancelled", b.ID)
	// }

	b.raise(BookingUpdatedEvent{
		BookingID: bookingId,
		StartDate: startDate,
		EndDate:   endDate,
	})
	return nil
}

func (b *ApartmentEntity) Cancel(bookingId string) error {
	// if b.Cancelled {
	// 	return fmt.Errorf("booking %v is already cancelled", b.ID)
	// }

	b.raise(BookingCancelledEvent{
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
func (b ApartmentEntity) Events() []Event {
	return b.changes
}

// Version returns the last version of the aggregate before changes.
func (b ApartmentEntity) Version() int {
	return b.version
}

func (b *ApartmentEntity) raise(event Event) {
	b.changes = append(b.changes, event)
	b.On(event, true)
}
