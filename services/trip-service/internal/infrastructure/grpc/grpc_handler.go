package grpc

import (
	"context"
	"log"

	"github.com/atcheri/ride-booking-go/services/trip-service/internal/domain/service"
	"github.com/atcheri/ride-booking-go/services/trip-service/internal/infrastructure/dto"
	"github.com/atcheri/ride-booking-go/shared/types"
	pb "github.com/atcheri/ride-booking-grpc-proto/golang/trip"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type gRPCHandler struct {
	pb.UnimplementedTripServiceServer
	service service.TripService
}

func NewGRPCHandler(server *grpc.Server, service service.TripService) {
	handler := &gRPCHandler{
		service: service,
	}

	pb.RegisterTripServiceServer(server, handler)
}

func (h *gRPCHandler) PreviewTrip(ctx context.Context, req *pb.PreviewTripRequest) (*pb.PreviewTripResponse, error) {
	pickup := req.GetStartLocation()
	destination := req.GetEndLocation()
	route, err := h.service.GetRoute(ctx,
		&types.Coordinate{
			Latitude:  pickup.GetLatitude(),
			Longitude: pickup.GetLongitude(),
		}, &types.Coordinate{
			Latitude:  destination.GetLatitude(),
			Longitude: destination.GetLongitude(),
		})

	if err != nil {
		log.Println(err)
		return nil, status.Errorf(codes.Internal, "failed to get route: %v", err)
	}

	estimatedFares := h.service.EstimateRoutePrices(ctx, route)
	fares, err := h.service.PersistTripFares(ctx, estimatedFares, route, req.GetUserID())
	if err != nil {
		log.Println(err)
		return nil, status.Errorf(codes.Internal, "failed to persis fares for route: %v", err)
	}

	return &pb.PreviewTripResponse{
		Route:     route.ToProto(),
		RideFares: dto.FareModelsToProto(fares),
	}, nil
}

func (h *gRPCHandler) CreateTrip(ctx context.Context, req *pb.CreateTripRequest) (*pb.CreateTripResponse, error) {
	fareID := req.GetRideFareID()
	userID := req.GetUserID()
	fare, err := h.service.GetAndValidateFare(ctx, fareID, userID)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to validate the fare: %v", err)
	}

	trip, err := h.service.CreateTrip(ctx, fare)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create the trip: %v", err)
	}

	// TODO: publish an event to the message queue

	return &pb.CreateTripResponse{
		TripID: trip.ID.Hex(),
	}, nil
}
