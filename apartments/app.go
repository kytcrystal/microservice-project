package apartments

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

func getRoot(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("got / request\n")
	io.WriteString(w, "Welcome to Apartments website!\n")
}

func (a *Application) apartmentsHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("[apartmentsHandler] received new request", r.Method, r.URL.Path)

	switch r.Method {
	case http.MethodGet:
		w.Header().Set("Content-Type", "application/json")
		allApartments := a.repo.ListAllApartments()
		json.NewEncoder(w).Encode(&allApartments)
	case http.MethodPost:
		var apartment Apartment
		err := json.NewDecoder(r.Body).Decode(&apartment)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		apartment = a.repo.SaveApartment(apartment)
		json.NewEncoder(w).Encode(&apartment)

		message := apartment
		a.publisher.SendMessage("apartment_created", message)

	case http.MethodDelete:
		var body struct {
			Id string `db:"id" json:"id"`
		}
		err := json.NewDecoder(r.Body).Decode(&body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		a.repo.DeleteApartment(body.Id)
		json.NewEncoder(w).Encode(&body)

		message := body
		a.publisher.SendMessage("apartment_deleted", message)

	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

type Application struct {
	repo      *ApartmentRepository
	publisher Publisher
}

func CreateApp() (*Application, error) {
	repo, err := NewApartmentRepository()
	if err != nil {
		log.Println("[apartments_repository] Failed to connect to database", err)
		return nil, err
	}

	const RABBIT_MQ_CONNECTION_STRING = "amqp://guest:guest@rabbitmq:5672/"
	apartmentPublisher, err := NewPublisher(RABBIT_MQ_CONNECTION_STRING)
	if err != nil {
		log.Println("[CreateApp] failed to setup rabbit mq publisher: will retry when first message is sent", err)
		apartmentPublisher = &RetryPublisher{}
	}
	apartmentApplication := Application{repo: repo, publisher: apartmentPublisher}
	return &apartmentApplication, nil
}

func (a *Application) StartApp() {
	http.HandleFunc("/", getRoot)
	http.HandleFunc("/api/apartments", a.apartmentsHandler)

	err := http.ListenAndServe(":3000", nil)
	if errors.Is(err, http.ErrServerClosed) {
		fmt.Printf("server closed\n")
	} else if err != nil {
		fmt.Printf("error starting server: %s\n", err)
		os.Exit(1)
	}
}
