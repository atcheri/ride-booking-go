package main

import (
	"encoding/json"
	"log"
	"net/http"

	grpcclient "github.com/atcheri/ride-booking-go/services/api-gateway/grpc_client"
	"github.com/atcheri/ride-booking-go/shared/contracts"
	"github.com/atcheri/ride-booking-go/shared/env"
	"github.com/atcheri/ride-booking-go/shared/messaging"

	pb "github.com/atcheri/ride-booking-grpc-proto/golang/driver"
)

var (
	driverServiceURL = env.GetString("DRIVER_SERVICE_URL", "driver-service:9092")
	connManager      = messaging.NewConnectionManager()
)

func handleDriversWebSocketWithRabbitMQ(rb *messaging.RabbitMQ) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		conn, err := connManager.Upgrade(w, r)
		if err != nil {
			log.Printf("drivers websocket upgrade failed: %v", err)
			return
		}

		defer conn.Close()

		userID := r.URL.Query().Get("userID")
		if userID == "" {
			log.Println("no user ID was provided")
			return
		}

		packageSlug := r.URL.Query().Get("packageSlug")
		if packageSlug == "" {
			log.Println("no package-slug was provided")
			return
		}

		// Add connection to the manager
		connManager.Add(userID, conn)

		// create a new grpc client
		driverService, err := grpcclient.NewDriverServiceClient(driverServiceURL)
		if err != nil {
			log.Fatal(err)
		}

		// closing the grpc connection after unregistering the driver when the ws connection is closed
		defer func() {
			connManager.Remove(userID)
			driverService.Client.UnregisterDriver(r.Context(), &pb.RegisterDriverRequest{
				DriverID:    userID,
				PackageSlug: packageSlug,
			})
			driverService.Close()
			log.Println("driver unregistered: ", userID)
		}()

		driverData, err := driverService.Client.RegisterDriver(r.Context(), &pb.RegisterDriverRequest{
			DriverID:    userID,
			PackageSlug: packageSlug,
		})
		if err != nil {
			log.Printf("error registering the driver: %v", err)
		}

		if err := connManager.SendMessage(userID, contracts.WSMessage{
			Type: contracts.DriverCmdRegister,
			Data: driverData.Driver,
		}); err != nil {
			log.Printf("error sending message to the driver: %v", err)
			return
		}

		// initialize the queue consumers
		queues := []string{
			messaging.DriverCmdTripRequestQueue,
		}

		for _, q := range queues {
			consumer := messaging.NewQueueConsumer(rb, connManager, q)

			if err := consumer.Start(); err != nil {
				log.Printf("failed to start consumer for the queue: %s; err: %v", q, err)
			}
		}

		for {
			_, message, err := conn.ReadMessage()
			if err != nil {
				log.Printf("error reading message from drivers websocket: %V", err)
				break
			}

			type driverMessage struct {
				Type string          `json:"type"`
				Data json.RawMessage `json:"data"`
			}
			var driverMsg driverMessage
			if err := json.Unmarshal(message, &driverMsg); err != nil {
				log.Printf("error unmarshalling driver message: %v", err)
				continue
			}

			// handle the different message by type
			switch driverMsg.Type {
			case contracts.DriverCmdLocation:
				// handle driver location update
				continue
			case contracts.DriverCmdTripAccept, contracts.DriverCmdTripDecline:
				// forward the message to RabbitMQ
				if err := rb.Publish(r.Context(), driverMsg.Type, contracts.AmqpMessage{
					OwnerID: userID,
					Data:    driverMsg.Data,
				}); err != nil {
					log.Printf("error publishing the message to RabbitMQ: %v", err)
				}
			default:
				log.Printf("unknown message type: %s", driverMsg.Type)
			}

		}
	}
}

func handleRidersWebSocketWithRabbitMQ(rb *messaging.RabbitMQ) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		conn, err := connManager.Upgrade(w, r)

		log.Println("ws connection established")

		if err != nil {
			log.Printf("riders websocket upgrade failed: %v", err)
			return
		}

		defer conn.Close()

		userID := r.URL.Query().Get("userID")
		if userID == "" {
			log.Println("no user ID was provided")
			return
		}

		// Add connection to the manager
		connManager.Add(userID, conn)
		defer connManager.Remove(userID)

		for {
			_, message, err := conn.ReadMessage()
			if err != nil {
				log.Printf("error reading message from riders websocket: %V", err)
				break
			}

			log.Printf("received message from riders websocket: %s", message)
		}
	}
}
