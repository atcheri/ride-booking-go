package services

import (
	"context"

	"github.com/atcheri/ride-booking-go/services/trip-service/internal/domain/models"
	"github.com/atcheri/ride-booking-go/services/trip-service/internal/domain/repositories"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type TripService struct {
	tripRepository repositories.TripRepository
}

func NewTripService(repo repositories.TripRepository) *TripService {
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
