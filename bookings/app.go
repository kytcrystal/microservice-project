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

	const APARTMENTS_CREATED_QUEUE = "apartment_created"
	createdMessages, err := createQueue(APARTMENTS_CREATED_QUEUE, channel)
	if err != nil {
		return fmt.Errorf("failed to create apartment_created queue or channel: %w", err)
	}

	go func() {
		for d := range createdMessages {
			var message Apartment
			if err := json.Unmarshal(d.Body, &message); err != nil {
				log.Printf("[error] failed to parse message body: %v", err)
				continue
			}

			log.Printf("Received a message: %+v", message)
			
			SaveApartment(message)
		}
	}()
	log.Printf("Server is listening to queue `%s`", APARTMENTS_CREATED_QUEUE)


	const APARTMENTS_DELETED_QUEUE = "apartment_deleted"
	deletedMessages, err := createQueue(APARTMENTS_DELETED_QUEUE, channel)
	if err != nil {
		return fmt.Errorf("failed to create apartment_deleted queue or channel: %w", err)
	}

	go func() {
		for d := range deletedMessages {
			var message struct{ Id string }
			if err := json.Unmarshal(d.Body, &message); err != nil {
				log.Printf("[error] failed to parse message body: %v", err)
				continue
			}

			log.Printf("Received a message: %+v", message)
			
			DeleteApartment(message.Id)
		}
	}()
	log.Printf("Server is listening to queue `%s`", APARTMENTS_DELETED_QUEUE)

	return nil
}

func createQueue(queueName string,channel *amqp.Channel) (<-chan amqp.Delivery, error) {

	queue, err := channel.QueueDeclare(
		queueName, 	// name
		false,		// durable
		false,		// delete when unused
		false, 		// exclusive
		false, 		// no-wait
		nil, 		// arguments
	)
	if err != nil {
		return nil, err
	}

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
		return nil, err
	}

	return msgs, nil
}
