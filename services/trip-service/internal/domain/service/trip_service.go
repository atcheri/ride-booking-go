package service

import (
	"context"

	"github.com/atcheri/ride-booking-go/services/trip-service/internal/domain/models"
	"github.com/atcheri/ride-booking-go/services/trip-service/internal/infrastructure/dto"
	"github.com/atcheri/ride-booking-go/shared/types"
	pb "github.com/atcheri/ride-booking-grpc-proto/golang/driver"
)

type TripService interface {
	CreateTrip(ctx context.Context, trip *models.RideFareModel) (*models.TripModel, error)
	GetRoute(ctx context.Context, pickup, destination *types.Coordinate) (*dto.OsrmApiResponse, error)
	EstimateRoutePrices(ctx context.Context, route *dto.OsrmApiResponse) []*models.RideFareModel
	PersistTripFares(
		ctx context.Context,
		fares []*models.RideFareModel,
		route *dto.OsrmApiResponse,
		userID string,
	) ([]*models.RideFareModel, error)
	GetAndValidateFare(ctx context.Context, fareID, userID string) (*models.RideFareModel, error)
	GetTripByID(ctx context.Context, id string) (*models.TripModel, error)
	UpdateTrip(ctx context.Context, id string, status string, driver *pb.Driver) error
}
