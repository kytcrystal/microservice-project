package apartments

import (
	"errors"
	"log"
	"net/http"
)

type Application struct {
	repo      *ApartmentRepository
	publisher Publisher
	*http.ServeMux
}

func CreateApp() (*Application, error) {
	repo, err := NewApartmentRepository()
	if err != nil {
		log.Println("[apartments_repository] Failed to connect to database", err)
		return nil, err
	}

	apartmentPublisher, err := NewPublisher(
		MQ_CONNECTION_STRING, 
		MQ_APARTMENT_CREATED_EXCHANGE,
		MQ_APARTMENT_DELETED_EXCHANGE,
	)
	if err != nil {
		apartmentPublisher = &RetryPublisher{}
	}

	return &Application{
		repo:      repo,
		publisher: apartmentPublisher,
		ServeMux:  http.DefaultServeMux,
	}, nil
}

func (a *Application) StartApp() error {
	a.CustomHandleFunc("/api/apartments", a.apartmentsHandler)

	err := http.ListenAndServe(":3000", a.ServeMux)
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		return err
	}
	return nil
}
