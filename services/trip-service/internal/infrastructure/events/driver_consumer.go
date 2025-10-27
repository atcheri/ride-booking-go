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
	return c.rabbitmq.Consume(messaging.DriverTripResponseQueue, func(ctx context.Context, msg amqp.Delivery) error {
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
			if err := c.handleTripDeclined(ctx, payload.TripID, payload.RiderID); err != nil {
				log.Printf("failed to decline the trip: %v", err)
				return err
			}

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

	// 4. notify the rider that the driver has been assigned
	if err := c.rabbitmq.Publish(ctx, contracts.TripEventDriverAssigned, contracts.AmqpMessage{
		OwnerID: trip.UserID,
		Data:    marshalledTrip,
	}); err != nil {
		return err
	}

	// 5. create the json data to be sent to the payment processor
	marshalledPayload, err := json.Marshal(messaging.PaymentTripResponseData{
		TripID:   tripID,
		UserID:   trip.UserID,
		DriverID: driver.GetId(),
		Amount:   trip.RideFare.TotalPriceInCents,
		Currency: "EUR",
	})
	if err != nil {
		log.Printf("failed to marshall trip payment response data: %v", err)
		return err
	}

	// 6. publish the payment session creation event with the payload
	if err := c.rabbitmq.Publish(ctx, contracts.PaymentCmdCreateSession, contracts.AmqpMessage{
		OwnerID: trip.UserID,
		Data:    marshalledPayload,
	}); err != nil {
		log.Printf("failed to publish the payment session creation event: %v", err)
		return err
	}

	return nil
}

func (c *DriverConsumer) handleTripDeclined(ctx context.Context, tripID, riderID string) error {
	// when a driver declines, we should try to find another done
	trip, err := c.service.GetTripByID(ctx, tripID)
	if err != nil {
		return err
	}

	newPayload := messaging.TripEventData{Trip: trip.ToProto()}
	marshalledPayload, err := json.Marshal(newPayload)
	if err != nil {
		return err
	}

	// ... the easiest way would be by re-publishing the event again
	if err := c.rabbitmq.Publish(ctx, contracts.TripEventDriverNotInterested, contracts.AmqpMessage{
		OwnerID: riderID,
		Data:    marshalledPayload,
	}); err != nil {
		return err
	}

	return nil
}
