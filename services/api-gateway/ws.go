package main

import (
	"log"
	"net/http"

	grpcclient "github.com/atcheri/ride-booking-go/services/api-gateway/grpc_client"
	"github.com/atcheri/ride-booking-go/shared/contracts"
	"github.com/atcheri/ride-booking-go/shared/env"
	"github.com/gorilla/websocket"

	pb "github.com/atcheri/ride-booking-grpc-proto/golang/driver"
)

var (
	upgrader = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
	driverServiceURL = env.GetString("DRIVER_SERVICE_URL", "driver-service:9092")
)

func handleDriversWebSocket(w http.ResponseWriter, r *http.Request) {
	log.Println("reaching at least here")
	conn, err := upgrader.Upgrade(w, r, nil)

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

	// create a new grpc client
	driverService, err := grpcclient.NewDriverServiceClient(driverServiceURL)
	if err != nil {
		log.Fatal(err)
	}

	// closing the grpc connection after unregistering the driver when the ws connection is closed
	defer func() {
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

	msg := contracts.WSMessage{
		Type: "driver.cmd.register",
		Data: driverData.Driver,
	}

	if err := conn.WriteJSON(msg); err != nil {
		log.Printf("error sending message to the driver: %v", err)
		return
	}

	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			log.Printf("error reading message from drivers websocket: %V", err)
			break
		}

		log.Printf("received message from drivers websocket: %s", message)
	}

}

func handleRidersWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)

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

	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			log.Printf("error reading message from riders websocket: %V", err)
			break
		}

		log.Printf("received message from riders websocket: %s", message)
	}
}
