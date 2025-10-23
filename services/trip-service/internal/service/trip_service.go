package service

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/atcheri/ride-booking-go/services/trip-service/internal/domain/models"
	"github.com/atcheri/ride-booking-go/services/trip-service/internal/domain/repository"
	"github.com/atcheri/ride-booking-go/services/trip-service/internal/infrastructure/dto"
	"github.com/atcheri/ride-booking-go/shared/types"
	pb "github.com/atcheri/ride-booking-grpc-proto/golang/trip"
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
		Driver:   &pb.Driver{}, // populating the struct with an empty driver on trip creation
	}

	return s.tripRepository.CreateTrip(ctx, trip)
}

func (s *TripService) GetRoute(ctx context.Context, pickup, destination *types.Coordinate) (*dto.OsrmApiResponse, error) {
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

	var routeResp dto.OsrmApiResponse
	if err := json.Unmarshal(body, &routeResp); err != nil {
		return nil, fmt.Errorf("failed to parse response: %v", err)
	}

	return &routeResp, nil
}

func (s *TripService) EstimateRoutePrices(ctx context.Context, route *dto.OsrmApiResponse) []*models.RideFareModel {
	baseFares := []*models.RideFareModel{
		{
			PackageSlug:       "suv",
			TotalPriceInCents: 200,
		},
		{
			PackageSlug:       "sedan",
			TotalPriceInCents: 350,
		},
		{
			PackageSlug:       "van",
			TotalPriceInCents: 400,
		},
		{
			PackageSlug:       "luxury",
			TotalPriceInCents: 1000,
		},
	}

	estimatedFares := make([]*models.RideFareModel, len(baseFares))
	for i, fare := range baseFares {
		estimatedFares[i] = estimateRouteFare(fare, route)
	}

	return estimatedFares
}

func (s *TripService) PersistTripFares(ctx context.Context, fares []*models.RideFareModel, userID string) ([]*models.RideFareModel, error) {
	faresToPersist := make([]*models.RideFareModel, len(fares))

	for i, f := range fares {
		fare := &models.RideFareModel{
			UserID:            userID,
			ID:                primitive.NewObjectID(),
			TotalPriceInCents: f.TotalPriceInCents,
			PackageSlug:       f.PackageSlug,
		}

		if err := s.tripRepository.SaveTripFare(ctx, fare); err != nil {
			return nil, fmt.Errorf("failed to persis trip fare: %w", err)
		}

		faresToPersist[i] = fare
	}

	return faresToPersist, nil
}

func (s *TripService) GetAndValidateFare(ctx context.Context, fareID, userID string) (*models.RideFareModel, error) {
	fare, err := s.tripRepository.GetFareByID(ctx, fareID)
	if err != nil {
		return nil, fmt.Errorf("failed to get the fare for ID %s : %v", fareID, err)
	}

	// validate the fare
	if fare == nil {
		return nil, fmt.Errorf("fare does not exist")
	}

	// ... and the fare ownership
	if userID != fare.UserID {
		return nil, fmt.Errorf("fare does not belong to the user")
	}

	return fare, nil
}

func estimateRouteFare(fare *models.RideFareModel, route *dto.OsrmApiResponse) *models.RideFareModel {
	pricingCfg := dto.DefaultPricingConfig()
	carPackagePrice := fare.TotalPriceInCents
	distanceKm := route.Routes[0].Distance
	durationInMinutes := route.Routes[0].Duration
	distanceFare := distanceKm * pricingCfg.PricePerUnitOfDistance
	timeFare := durationInMinutes * pricingCfg.PricingPerUnitOftime
	totalPrice := carPackagePrice + distanceFare + timeFare

	return &models.RideFareModel{
		TotalPriceInCents: totalPrice,
		PackageSlug:       fare.PackageSlug,
	}
}
