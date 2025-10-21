package service

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/atcheri/ride-booking-go/services/trip-service/internal/domain/models"
	"github.com/atcheri/ride-booking-go/services/trip-service/internal/domain/repository"
	"github.com/atcheri/ride-booking-go/shared/types"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var (
	osrmApiV1Endpoint = "http://router.project-osrm.org/route/v1/driving"
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

func (s *TripService) GetRoute(ctx context.Context, pickup, destination *types.Coordinate) (*types.OsrmApiResponse, error) {
	url := fmt.Sprintf(
		"%s/%f,%f;%f,%f?overview=full&geometries=geojson",
		osrmApiV1Endpoint,
		pickup.Longitude, pickup.Latitude,
		destination.Longitude, destination.Latitude,
	)

	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch route from OSRM API: %v", err)
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read the response: %v", err)
	}

	var routeResp types.OsrmApiResponse
	if err := json.Unmarshal(body, &routeResp); err != nil {
		return nil, fmt.Errorf("failed to parse response: %v", err)
	}

	return &routeResp, nil
}
