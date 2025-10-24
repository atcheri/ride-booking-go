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
	grpcserver "google.golang.org/grpc"
)

var (
	GrpcAddr = env.GetString("HTTP_ADDR", ":9092")
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

	log.Printf("Starting the driver-service grpc server on port: %s", GrpcAddr)

	lis, err := net.Listen("tcp", GrpcAddr)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

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
