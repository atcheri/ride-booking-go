package repository

import (
	"context"
	"fmt"

	"github.com/atcheri/ride-booking-go/services/trip-service/internal/domain/models"
)

type inMemoryRepository struct {
	trips     map[string]*models.TripModel
	rideFares map[string]*models.RideFareModel
}

func NewInMemoryRepository() *inMemoryRepository {
	return &inMemoryRepository{
		trips:     make(map[string]*models.TripModel),
		rideFares: make(map[string]*models.RideFareModel),
	}
}

func (r *inMemoryRepository) CreateTrip(ctx context.Context, trip *models.TripModel) (*models.TripModel, error) {
	r.trips[trip.ID.Hex()] = trip

	return trip, nil
}

func (r *inMemoryRepository) SaveTripFare(ctx context.Context, fare *models.RideFareModel) error {
	r.rideFares[fare.ID.Hex()] = fare

	return nil
}

func (r *inMemoryRepository) GetFareByID(ctx context.Context, fareID string) (*models.RideFareModel, error) {
	fare, ok := r.rideFares[fareID]
	if !ok {
		return nil, fmt.Errorf("fare with id %s not found in the in-memory DB", fareID)
	}

	return fare, nil
}
