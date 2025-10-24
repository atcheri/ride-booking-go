package service

import (
	math "math/rand/v2"
	"sync"

	"github.com/icrowley/fake"
	"github.com/mmcloughlin/geohash"

	"github.com/atcheri/ride-booking-go/services/driver-service/internal/utils"
	sharedUtils "github.com/atcheri/ride-booking-go/shared/util"
	pb "github.com/atcheri/ride-booking-grpc-proto/golang/driver"
)

type driverInMap struct {
	Driver *pb.Driver
	// Index int
	// TODO: add a route
}

type DriverService struct {
	drivers []*driverInMap
	mu      sync.RWMutex
}

func NewDriverService() *DriverService {
	return &DriverService{
		drivers: make([]*driverInMap, 0),
	}
}

func (s *DriverService) RegisterDriver(driverId string, packageSlug string) (*pb.Driver, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	randomIndex := math.IntN(len(utils.PredefinedRoutes))
	randomRoute := utils.PredefinedRoutes[randomIndex]

	geohash := geohash.Encode(randomRoute[0][0], randomRoute[0][1])

	driver := s.generateFakeDriver(driverId, randomIndex, packageSlug, geohash, &pb.Location{
		Latitude:  randomRoute[0][0],
		Longitude: randomRoute[0][1],
	})

	s.drivers = append(s.drivers, &driverInMap{Driver: driver})

	return driver, nil
}

func (s *DriverService) UnregisterDriver(driverId string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	for i, driver := range s.drivers {
		if driver.Driver.Id == driverId {
			s.drivers = append(s.drivers[:i], s.drivers[i+1:]...)
		}
	}
}

func (s *DriverService) generateFakeDriver(id string, i int, slug string, geohash string, location *pb.Location) *pb.Driver {
	return &pb.Driver{
		Id:             id,
		Name:           fake.FullName(),
		ProfilePicture: sharedUtils.GetRandomAvatar(i),
		CarPlate:       utils.GenerateRandomPlate(),
		Geohash:        geohash,
		PackageSlug:    slug,
		Location:       location,
	}
}
