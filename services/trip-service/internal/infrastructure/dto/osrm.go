package dto

import (
	"github.com/atcheri/ride-booking-go/services/trip-service/internal/domain/models"
	pb "github.com/atcheri/ride-booking-grpc-proto/golang/trip"
)

type OsrmApiResponse struct {
	Routes []struct {
		Distance float64 `json:"distance"`
		Duration float64 `json:"duration"`
		Geometry struct {
			Coordinates [][]float64 `json:"coordinates"`
		} `json:"geometry"`
	} `json:"routes"`
}

func (o *OsrmApiResponse) ToProto() *pb.Route {
	route := o.Routes[0]
	geometry := route.Geometry.Coordinates
	coordinates := make([]*pb.Coordinate, len(geometry))
	for i, coord := range geometry {
		coordinates[i] = &pb.Coordinate{
			Latitude:  coord[0],
			Longitude: coord[1],
		}
	}

	return &pb.Route{
		Geometry: []*pb.Geometry{
			{
				Coordinates: coordinates,
			},
		},
		Distance: route.Distance,
		Duration: route.Duration,
	}
}

func (o *OsrmApiResponse) ToDomain() *models.TripRoute {
	route := o.Routes[0]
	geometry := route.Geometry.Coordinates
	coordinates := make([][]float64, len(geometry))
	for i, coord := range geometry {
		coordinates[i] = []float64{coord[0], coord[1]}
	}
	return &models.TripRoute{
		Distance: route.Distance,
		Duration: route.Duration,
		Geometry: struct {
			Coordinates [][]float64
		}{
			Coordinates: coordinates,
		},
	}
}
