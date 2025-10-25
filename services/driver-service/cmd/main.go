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
	amqp "github.com/rabbitmq/amqp091-go"
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

	// RAbbitMQ connection
	conn, err := amqp.Dial(rabbitmqURI)
	if err != nil {
		log.Fatalf("failed to connect to rabbitmq")
	}
	defer conn.Close()

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
