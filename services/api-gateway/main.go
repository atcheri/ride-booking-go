package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/atcheri/ride-booking-go/shared/env"
	"github.com/atcheri/ride-booking-go/shared/messaging"
)

var (
	httpAddr    = env.GetString("HTTP_ADDR", ":8081")
	rabbitmqURI = env.GetString("RABBITMQ_DEFAULT_URI", "amqp://guest:guest@rabbitmq:56723/")
)

func main() {
	log.Printf("Starting API Gateway on port: %s", httpAddr)

	// RabbitMQ connection
	rabbitMQ, err := messaging.NewRabbitMQ(rabbitmqURI)
	if err != nil {
		log.Fatal(err)
	}
	defer rabbitMQ.Close()

	log.Println("api-gateway connected to RabbitMQ")

	mux := http.NewServeMux()

	mux.HandleFunc("GET /hello", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Hello from API Gateway"))
	})

	mux.HandleFunc("POST /trip/preview", enableCors(handleTripPreview))
	mux.HandleFunc("POST /trip/start", enableCors(handleStartTrip))
	mux.HandleFunc("/ws/drivers", handleDriversWebSocketWithRabbitMQ(rabbitMQ))
	mux.HandleFunc("/ws/riders", handleRidersWebSocketWithRabbitMQ(rabbitMQ))

	server := &http.Server{
		Addr:    httpAddr,
		Handler: mux,
	}

	serverErrorCh := make(chan error, 1)

	go func() {
		log.Printf("APIGateway listening on port %s", httpAddr)
		serverErrorCh <- server.ListenAndServe()
	}()

	shutDownCh := make(chan os.Signal, 1)
	signal.Notify(shutDownCh, os.Interrupt, syscall.SIGTERM)

	select {
	case err := <-serverErrorCh:
		log.Printf("error starting the APIGateway: %v", err)
	case sig := <-shutDownCh:
		log.Printf("APIGateway is shutting down due to %v signal", sig)
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		if err := server.Shutdown(ctx); err != nil {
			log.Printf("failed to gracefully shutdown APIGateway: %v", err)
			server.Close()
		}
	}

}
