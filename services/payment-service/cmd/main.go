package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/atcheri/ride-booking-go/services/payment-service/internal/infrastructure/payment"
	"github.com/atcheri/ride-booking-go/services/payment-service/internal/service"
	"github.com/atcheri/ride-booking-go/services/payment-service/pkg/types"
	"github.com/atcheri/ride-booking-go/shared/env"
	"github.com/atcheri/ride-booking-go/shared/messaging"
)

var (
	rabbitmqURI     = env.GetString("RABBITMQ_DEFAULT_URI", "amqp://guest:guest@rabbitmq:56723/")
	appURL          = env.GetString("APP_URL", "http://localhost:3000")
	stripeSecretKey = env.GetString("STRIPE_SECRET_KEY", "")
	successURL      = env.GetString("STRIPE_SUCCESS_URL", appURL+"?payment=success")
	cancelURL       = env.GetString("STRIPE_CANCEL_URL", appURL+"?payment=cancel")
)

func main() {
	// Setup graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

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

	paymentProcessor := payment.NewStripeClient(stripeCfg)
	paymentService := service.NewPaymentService(paymentProcessor)
	// FIXME: remove this later
	log.Println(paymentService)

	// RabbitMQ connection
	rabbitmq, err := messaging.NewRabbitMQ(rabbitmqURI)
	if err != nil {
		log.Fatal(err)
	}
	defer rabbitmq.Close()

	log.Println("payment-service connected to RabbitMQ")

	// Wait for shutdown signal
	<-ctx.Done()
	log.Println("Shutting down payment service...")
}
