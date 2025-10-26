package repository

import (
	"context"

	"github.com/atcheri/ride-booking-go/services/trip-service/internal/domain/models"
	pb "github.com/atcheri/ride-booking-grpc-proto/golang/driver"
)

type TripRepository interface {
	CreateTrip(ctx context.Context, trip *models.TripModel) (*models.TripModel, error)
	SaveTripFare(ctx context.Context, fare *models.RideFareModel) error
	GetFareByID(ctx context.Context, fareID string) (*models.RideFareModel, error)
	GetTripByID(ctx context.Context, id string) (*models.TripModel, error)
	UpdateTrip(ctx context.Context, tripID string, status string, driver *pb.Driver) error
}
