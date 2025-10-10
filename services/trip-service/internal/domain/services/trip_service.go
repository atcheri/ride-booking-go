package services

import (
	"context"

	domain "github.com/atcheri/ride-booking-go/services/trip-service/internal/domain/models"
)

type TripService interface {
	CreateTrip(ctx context.Context, trip *domain.RideFareModel) (*domain.TripModel, error)
}
