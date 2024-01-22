package eventsourcing

import (
	"context"
	"log"
)

type Service struct {
	repo EventRepository
}

func NewService() (*Service, error) {
	repo, err := NewRepository()
	if err != nil {
		return nil, err
	}
	return &Service{
		repo: *repo,
	}, nil
}

func (s *Service) Rollback(
	ctx context.Context,
	bookingID string,
) error {
	_, err := s.repo.db.Exec(`DELETE FROM booking_events 
	where sequence_nr > (
		select sequence_nr 
			from booking_events be 
		where be.payload ->> 'id' = $1
		order by be.created_at desc limit 1 
	)`, bookingID)
	return err
}

func (s *Service) CreateApartment(
	ctx context.Context,
	apartmentID string,
	apartmentName string,
) error {
	entity, err := s.repo.Load(ctx, apartmentID)
	if err != nil {
		return err
	}
	if entity == nil || entity.ApartmentID == "" {
		return s.repo.SaveEvent(ctx, apartmentID, ApartmentCreatedEvent{
			ApartmentID:   apartmentID,
			ApartmentName: apartmentName,
		})
	}
	log.Println("entity already exist: skip insert", apartmentID)
	return nil
}

func (s *Service) DeleteApartment(
	ctx context.Context,
	apartmentID string,
) error {
	_, err := s.repo.Load(ctx, apartmentID)
	if err != nil {
		return err
	}
	return nil
}

func (s *Service) CancelBooking(
	ctx context.Context,
	apartmentId string,
	bookingId string,
) error {
	return nil
}

func (s *Service) AddBooking(
	ctx context.Context,
	apartmentId string,
	userId string,
	bookingId string,
	startDate string,
	endDate string,
) error {
	b, err := s.repo.Load(ctx, apartmentId)
	if err != nil {
		return err
	}

	if err := b.CreateBooking(bookingId, apartmentId, userId, startDate, endDate); err != nil {
		return err
	}

	if err := s.repo.Save(ctx, b); err != nil {
		return err
	}
	return nil
}

func (s *Service) UpdateBooking(
	ctx context.Context,
	apartmentId string,
	bookingId string,
	startDate string,
	endDate string,
) error {
	b, err := s.repo.Load(ctx, apartmentId)
	if err != nil {
		return err
	}

	if err := b.Update(bookingId, startDate, endDate); err != nil {
		return err
	}

	if err := s.repo.Save(ctx, b); err != nil {
		return err
	}
	return nil
}
