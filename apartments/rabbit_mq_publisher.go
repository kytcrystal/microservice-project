package apartments

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Publisher struct {
	connection *amqp.Connection
	channel    *amqp.Channel
}

func NewPublisher(dsn string) (*Publisher, error) {

	conn, err := amqp.Dial(dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to create rabbit mq connection: %w", err)
	}

	channel, err := conn.Channel()
	if err != nil {
		return nil, fmt.Errorf("failed to create rabbit mq channel: %w", err)
	}

	return &Publisher{
		connection: conn,
		channel:    channel,
	}, nil

}

func (p *Publisher) SendMessage(queueName string, message interface{}) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := p.channel.QueueDeclare(
		queueName, // name
		false,     // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // no-wait
		nil,       // arguments
	)
	if err != nil {
		return  fmt.Errorf("queue does not exist or is misconfigured: %v", err)
	}

	err = p.channel.PublishWithContext(ctx,
		"",           // exchange
		queueName, // routing key
		false,        // mandatory
		false,        // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        toJsonBytes(message),
		})

	if err != nil {
		return fmt.Errorf("failed to send message: %w", err)
	}
	return nil
}

func (p *Publisher) Close() error {
	p.channel.Close()
	p.connection.Close()
	return nil
}

func toJsonBytes(v interface{}) []byte {
	value, _ := json.Marshal(v)
	return value
}
