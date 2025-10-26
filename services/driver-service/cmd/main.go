package main

import (
	"context"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/atcheri/ride-booking-go/services/driver-service/internal/infrastructure/grpc"
	"github.com/atcheri/ride-booking-go/services/driver-service/internal/service"
	"github.com/atcheri/ride-booking-go/shared/env"
	"github.com/atcheri/ride-booking-go/shared/messaging"
	grpcserver "google.golang.org/grpc"
)

var (
	gRPCAddr    = env.GetString("HTTP_ADDR", ":9092")
	rabbitmqURI = env.GetString("RABBITMQ_DEFAULT_URI", "amqp://guest:guest@rabbitmq:56723/")
)

func main() {
	driverService := service.NewDriverService()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		sigCh := make(chan os.Signal, 1)
		signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)
		<-sigCh
		cancel()
	}()

	log.Printf("Starting the driver-service grpc server on port: %s", gRPCAddr)

	lis, err := net.Listen("tcp", gRPCAddr)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	// RabbitMQ connection
	rabbitMQ, err := messaging.NewRabbitMQ(rabbitmqURI)
	if err != nil {
		log.Fatal(err)
	}
	defer rabbitMQ.Close()

	log.Println("drver-service connected to RabbitMQ")

	consumer := service.NewTripConsumerService(rabbitMQ, driverService)
	go func() {
		if err := consumer.Listen(); err != nil {
			log.Fatalf("failed to listen to rabbitmq: %v", err)
		}
	}()

	grpcServer := grpcserver.NewServer( /*OPTIONS*/ )
	grpc.NewGrpcHandler(grpcServer, driverService)

	log.Printf("starting gRPC driver-service on port: %s", lis.Addr().String())

	go func() {
		if err := grpcServer.Serve(lis); err != nil {
			log.Printf("failed to serve the gRPC server: %v", err)
			cancel()
		}
	}()

	<-ctx.Done()
	log.Println("shutting down the gRPC driver-service server")
	grpcServer.GracefulStop()
}
