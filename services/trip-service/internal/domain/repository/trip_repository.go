package repository

import (
	"context"

	"github.com/atcheri/ride-booking-go/services/trip-service/internal/domain/models"
)

type TripRepository interface {
	CreateTrip(ctx context.Context, trip *models.TripModel) (*models.TripModel, error)
	SaveTripFare(ctx context.Context, fare *models.RideFareModel) error
}
