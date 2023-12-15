package search

import (
	"context"
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

	// TODO: make this configurable via env variable so it's easy to Dockerize the app
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		return fmt.Errorf("failed to create rabbit mq connection: %w", err)
	}
	//TODO: when stop the app we should close the connection

	channel, err := conn.Channel()
	if err != nil {
		return fmt.Errorf("failed to create rabbit mq channel: %w", err)
	}
	//TODO: we should also close the channel later on

	queue, err := channel.QueueDeclare(
		"apparments-queue", // name
		false,              // durable
		false,              // delete when unused
		false,              // exclusive
		false,              // no-wait
		nil,                // arguments
	)
	if err != nil {
		return fmt.Errorf("failed to create new queue: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = channel.PublishWithContext(ctx,
		"",         // exchange
		queue.Name, // routing key
		false,      // mandatory
		false,      // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        toJsonBytes(SampleMessage{ID: "123", Number: 123, CreatedAt: time.Now()}),
		})

	if err != nil {
		return fmt.Errorf("failed to send message: %w", err)
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
