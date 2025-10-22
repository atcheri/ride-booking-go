package main

import (
	"encoding/json"
	"log"
	"net/http"

	grpcclient "github.com/atcheri/ride-booking-go/services/api-gateway/grpc_client"
	"github.com/atcheri/ride-booking-go/shared/contracts"
	"github.com/atcheri/ride-booking-go/shared/env"
)

var (
	tripServiceURL = env.GetString("TRIP_SERVICE_URL", "trip-service:9093")
)

func handleTripPreview(w http.ResponseWriter, r *http.Request) {
	var body tripPreviewRequest

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, "failed to parse JSON data", http.StatusBadRequest)
		return
	}

	defer r.Body.Close()

	// validation
	if body.UserId == "" {
		http.Error(w, "user ID is required", http.StatusBadRequest)
		return
	}

	tripServiceClient, err := grpcclient.NewTripServiceClient(tripServiceURL)
	if err != nil {
		log.Fatal(err)
	}

	// Closing the connection on each request
	defer tripServiceClient.Close()

	tripPreviewResp, err := tripServiceClient.Client.PreviewTrip(r.Context(), body.ToProto())
	if err != nil {
		log.Printf("failed to preview the trip: %v", err)
		http.Error(w, "failed to preview the trip", http.StatusInternalServerError)
		return
	}

	response := contracts.APIResponse{Data: tripPreviewResp}

	writeJSON(w, http.StatusCreated, response)
}
