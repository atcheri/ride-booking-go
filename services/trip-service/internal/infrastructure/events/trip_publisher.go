package events

import (
	"context"
	"encoding/json"
	"log"

	"github.com/atcheri/ride-booking-go/services/trip-service/internal/domain/models"
	"github.com/atcheri/ride-booking-go/shared/contracts"
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

func (p *TripEventPubliser) PublishTripCreated(ctx context.Context, trip *models.TripModel) error {
	routingKey := contracts.TripEventCreated
	log.Printf("publish created trip: %+v", trip)
	payload := messaging.TripEventData{
		Trip: trip.ToProto(),
	}
	tripEventJSON, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	return p.rabbitmq.Publish(ctx, routingKey, contracts.AmqpMessage{
		OwnerID: trip.UserID,
		Data:    tripEventJSON,
	})
}
