package messaging

import (
	"context"
	"fmt"
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

type RabbitMQ struct {
	conn    *amqp.Connection
	Channel *amqp.Channel
}

func NewRabbitMQ(uri string) (*RabbitMQ, error) {
	conn, err := amqp.Dial(uri)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to RabbitMQ: %v", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		conn.Close()
		return nil, fmt.Errorf("failed to create the RabbitMQ channel: %v", err)
	}

	rabbitMQ := &RabbitMQ{
		conn:    conn,
		Channel: ch,
	}

	if err := rabbitMQ.setupExchangesAndQueues(); err != nil {
		conn.Close()
		return nil, fmt.Errorf("failed to setup rabbitMQ exchanges and queues: %v", err)
	}

	return rabbitMQ, nil
}

func (r *RabbitMQ) Close() {
	if r.conn != nil {
		r.conn.Close()
	}
	if r.Channel != nil {
		r.Channel.Close()
	}
}

func (r *RabbitMQ) setupExchangesAndQueues() error {
	_, err := r.Channel.QueueDeclare(
		"hello", // name
		true,    // durable
		false,   // delete when unused
		false,   // exclusive
		false,   // no-wait
		nil,     // arguments
	)

	if err != nil {
		log.Fatal(err)
		return err
	}

	return nil
}

func (r *RabbitMQ) Publish(ctx context.Context, exchangeName, queueName, message string) error {
	return r.Channel.PublishWithContext(ctx,
		exchangeName,
		queueName,
		false,
		false,
		amqp.Publishing{
			ContentType:  "text/plain",
			Body:         []byte(message),
			DeliveryMode: amqp.Persistent,
		})
}
