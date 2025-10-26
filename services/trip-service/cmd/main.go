package main

import (
	"context"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/atcheri/ride-booking-go/services/trip-service/internal/infrastructure/events"
	"github.com/atcheri/ride-booking-go/services/trip-service/internal/infrastructure/grpc"
	"github.com/atcheri/ride-booking-go/services/trip-service/internal/infrastructure/repository"
	"github.com/atcheri/ride-booking-go/services/trip-service/internal/service"
	"github.com/atcheri/ride-booking-go/shared/env"
	"github.com/atcheri/ride-booking-go/shared/messaging"
	grpcserver "google.golang.org/grpc"
)

var (
	gRPCAddr    = env.GetString("HTTP_ADDR", ":9093")
	rabbitmqURI = env.GetString("RABBITMQ_DEFAULT_URI", "amqp://guest:guest@rabbitmq:56723/")
)

func main() {
	inMemoryRepository := repository.NewInMemoryRepository()
	tripService := service.NewTripService(inMemoryRepository)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// this go routine catches the sigterm or interrup signals and calls the context cancel function
	// that will allow the gracefull shutdown at the end
	go func() {
		sigCh := make(chan os.Signal, 1)
		signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)
		<-sigCh
		cancel()
	}()

	log.Printf("Starting the trip-service grpc server on port: %s", gRPCAddr)

	lis, err := net.Listen("tcp", gRPCAddr)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	// RAbbitMQ connection
	rabbitMQ, err := messaging.NewRabbitMQ(rabbitmqURI)
	if err != nil {
		log.Fatal(err)
	}
	defer rabbitMQ.Close()

	log.Println("trip-service connected to RabbitMQ")

	publisher := events.NewTripEventPublisher(rabbitMQ)

	// starting the gRPc server
	grpcServer := grpcserver.NewServer( /*OPTIONS*/ )
	grpc.NewGRPCHandler(grpcServer, tripService, publisher)

	log.Printf("starting gRPC trip-service on port: %s", lis.Addr().String())

	go func() {
		if err := grpcServer.Serve(lis); err != nil {
			log.Printf("failed to serve the gRPC server: %v", err)
			cancel()
		}
	}()

	// wait for the shutdown signal
	<-ctx.Done()
	log.Println("shutting down the gRPC trip-service server")
	grpcServer.GracefulStop()
}
