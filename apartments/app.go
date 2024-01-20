package apartments

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
)

func getRoot(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("got / request\n")
	io.WriteString(w, "Welcome to Apartments website!\n")
}

func (a *Application) apartmentsHandler(w http.ResponseWriter, r *http.Request) {

	switch r.Method {
	case http.MethodGet:
		fmt.Printf("got /api/apartments GET request\n")
		w.Header().Set("Content-Type", "application/json")
		allApartments := ListAllApartments()
		json.NewEncoder(w).Encode(&allApartments)
	case http.MethodPost:
		fmt.Printf("got /api/apartments POST request\n")
		var apartment Apartment
		err := json.NewDecoder(r.Body).Decode(&apartment)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		apartment = SaveApartment(apartment)
		json.NewEncoder(w).Encode(&apartment)

		message := apartment
		a.publisher.SendMessage("apartment_created", message)

	case http.MethodDelete:
		fmt.Printf("got /api/apartments DELETE request\n")
		var body struct{ Id string }
		err := json.NewDecoder(r.Body).Decode(&body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		fmt.Println("here")
		w.Header().Set("Content-Type", "application/json")
		DeleteApartment(body.Id)
		json.NewEncoder(w).Encode(&body)

		message := body
		a.publisher.SendMessage("apartment_deleted", message)

	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

type Application struct {
	publisher *Publisher
}

func CreateApp() (*Application, error) {
	const RABBIT_MQ_CONNECTION_STRING = "amqp://guest:guest@rabbitmq:5672/"
	apartmentPublisher, err := NewPublisher(RABBIT_MQ_CONNECTION_STRING)
	if err != nil {
		return nil, err
	}
	apartmentApplication := Application{publisher: apartmentPublisher}
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
