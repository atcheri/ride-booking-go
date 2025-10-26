package service

import (
	"context"
	"encoding/json"
	"log"

	amqp "github.com/rabbitmq/amqp091-go"

	"github.com/atcheri/ride-booking-go/shared/contracts"
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
	return c.rabbitmq.Consume(messaging.FindAvailableDriversQueue, func(ctx context.Context, msg amqp.Delivery) error {
		var tripEvent contracts.AmqpMessage
		if err := json.Unmarshal(msg.Body, &tripEvent); err != nil {
			log.Printf("failed to unmarshal the message: %v", err)
			return err
		}

		var payload messaging.TripEventData
		if err := json.Unmarshal(tripEvent.Data, &payload); err != nil {
			log.Printf("failed to unmarshal the message %v", err)
			return err
		}

		log.Printf("driver received a message: %v", payload)
		return nil
	})
}
