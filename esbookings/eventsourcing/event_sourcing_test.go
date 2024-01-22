package eventsourcing

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
)

func TestBookingService(t *testing.T) {

	t.Run("it should load a newly created booking", func(t *testing.T) {
		var repo, err = NewRepository()
		assert.NoError(t, err)

		var apartmentId = uuid.NewString()

		var apartmentCreatedEvent = ApartmentCreatedEvent{
			ApartmentID:   apartmentId,
			ApartmentName: "Test Name",
		}

		err = repo.SaveEvent(context.Background(), apartmentId, apartmentCreatedEvent)
		assert.NoError(t, err)

		var event = BookingCreatedEvent{
			BookingID:   uuid.NewString(),
			ApartmentID: apartmentId,
			UserID:      uuid.NewString(),
			StartDate:   "2023-01-02",
			EndDate:     "2023-01-05",
		}

		err = repo.SaveEvent(context.Background(), apartmentId, event)
		assert.NoError(t, err)

		entity, err := repo.Load(context.Background(), apartmentId)
		assert.NoError(t, err)

		assert.NotNil(t, entity)
		assert.Equal(t, apartmentCreatedEvent.ApartmentID, entity.ApartmentID)
		assert.Equal(t, apartmentCreatedEvent.ApartmentName, entity.Name)

		booking := entity.Bookings[0]
		assert.Equal(t, event.BookingID, booking.BookingID)
		assert.Equal(t, event.ApartmentID, booking.ApartmentID)
		assert.Equal(t, event.UserID, booking.UserID)
		assert.Equal(t, event.StartDate, booking.StartDate)
		assert.Equal(t, event.EndDate, booking.EndDate)
		assert.Equal(t, false, booking.Cancelled)
	})

	t.Run("service should allow to add a new booking", func(t *testing.T) {
		var repo, err = NewRepository()
		assert.NoError(t, err)

		var apartmentId = uuid.NewString()

		var service = Service{
			repo: *repo,
		}

		var bookingId = uuid.NewString()

		err = service.CreateApartment(context.Background(), apartmentId, uuid.NewString())
		assert.NoError(t, err)

		err = service.AddBooking(context.Background(),
			apartmentId,      // apartmentId
			uuid.NewString(), // userId
			bookingId,        // bookingId
			"2023-01-02",     // startDate
			"2023-01-05",     // endDate
		)
		assert.NoError(t, err)

		entity, err := repo.Load(context.Background(), apartmentId)
		assert.NoError(t, err)

		assert.NotNil(t, entity)
		assert.Equal(t, apartmentId, entity.ApartmentID)
		assert.NotEmpty(t, entity.Name)

		assert.Len(t, entity.Bookings, 1)
		booking := entity.Bookings[0]
		assert.Equal(t, bookingId, booking.BookingID)
	})

	t.Run("service should allow to rollback application", func(t *testing.T) {
		var repo, err = NewRepository()
		assert.NoError(t, err)

		var apartmentId = uuid.NewString()

		var service = Service{
			repo: *repo,
		}

		var bookingId = uuid.NewString()

		err = service.CreateApartment(context.Background(), apartmentId, uuid.NewString())
		assert.NoError(t, err)

		err = service.AddBooking(context.Background(),
			apartmentId,      // apartmentId
			uuid.NewString(), // userId
			bookingId,        // bookingId
			"2023-01-02",     // startDate
			"2023-01-05",     // endDate
		)
		assert.NoError(t, err)

		entity, err := repo.Load(context.Background(), apartmentId)
		assert.NoError(t, err)
		assert.Len(t, entity.Bookings, 1)

		err = service.AddBooking(context.Background(),
			apartmentId,      // apartmentId
			uuid.NewString(), // userId
			uuid.NewString(), // bookingId
			"2023-02-02",     // startDate
			"2023-02-05",     // endDate
		)
		assert.NoError(t, err)

		entity, err = repo.Load(context.Background(), apartmentId)
		assert.NoError(t, err)
		assert.Len(t, entity.Bookings, 2)

		err = service.Rollback(context.Background(),
			bookingId,
		)
		assert.NoError(t, err)

		entity, err = repo.Load(context.Background(), apartmentId)
		assert.NoError(t, err)
		assert.Len(t, entity.Bookings, 1)
		assert.Equal(t, bookingId, entity.Bookings[0].BookingID)

	})

	t.Run("service should not allow booking if apartment is not available", func(t *testing.T) {
		var repo, err = NewRepository()
		assert.NoError(t, err)

		var apartmentId = uuid.NewString()

		var service = Service{
			repo: *repo,
		}

		var bookingId = uuid.NewString()

		err = service.CreateApartment(context.Background(), apartmentId, uuid.NewString())
		assert.NoError(t, err)

		err = service.AddBooking(context.Background(),
			apartmentId,      // apartmentId
			uuid.NewString(), // userId
			bookingId,        // bookingId
			"2023-01-02",     // startDate
			"2023-01-05",     // endDate
		)
		assert.NoError(t, err)

		entity, err := repo.Load(context.Background(), apartmentId)
		assert.NoError(t, err)
		assert.Len(t, entity.Bookings, 1)

		err = service.AddBooking(context.Background(),
			apartmentId,      // apartmentId
			uuid.NewString(), // userId
			uuid.NewString(), // bookingId
			"2023-01-01",     // startDate
			"2023-01-04",     // endDate
		)
		assert.Error(t, err, errors.New("apartment not available for booking"))

		entity, _ = repo.Load(context.Background(), apartmentId)
		assert.Len(t, entity.Bookings, 1)
	})

	// t.Run("it should load allow to cancel a booking", func(t *testing.T) {
	// 	var repo, err = NewRepository()
	// 	assert.NoError(t, err)

	// 	var service = Service{
	// 		repo: *repo,
	// 	}

	// 	var event = BookingCreatedEvent{
	// 		ID:          uuid.NewString(),
	// 		ApartmentID: uuid.NewString(),
	// 		UserID:      uuid.NewString(),
	// 		StartDate:   "2023-01-02",
	// 		EndDate:     "2023-01-05",
	// 	}

	// 	err = repo.SaveEvent(context.Background(), event.ID, event)
	// 	assert.NoError(t, err)

	// 	err = service.CancelBooking(context.Background(), event.ID)
	// 	assert.NoError(t, err)

	// 	entity, err := repo.Load(context.Background(), event.ID)
	// 	assert.NoError(t, err)

	// 	assert.NotNil(t, entity)
	// 	assert.Equal(t, event.ID, entity.ID)
	// 	assert.Equal(t, event.ApartmentID, entity.ApartmentID)
	// 	assert.Equal(t, event.UserID, entity.UserID)
	// 	assert.Equal(t, true, entity.Cancelled)
	// })
}
