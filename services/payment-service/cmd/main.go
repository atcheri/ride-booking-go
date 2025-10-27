package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/atcheri/ride-booking-go/services/payment-service/internal/events"
	"github.com/atcheri/ride-booking-go/services/payment-service/internal/infrastructure/payment"
	"github.com/atcheri/ride-booking-go/services/payment-service/internal/service"
	"github.com/atcheri/ride-booking-go/services/payment-service/pkg/types"
	"github.com/atcheri/ride-booking-go/shared/env"
	"github.com/atcheri/ride-booking-go/shared/messaging"
	"github.com/atcheri/ride-booking-go/shared/tracing"
)

var (
	serviceName     = "payment-service"
	environment     = env.GetString("ENVIRONMENT", "development")
	jaegerEndpoint  = env.GetString("JAEGER_ENDPOINT", "http://jaeger:14268/api/traces")
	rabbitmqURI     = env.GetString("RABBITMQ_DEFAULT_URI", "amqp://guest:guest@rabbitmq:56723/")
	appURL          = env.GetString("APP_URL", "http://localhost:3000")
	stripeSecretKey = env.GetString("STRIPE_SECRET_KEY", "")
	successURL      = env.GetString("STRIPE_SUCCESS_URL", appURL+"?payment=success")
	cancelURL       = env.GetString("STRIPE_CANCEL_URL", appURL+"?payment=cancel")
)

func main() {
	// Initialize Tracing
	tracerConfig := tracing.NewConfig(serviceName, environment, jaegerEndpoint)
	shutDownTracer, err := tracing.InitTracer(tracerConfig)
	if err != nil {
		log.Fatalf("failed to initialize the tracer: %v", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	defer shutDownTracer(ctx)

	// Setup graceful shutdown

	go func() {
		sigCh := make(chan os.Signal, 1)
		signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)
		<-sigCh
		cancel()
	}()

	// Stripe config
	stripeCfg := &types.PaymentConfig{
		StripeSecretKey: stripeSecretKey,
		SuccessURL:      successURL,
		CancelURL:       cancelURL,
	}

	if stripeCfg.StripeSecretKey == "" {
		log.Fatalf("STRIPE_SECRET_KEY is not set")
		return
	}

	// RabbitMQ connection
	rabbitmq, err := messaging.NewRabbitMQ(rabbitmqURI)
	if err != nil {
		log.Fatal(err)
	}
	defer rabbitmq.Close()

	log.Println("payment-service connected to RabbitMQ")

	paymentProcessor := payment.NewStripeClient(stripeCfg)
	paymentService := service.NewPaymentService(paymentProcessor)

	// start the trip-payment-consumer
	tripPaymentConsumer := events.NewTripConsumer(rabbitmq, paymentService)
	go tripPaymentConsumer.Listen()

	// Wait for shutdown signal
	<-ctx.Done()
	log.Println("Shutting down payment service...")
}
