package apartments

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

func (a *Application) CustomHandleFunc(pattern string, handle func(*http.Request) (any, error)) {
	a.HandleFunc(pattern, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		value, err := handle(r)
		if err != nil {
			log.Println("Encountered error while handling request", r.Method, r.URL.Path, err)
			value = struct{ Error string }{err.Error()}
			w.WriteHeader(http.StatusInternalServerError)
		}
		json.NewEncoder(w).Encode(value)
	})
}

func (a *Application) apartmentsHandler(r *http.Request) (any, error) {
	log.Println("[apartmentsHandler] received new request", r.Method, r.URL.Path)

	switch r.Method {
	case http.MethodGet:
		return a.repo.ListAllApartments()
	case http.MethodPost:
		return a.createApartment(r)
	case http.MethodDelete:
		return a.deleteApartment(r)
	default:
		return nil, fmt.Errorf("method not allowed: %s", r.Method)
	}
}

func (a *Application) createApartment(r *http.Request) (any, error) {
	var apartment Apartment
	err := json.NewDecoder(r.Body).Decode(&apartment)
	if err != nil {
		return nil, err
	}

	newApartment := a.repo.SaveApartment(apartment)
	a.publisher.SendMessage(MQ_APARTMENT_CREATED_EXCHANGE, newApartment)
	return newApartment, nil
}

func (a *Application) deleteApartment(r *http.Request) (any, error) {
	var body struct {
		Id string `db:"id" json:"id"`
	}
	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		return nil, err
	}
	a.repo.DeleteApartment(body.Id)

	a.publisher.SendMessage(MQ_APARTMENT_DELETED_EXCHANGE, body)
	return body, nil
}
