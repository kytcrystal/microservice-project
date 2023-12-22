package search

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Application struct {
	// in here we will be able to add "objects" will need to be shared across our application
	// for instance shared database connection, or similar things;
	// Having it here rather than it in main make it more easy to test, for instance in a unit test.
}

func NewApplication() Application {
	return Application{}
}

func (a *Application) Run() error {
	log.Println("Starting Search Application")

	publisher, err := NewPublisher(RABBIT_MQ_CONNECTION_STRING, APARTMENTS_QUEUE_NAME)
	if err != nil {
		return fmt.Errorf("failed to create rabbit mq publisher: %w", err)
	}
	defer publisher.Close()

	message := SampleMessage{ID: "123", Number: 123, CreatedAt: time.Now()}
	if err := publisher.SendMessage(message); err != nil {
		return fmt.Errorf("failed to send message: %w", err)
	}

	conn, err := amqp.Dial(RABBIT_MQ_CONNECTION_STRING)
	if err != nil {
		return fmt.Errorf("failed to create rabbit mq connection: %w", err)
	}

	channel, err := conn.Channel()
	if err != nil {
		return fmt.Errorf("failed to create rabbit mq channel: %w", err)
	}

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

	// this is a thing that allow us to wait forever
	var forever chan struct{}

	go func() {
		for d := range msgs {
			var message SampleMessage
			if err := json.Unmarshal(d.Body, &message); err != nil {
				log.Printf("[error] failed to parse message body: %v", err)
			}

			log.Printf("Received a message: %+v", message)
		}
	}()

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever

	return nil
}

func toJsonBytes(v interface{}) []byte {
	value, _ := json.Marshal(v)
	return value
}

type SampleMessage struct {
	ID        string
	Number    int
	CreatedAt time.Time
}
