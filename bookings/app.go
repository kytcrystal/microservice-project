package bookings

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Application struct {
	publisher Publisher
	*http.ServeMux
}

func NewApplication() (*Application, error) {
	const RABBIT_MQ_CONNECTION_STRING = "amqp://guest:guest@rabbitmq:5672/"
	apartmentPublisher, err := NewPublisher(RABBIT_MQ_CONNECTION_STRING)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize rabbit mq: %w", err)
	}

	return &Application{
		publisher: *apartmentPublisher,
		ServeMux:  http.DefaultServeMux,
	}, nil
}

func (a *Application) Run() error {
	var port = "3000"

	a.CustomHandleFunc("/api/bookings", a.bookingsHandler)

	StartListener()

	log.Println("[app:run] starting booking service at port", port)
	err := http.ListenAndServe(":"+port, a.ServeMux)
	return err
}

func (a *Application) bookingsHandler(w http.ResponseWriter, r *http.Request) (any, error) {
	log.Println("[bookingsHandler] received request: ", r.Method, r.URL.Path)
	switch r.Method {

	case http.MethodGet:
		return ListAllBookings()

	case http.MethodPost:
		var booking Booking
		err := json.NewDecoder(r.Body).Decode(&booking)
		if err != nil {
			return nil, err
		}
		return CreateBooking(booking)
	case http.MethodDelete:
		var body struct{ ID string }
		err := json.NewDecoder(r.Body).Decode(&body)
		if err != nil {
			return nil, err
		}
		if err := CancelBooking(body.ID); err != nil {
			return nil, err
		}
		return body, nil
	case http.MethodPatch:
		var booking Booking
		err := json.NewDecoder(r.Body).Decode(&booking)
		if err != nil {
			return nil, err
		}
		return ChangeBooking(booking)
	default:
		return nil, fmt.Errorf("method not allowed")
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

func createQueue(queueName string, channel *amqp.Channel) (<-chan amqp.Delivery, error) {

	queue, err := channel.QueueDeclare(
		queueName, // name
		false,     // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // no-wait
		nil,       // arguments
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

func (a *Application) CustomHandleFunc(pattern string, handle func(http.ResponseWriter, *http.Request) (any, error)) {
	a.HandleFunc(pattern, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		value, err := handle(w, r)
		if err != nil {
			log.Println("Encountered error while handling request", r.Method, r.URL.Path, err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(value)
	})
}
