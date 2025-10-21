package main

import (
	"log"
	"net/http"

	h "github.com/atcheri/ride-booking-go/services/trip-service/internal/infrastructure/http"
	"github.com/atcheri/ride-booking-go/services/trip-service/internal/infrastructure/repository"
	"github.com/atcheri/ride-booking-go/services/trip-service/internal/service"
	"github.com/atcheri/ride-booking-go/shared/env"
)

var (
	httpAddr = env.GetString("HTTP_ADDR", ":8083")
)

func main() {
	inMemoryRepository := repository.NewInMemoryRepository()
	tripService := service.NewTripService(inMemoryRepository)
	httphandler := h.HttpHandler{Service: tripService}

	log.Printf("Starting API Gateway on port: %s", httpAddr)

	mux := http.NewServeMux()

	mux.HandleFunc("POST /preview", httphandler.HandleTripPreview)

	server := &http.Server{
		Addr:    httpAddr,
		Handler: mux,
	}

	if err := server.ListenAndServe(); err != nil {
		log.Printf("TRIP SERVICE HTTP server error: %v", err)
	}
}
