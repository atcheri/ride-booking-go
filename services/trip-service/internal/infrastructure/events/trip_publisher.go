package events

import (
	"context"

	"github.com/atcheri/ride-booking-go/shared/messaging"
)

type TripEventPubliser struct {
	rabbitmq *messaging.RabbitMQ
}

func NewTripEventPublisher(rabbitmq *messaging.RabbitMQ) *TripEventPubliser {
	return &TripEventPubliser{
		rabbitmq: rabbitmq,
	}
}

func (p *TripEventPubliser) PublishTripCreated(ctx context.Context) error {
	exchangeName := ""
	queueName := "hello"
	body := "Hello World!"

	return p.rabbitmq.Publish(ctx, exchangeName, queueName, body)
}
