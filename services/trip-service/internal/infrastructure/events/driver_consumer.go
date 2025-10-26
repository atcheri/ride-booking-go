package events

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	amqp "github.com/rabbitmq/amqp091-go"

	"github.com/atcheri/ride-booking-go/services/trip-service/internal/domain/service"
	"github.com/atcheri/ride-booking-go/shared/contracts"
	"github.com/atcheri/ride-booking-go/shared/messaging"
	"github.com/atcheri/ride-booking-grpc-proto/golang/driver"
)

type DriverConsumer struct {
	rabbitmq *messaging.RabbitMQ
	service  service.TripService
}

func NewDriverConsumer(rabbitmq *messaging.RabbitMQ, service service.TripService) *DriverConsumer {
	return &DriverConsumer{
		rabbitmq: rabbitmq,
		service:  service,
	}
}

func (c *DriverConsumer) Listen() error {
	return c.rabbitmq.Consume(messaging.FindAvailableDriversQueue, func(ctx context.Context, msg amqp.Delivery) error {
		var message contracts.AmqpMessage
		if err := json.Unmarshal(msg.Body, &message); err != nil {
			log.Printf("failed to unmarshal the message: %v", err)
			return err
		}

		var payload messaging.DriverTripResponseData
		if err := json.Unmarshal(message.Data, &payload); err != nil {
			log.Printf("failed to unmarshal the message %v", err)
			return err
		}

		log.Printf("driver response received message: %+v", payload)

		switch msg.RoutingKey {
		case contracts.DriverCmdTripAccept:
			if err := c.handleTripAccepted(ctx, payload.TripID, payload.Driver); err != nil {
				log.Printf("failed to handle the accepted trip: %v", err)
			}
		case contracts.DriverCmdTripDecline:
			// TODO: add a trip decline handler function
			log.Println("trip declined")
			return nil
		}

		log.Printf("unknown trip event: %+v", payload)

		return nil
	})
}

func (c *DriverConsumer) handleTripAccepted(ctx context.Context, tripID string, driver *driver.Driver) error {
	// 1. fetch the trip
	trip, err := c.service.GetTripByID(ctx, tripID)
	if err != nil {
		return err
	}

	if trip == nil {
		return fmt.Errorf("trip with id %s was not found", tripID)
	}

	// 2. update the trip
	if err := c.service.UpdateTrip(ctx, tripID, "accepted", driver); err != nil {
		log.Printf("failed to update the trip: %v", err)
		return err
	}

	// refetch the updated trip again, similar to 1.
	trip, err = c.service.GetTripByID(ctx, tripID)
	if err != nil {
		return err
	}

	// 3. driver is assigned -> publish the event to RB
	marshalledTrip, err := json.Marshal(trip)
	if err != nil {
		return err
	}

	if err := c.rabbitmq.Publish(ctx, contracts.TripEventDriverAssigned, contracts.AmqpMessage{
		OwnerID: trip.UserID,
		Data:    marshalledTrip,
	}); err != nil {
		return err
	}

	// TODO: notify the payment service to start a payment link

	return nil
}
