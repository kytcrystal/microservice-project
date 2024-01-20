package bookings

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	amqp "github.com/rabbitmq/amqp091-go"
)

var listOfBookings []Booking

func Run() error {
	var port = "3000"

	http.HandleFunc("/api/bookings", bookingsHandler)

	StartListener()

	log.Println("[app:run] starting booking service at port", port)
	err := http.ListenAndServe(":"+port, nil)
	return err
}

func bookingsHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("[bookingsHandler] received request: ", r.Method, r.URL.Path)

	w.Header().Set("Content-Type", "application/json")

	switch r.Method {

	case http.MethodGet:
		allBookings := ListAllBookings()
		json.NewEncoder(w).Encode(&allBookings)

	case http.MethodPost:
		var booking Booking
		err := json.NewDecoder(r.Body).Decode(&booking)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		newBooking, err := CreateBooking(booking)
		if err != nil {
			fmt.Fprintf(w, "error in booking %v", err)
			return
		}
		json.NewEncoder(w).Encode(&newBooking)
	case http.MethodDelete:
		var body struct{ ID string }
		err := json.NewDecoder(r.Body).Decode(&body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		CancelBooking(body.ID)
		json.NewEncoder(w).Encode(&body)

	case http.MethodPatch:
		var booking Booking
		err := json.NewDecoder(r.Body).Decode(&booking)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		ChangeBooking(booking)
		json.NewEncoder(w).Encode(&booking)

	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func StartListener() error {
	log.Println("Starting Listener For Apartment Messages")
	const RABBIT_MQ_CONNECTION_STRING = "amqp://guest:guest@rabbitmq:5672/"

	conn, err := amqp.Dial(RABBIT_MQ_CONNECTION_STRING)
	if err != nil {
		return fmt.Errorf("failed to create rabbit mq connection: %w", err)
	}

	channel, err := conn.Channel()
	if err != nil {
		return fmt.Errorf("failed to create rabbit mq channel: %w", err)
	}

	const APARTMENTS_QUEUE_NAME = "apartment_created"
	queue, err := channel.QueueDeclare(
		APARTMENTS_QUEUE_NAME, // name
		false,                 // durable
		false,                 // delete when unused
		false,                 // exclusive
		false,                 // no-wait
		nil,                   // arguments
	)
	if err != nil {
		return fmt.Errorf("failed to create new queue: %v", err)
	}

	// Just as a simple example let's try to set up the listener as well

	msgs, err := channel.Consume(
		queue.Name, // queue
		"",         // consumer
		true,       // auto-ack
		false,      // exclusive
		false,      // no-local
		false,      // no-wait
		nil,        // args
	)

	if err != nil {
		return fmt.Errorf("failed to create message consumer: %w", err)
	}

	go func() {
		for d := range msgs {
			var message Apartment
			if err := json.Unmarshal(d.Body, &message); err != nil {
				log.Printf("[error] failed to parse message body: %v", err)
				continue
			}

			log.Printf("Received a message: %+v", message)
			
			SaveApartment(message)
		}
	}()

	log.Printf("Server is listening to queue `%s`", queue.Name)
	return nil
}
