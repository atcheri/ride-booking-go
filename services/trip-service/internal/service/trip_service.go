package service

import (
	"context"

	"github.com/atcheri/ride-booking-go/services/trip-service/internal/domain/models"
	"github.com/atcheri/ride-booking-go/services/trip-service/internal/domain/repository"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type TripService struct {
	tripRepository repository.TripRepository
}

func NewTripService(repo repository.TripRepository) *TripService {
	return &TripService{
		tripRepository: repo,
	}
}

func (s *TripService) CreateTrip(ctx context.Context, ride *models.RideFareModel) (*models.TripModel, error) {
	trip := &models.TripModel{
		ID:       primitive.NewObjectID(),
		UserId:   ride.UserID,
		Status:   "pending",
		RideFare: ride,
	}

	return s.tripRepository.CreateTrip(ctx, trip)

}
