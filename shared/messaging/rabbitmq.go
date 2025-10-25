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

type MessageHandler func(context.Context, amqp.Delivery) error

func (r *RabbitMQ) Consume(queueName string, handler MessageHandler) error {
	err := r.Channel.Qos(
		1,     // prefetchCount: Limit to 1 unacknowledged message per consumer; set to 1 for a fair dispatch
		0,     // prefetchSize: no specific limit on message size
		false, // global: apply prefetchCount to each consumer individually
	)
	if err != nil {
		return fmt.Errorf("failed to set QoS: %v", err)
	}

	msgs, err := r.Channel.Consume(
		queueName, // queue
		"",        // consumer
		false,     // auto-ack
		false,     // exclusive
		false,     // no-local
		false,     // no-wait
		nil,       // args
	)
	if err != nil {
		return err
	}

	ctx := context.Background()

	go func() {
		for msg := range msgs {
			log.Printf("received a message: %s", msg.Body)

			if err := handler(ctx, msg); err != nil {
				log.Printf("ERROR: Failed to handle message: %v. Message body: %s", err, msg.Body)

				if nackErr := msg.Nack(false, false); nackErr != nil {
					log.Printf("ERROR: failed to Nack message: %V", nackErr)
				}

				continue
			}

			// Acknowledging explicitely for each consumed message
			if ackErr := msg.Ack(false); ackErr != nil {
				log.Printf("ERROR: Failed to Ack message: %v. Message body: %s", ackErr, msg.Body)
			}
		}
	}()

	return nil

}
