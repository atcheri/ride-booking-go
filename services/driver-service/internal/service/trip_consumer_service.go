package service

import (
	"context"
	"log"

	amqp "github.com/rabbitmq/amqp091-go"

	"github.com/atcheri/ride-booking-go/shared/messaging"
)

type TripConsumerService struct {
	rabbitmq *messaging.RabbitMQ
}

func NewTripConsumerService(rabbitmq *messaging.RabbitMQ) *TripConsumerService {
	return &TripConsumerService{
		rabbitmq: rabbitmq,
	}
}

func (c *TripConsumerService) Listen() error {
	queueName := "hello"
	return c.rabbitmq.Consume(queueName, func(ctx context.Context, msg amqp.Delivery) error {
		log.Printf("driver received a message: %v", msg)
		return nil
	})
}
