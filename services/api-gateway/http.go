package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"

	"github.com/stripe/stripe-go/v81"
	"github.com/stripe/stripe-go/v81/webhook"

	grpcclient "github.com/atcheri/ride-booking-go/services/api-gateway/grpc_client"
	"github.com/atcheri/ride-booking-go/shared/contracts"
	"github.com/atcheri/ride-booking-go/shared/env"
	"github.com/atcheri/ride-booking-go/shared/messaging"
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

func handleStartTrip(w http.ResponseWriter, r *http.Request) {
	var body startTripRequest

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, "failed to parse JSON data", http.StatusBadRequest)
		return
	}

	defer r.Body.Close()

	// validation
	if body.UserID == "" {
		http.Error(w, "user ID is required", http.StatusBadRequest)
		return
	}

	tripServiceClient, err := grpcclient.NewTripServiceClient(tripServiceURL)
	if err != nil {
		log.Fatal(err)
	}

	// Closing the connection on each request
	defer tripServiceClient.Close()

	ceateTripResp, err := tripServiceClient.Client.CreateTrip(r.Context(), body.ToProto())
	if err != nil {
		log.Printf("failed to create the trip: %v", err)
		http.Error(w, "failed to create the trip", http.StatusInternalServerError)
		return
	}

	response := contracts.APIResponse{Data: ceateTripResp}

	writeJSON(w, http.StatusCreated, response)
}

func handleStripWebhookWithRabbitMQ(rb *messaging.RabbitMQ) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "failed to read the request body", http.StatusInternalServerError)
			return
		}
		defer r.Body.Close()

		webhookKey := env.GetString("STRIPE_WEBHOOK_KEY", "")
		if webhookKey == "" {
			log.Printf("stripe webhook key is required")
			http.Error(w, "failed to process webhook", http.StatusInternalServerError)
			return
		}

		event, err := webhook.ConstructEventWithOptions(
			body,
			r.Header.Get("Stripe-Signature"),
			webhookKey,
			webhook.ConstructEventOptions{
				IgnoreAPIVersionMismatch: true,
			},
		)

		switch event.Type {
		case "checkout.session.completed":
			var session stripe.CheckoutSession

			err := json.Unmarshal(event.Data.Raw, &session)
			if err != nil {
				log.Printf("Error parsing webhook JSON: %v", err)
				http.Error(w, "Invalid payload", http.StatusBadRequest)
				return
			}

			payload := messaging.PaymentStatusUpdateData{
				TripID:   session.Metadata["trip_id"],
				UserID:   session.Metadata["user_id"],
				DriverID: session.Metadata["driver_id"],
			}

			payloadBytes, err := json.Marshal(payload)
			if err != nil {
				log.Printf("Error marshalling payload: %v", err)
				http.Error(w, "Failed to marshal payload", http.StatusInternalServerError)
				return
			}

			message := contracts.AmqpMessage{
				OwnerID: session.Metadata["user_id"],
				Data:    payloadBytes,
			}

			if err := rb.Publish(
				r.Context(),
				contracts.PaymentEventSuccess,
				message,
			); err != nil {
				log.Printf("Error publishing payment event: %v", err)
				http.Error(w, "Failed to publish payment event", http.StatusInternalServerError)
				return
			}
		}
	}
}
