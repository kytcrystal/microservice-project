package bookings

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Publisher interface {
	SendMessage(queueName string, message interface{}) error
}

type RetryPublisher struct {
	Publisher
}

func (p *RetryPublisher) SendMessage(exchangeName string, message interface{}) error {
	if p.Publisher == nil {
		log.Println("[RetryPublisher] rabbit mq publisher was not set up correcty, will attempt to create it now")
		publisher, err := NewPublisher(MQ_CONNECTION_STRING)
		if err != nil {
			log.Println("[RetryPublisher] failed to setup publisher, message will be skipped", exchangeName, message)
			return err
		}
		p.Publisher = publisher
	}
	return p.Publisher.SendMessage(exchangeName, message)
}

type SimplePublisher struct {
	connection *amqp.Connection
	channel    *amqp.Channel
}

func NewPublisher(dsn string, exchanges ...string) (Publisher, error) {
	conn, err := amqp.Dial(dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to create rabbit mq connection: %w", err)
	}

	channel, err := conn.Channel()
	if err != nil {
		return nil, fmt.Errorf("failed to create rabbit mq channel: %w", err)
	}

	publisher := &SimplePublisher{
		connection: conn,
		channel:    channel,
	}

	for _, exchange := range exchanges {
		err := publisher.declareFanoutExchange(exchange)
		if err != nil {
			log.Println("failed to create exchange: will retry when message is sent", exchange, err)
		}
	}

	return publisher, nil
}

func (p *SimplePublisher) SendMessage(exchangeName string, message interface{}) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := p.declareFanoutExchange(exchangeName)
	if err != nil {
		return fmt.Errorf("failed to declare exchange:  %w", err)
	}

	err = p.channel.PublishWithContext(ctx,
		exchangeName, // exchange
		"",           // routing key: empty cause we publish to the exchange
		false,        // mandatory
		false,        // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        toJsonBytes(message),
		})

	if err != nil {
		return fmt.Errorf("failed to send message: %w", err)
	}

	log.Println("succesfully sent message to to exchange", exchangeName, message)
	return nil
}

func (p *SimplePublisher) declareFanoutExchange(exchangeName string) error {
	err := p.channel.ExchangeDeclare(
		exchangeName, // name
		"fanout",     // type
		true,         // durable
		false,        // auto-deleted
		false,        // internal
		false,        // no-wait
		nil,          // arguments
	)
	return err
}

func (p *SimplePublisher) Close() error {
	p.channel.Close()
	p.connection.Close()
	return nil
}

func toJsonBytes(v interface{}) []byte {
	value, _ := json.Marshal(v)
	return value
}
