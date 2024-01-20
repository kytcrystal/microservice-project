package bookings

type BookingCreatedEvent struct {
	Booking *Booking
}

type BookingUpdatedEvent struct {
	Booking *Booking
}

type BookingCancelledEvent struct {
	BookingID string
}
