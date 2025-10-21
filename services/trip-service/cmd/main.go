package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

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

	serverErrorChan := make(chan error, 1)

	go func() {
		log.Printf("trip-service server listening on port %s", httpAddr)
		serverErrorChan <- server.ListenAndServe()
	}()

	shutDownCh := make(chan os.Signal, 1)
	signal.Notify(shutDownCh, os.Interrupt, syscall.SIGTERM)

	select {
	case err := <-serverErrorChan:
		log.Printf("error starting the trip-service server: %v", err)
	case sig := <-shutDownCh:
		log.Printf("trip-service server is shutting down due to %v signal", sig)
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

		defer cancel()
		if err := server.Shutdown(ctx); err != nil {
			log.Printf("failed to gracefully shutdown trip-service server: %v", err)
			server.Close()
		}
	}

}
