package grpc

import (
	"context"
	"log"
	"time"

	"github.com/atcheri/ride-booking-go/services/trip-service/internal/domain/models"
	"github.com/atcheri/ride-booking-go/services/trip-service/internal/domain/service"
	"github.com/atcheri/ride-booking-go/services/trip-service/internal/infrastructure/dto"
	"github.com/atcheri/ride-booking-go/shared/types"
	pb "github.com/atcheri/ride-booking-grpc-proto/golang/trip"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type gRPCHandler struct {
	pb.UnimplementedTripServiceServer
	service service.TripService
}

func NewGRPCHandler(server *grpc.Server, service service.TripService) *gRPCHandler {
	handler := &gRPCHandler{
		service: service,
	}

	pb.RegisterTripServiceServer(server, handler)

	return handler
}

func (h *gRPCHandler) PreviewTrip(ctx context.Context, req *pb.PreviewTripRequest) (*pb.PreviewTripResponse, error) {
	pickup := req.GetStartLocation()
	destination := req.GetEndLocation()
	resp, err := h.service.GetRoute(ctx,
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

	estimatedFares := h.service.EstimateRoutePrices(ctx, resp)
	fares, err := h.service.PersistTripFares(ctx, estimatedFares, req.GetUserID())
	if err != nil {
		log.Println(err)
		return nil, status.Errorf(codes.Internal, "failed to persis fares for route: %v", err)
	}

	return &pb.PreviewTripResponse{
		Route:     resp.ToProto(),
		RideFares: dto.FareModelsToProto(fares),
	}, nil
}

func (h *gRPCHandler) CreateTrip(ctx context.Context, req *pb.CreateTripRequest) (*pb.CreateTripResponse, error) {
	resp, err := h.service.CreateTrip(ctx, &models.RideFareModel{
		ID:                primitive.ObjectID{},
		UserID:            "",
		PackageSlug:       "",
		TotalPriceInCents: 0,
		ExpiresAt:         time.Time{},
	})

	if err != nil {
		log.Println(err)
		return nil, status.Errorf(codes.Internal, "failed to create trip: %v", err)
	}

	return &pb.CreateTripResponse{
		TripID: resp.ID.Hex(),
	}, nil
}
