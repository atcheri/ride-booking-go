package main

import (
	"log"
	"math/rand"
	"net/http"

	"github.com/atcheri/ride-booking-go/shared/contracts"
	"github.com/atcheri/ride-booking-go/shared/util"
	"github.com/gorilla/websocket"
)

var (
	upgrader = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
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

	type Driver struct {
		Id             string `json:"id"`
		Name           string `json:"name"`
		ProfilePicture string `json:"profilePicture"`
		CarPlate       string `json:"carPlate"`
		PackageSlug    string `json:"packageSlug"`
	}

	msg := contracts.WSMessage{
		Type: "driver.cmd.register",
		Data: Driver{
			Id:             userID,
			Name:           "Atch",
			ProfilePicture: util.GetRandomAvatar(rand.Intn(9) + 1),
			CarPlate:       "ENDO16",
			PackageSlug:    packageSlug,
		},
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
