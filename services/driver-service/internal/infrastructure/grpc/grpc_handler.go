package grpc

import (
	"context"

	"github.com/atcheri/ride-booking-go/services/driver-service/domain/service"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "github.com/atcheri/ride-booking-grpc-proto/golang/driver"
)

type gRPCHandler struct {
	pb.UnimplementedDriverServiceServer
	service service.DriverService
}

func NewGrpcHandler(server *grpc.Server, service service.DriverService) *gRPCHandler {
	handler := &gRPCHandler{
		service: service,
	}

	pb.RegisterDriverServiceServer(server, handler)

	return handler
}

func (h *gRPCHandler) RegisterDriver(ctx context.Context, req *pb.RegisterDriverRequest) (*pb.RegisterDriverResponse, error) {
	driver, err := h.service.RegisterDriver(req.GetDriverID(), req.GetPackageSlug())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to register driver: %v", err)
	}

	return &pb.RegisterDriverResponse{
		Driver: driver,
	}, nil
}

func (h *gRPCHandler) UnregisterDriver(ctx context.Context, req *pb.RegisterDriverRequest) (*pb.RegisterDriverResponse, error) {
	h.service.UnregisterDriver(req.GetDriverID())

	return &pb.RegisterDriverResponse{
		Driver: &pb.Driver{
			Id: req.GetDriverID(),
		},
	}, nil
}
