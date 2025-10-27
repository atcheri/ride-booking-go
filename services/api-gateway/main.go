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
	"github.com/atcheri/ride-booking-go/shared/tracing"
)

var (
	serviceName    = "api-gateway"
	environment    = env.GetString("ENVIRONMENT", "development")
	jaegerEndpoint = env.GetString("JAEGER_ENDPOINT", "http://jaeger:14268/api/traces")
	httpAddr       = env.GetString("HTTP_ADDR", ":8081")
	rabbitmqURI    = env.GetString("RABBITMQ_DEFAULT_URI", "amqp://guest:guest@rabbitmq:56723/")
)

func main() {
	log.Printf("Starting API Gateway on port: %s", httpAddr)

	// Initialize Tracing
	tracerConfig := tracing.NewConfig(serviceName, environment, jaegerEndpoint)
	shutDownTracer, err := tracing.InitTracer(tracerConfig)
	if err != nil {
		log.Fatalf("failed to initialize the tracer: %v", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	defer shutDownTracer(ctx)

	// RabbitMQ connection
	rabbitMQ, err := messaging.NewRabbitMQ(rabbitmqURI)
	if err != nil {
		log.Fatal(err)
	}
	defer rabbitMQ.Close()

	log.Println("api-gateway connected to RabbitMQ")

	mux := http.NewServeMux()

	mux.Handle("GET /hello", tracing.WrapHandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Hello from API Gateway"))
	}, "/hell"))
	mux.Handle("POST /trip/preview", tracing.WrapHandlerFunc(enableCors(handleTripPreview), "/trip/preview"))
	mux.Handle("POST /trip/start", tracing.WrapHandlerFunc(enableCors(handleStartTrip), "/trip/start"))
	mux.Handle("/ws/drivers", tracing.WrapHandlerFunc(handleDriversWebSocketWithRabbitMQ(rabbitMQ), "/ws/drivers"))
	mux.Handle("/ws/riders", tracing.WrapHandlerFunc(handleRidersWebSocketWithRabbitMQ(rabbitMQ), "/ws/riders"))
	mux.Handle("/webhook/stripe", tracing.WrapHandlerFunc(handleStripWebhookWithRabbitMQ(rabbitMQ), "/webhook/stripe"))

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
