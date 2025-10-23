package service

import (
	"context"

	domain "github.com/atcheri/ride-booking-go/services/trip-service/internal/domain/models"
	"github.com/atcheri/ride-booking-go/services/trip-service/internal/infrastructure/dto"
	"github.com/atcheri/ride-booking-go/shared/types"
)

type TripService interface {
	CreateTrip(ctx context.Context, trip *domain.RideFareModel) (*domain.TripModel, error)
	GetRoute(ctx context.Context, pickup, destination *types.Coordinate) (*dto.OsrmApiResponse, error)
	EstimateRoutePrices(ctx context.Context, route *dto.OsrmApiResponse) []*domain.RideFareModel
	PersistTripFares(ctx context.Context, fares []*domain.RideFareModel, userID string) ([]*domain.RideFareModel, error)
}
