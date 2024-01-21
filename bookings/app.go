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
	apartmentPublisher, err := NewPublisher(MQ_CONNECTION_STRING)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize rabbit mq: %w", err)
	}

	return &Application{
		publisher: *apartmentPublisher,
		ServeMux:  http.DefaultServeMux,
	}, nil
}

func (a *Application) Run() error {
	a.CustomHandleFunc("/api/bookings", a.bookingsHandler)

	StartListener()

	log.Println("[app:run] starting booking service at port", PORT)
	err := http.ListenAndServe(":"+PORT, a.ServeMux)
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
		bookingCreated, err := CreateBooking(booking)
		if err != nil {
			return nil, err
		}
		err = a.publisher.SendMessage("booking_created", BookingCreatedEvent{bookingCreated})
		if err != nil {
			return nil, err
		}
		return bookingCreated, nil
	case http.MethodDelete:
		var body struct{ ID string }
		err := json.NewDecoder(r.Body).Decode(&body)
		if err != nil {
			return nil, err
		}
		if err := CancelBooking(body.ID); err != nil {
			return nil, err
		}
		err = a.publisher.SendMessage("booking_cancelled", BookingCancelledEvent{body.ID})
		if err != nil {
			return nil, err
		}
		return body, nil
	case http.MethodPatch:
		var booking Booking
		err := json.NewDecoder(r.Body).Decode(&booking)
		if err != nil {
			return nil, err
		}
		updatedBooking, err := ChangeBooking(booking)
		if err != nil {
			return nil, err
		}
		err = a.publisher.SendMessage("booking_updated", BookingUpdatedEvent{updatedBooking})
		if err != nil {
			return nil, err
		}
		return updatedBooking, nil
	default:
		return nil, fmt.Errorf("method not allowed")
	}
}

func StartListener() error {
	log.Println("Starting Listener For Apartment Messages")

	conn, err := amqp.Dial(MQ_CONNECTION_STRING)
	if err != nil {
		return fmt.Errorf("failed to create rabbit mq connection: %w", err)
	}

	channel, err := conn.Channel()
	if err != nil {
		return fmt.Errorf("failed to create rabbit mq channel: %w", err)
	}

	createdMessages, err := createQueue(MQ_APARTMENT_CREATED_EXCHANGE, MQ_APARTMENT_CREATED_QUEUE, channel)
	if err != nil {
		return fmt.Errorf("failed to create apartment_created queue or channel: %w", err)
	}

	go func() {
		log.Println("initialized goroutine that process created messages: ", createdMessages)
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

	deletedMessages, err := createQueue(MQ_APARTMENT_DELETED_EXCHANGE, MQ_APARTMENT_DELETED_QUEUE, channel)
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

	return nil
}

func createQueue(
	exchangeName string,
	queueName string,
	channel *amqp.Channel,
) (<-chan amqp.Delivery, error) {

	queue, err := channel.QueueDeclare(
		queueName, // name
		true,      // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // no-wait
		nil,       // arguments
	)
	if err != nil {
		return nil, err
	}

	err = channel.QueueBind(
		queueName,    // queue name
		"",           // routing key
		exchangeName, // exchange
		false,
		nil,
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

	log.Printf("Server is listening to queue `%s` in exchange `%s`", queueName, exchangeName)
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
