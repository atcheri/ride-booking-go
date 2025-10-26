package service

import (
	"context"
	"encoding/json"
	"log"
	"math/rand/v2"

	amqp "github.com/rabbitmq/amqp091-go"

	"github.com/atcheri/ride-booking-go/services/driver-service/domain/service"
	"github.com/atcheri/ride-booking-go/shared/contracts"
	"github.com/atcheri/ride-booking-go/shared/messaging"
)

type TripConsumerService struct {
	rabbitmq *messaging.RabbitMQ
	service  service.DriverService
}

func NewTripConsumerService(rabbitmq *messaging.RabbitMQ, service service.DriverService) *TripConsumerService {
	return &TripConsumerService{
		rabbitmq: rabbitmq,
		service:  service,
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

		switch msg.RoutingKey {
		case contracts.TripEventCreated, contracts.TripEventDriverNotInterested:
			return c.handleFindAndNotifyDrivers(ctx, payload)
		}

		log.Printf("unknown trip event: %+v", payload)

		return nil
	})
}

func (c *TripConsumerService) handleFindAndNotifyDrivers(ctx context.Context, payload messaging.TripEventData) error {
	suitableDriverIDs := c.service.FindAvailableDrivers(payload.Trip.SelectedFare.PackageSlug)

	log.Printf("found %d suitable drivers", len(suitableDriverIDs))

	if len(suitableDriverIDs) == 0 {
		// notify the rider that no drivers are available
		if err := c.rabbitmq.Publish(ctx, contracts.TripEventNoDriversFound, contracts.AmqpMessage{
			OwnerID: payload.Trip.UserID,
		}); err != nil {
			log.Printf("failed to publish the message to the exchange: %v", err)
			return err
		}

		return nil
	}

	// get a random index from the matching drivers
	randomIndex := rand.IntN(len(suitableDriverIDs))
	firstSuitableDriver := suitableDriverIDs[randomIndex]
	marchalledEvent, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	// notify the driver about a trip request
	if err := c.rabbitmq.Publish(ctx, contracts.DriverCmdTripRequest, contracts.AmqpMessage{
		OwnerID: firstSuitableDriver,
		Data:    marchalledEvent,
	}); err != nil {
		log.Printf("failed to publish the message to the exchange: %v", err)
		return err
	}

	return nil
}
