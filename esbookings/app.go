package esbookings

import (
	"context"
	"encoding/json"
	"esbookings/eventsourcing"
	"fmt"
	"log"
	"net/http"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Application struct {
	publisher Publisher
	service   *eventsourcing.Service
	*http.ServeMux
}

func NewApplication() (*Application, error) {
	apartmentPublisher, err := NewPublisher(
		MQ_CONNECTION_STRING,
		MQ_BOOKING_CREATED_EXCHANGE,
		MQ_BOOKING_UPDATED_EXCHANGE,
		MQ_BOOKING_CANCELLED_EXCHANGE,
	)
	if err != nil {
		log.Println("[NewApplication] failed to setup rabbit mq publisher: will retry when first message is sent", err)
		apartmentPublisher = &RetryPublisher{}
	}

	service, err := eventsourcing.NewService()
	if err != nil {
		return nil, err
	}

	return &Application{
		publisher: apartmentPublisher,
		service:   service,
		ServeMux:  http.DefaultServeMux,
	}, nil
}

func (a *Application) Run() error {
	a.CustomHandleFunc("/api/bookings", a.bookingsHandler)
	a.CustomHandleFunc("/api/rollback", a.rollbackHandler)

	if err := a.StartListener(); err != nil {
		log.Println("[app:run] Failed to startup booking listeners", err)
		return err
	}

	log.Println("[app:run] starting booking service at port", PORT)
	return http.ListenAndServe(":"+PORT, a.ServeMux)
}

func (a *Application) StartListener() error {
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

			a.service.CreateApartment(context.Background(), message.Id, message.Apartment_Name)
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
			a.service.DeleteApartment(context.Background(), message.Id)
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
