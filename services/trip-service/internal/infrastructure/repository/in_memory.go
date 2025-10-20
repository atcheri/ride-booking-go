package repository

import (
	"context"

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
