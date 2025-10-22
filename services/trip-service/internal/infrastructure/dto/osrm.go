package dto

import (
	"github.com/atcheri/ride-booking-grpc-proto/golang/trip"
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

func (o *OsrmApiResponse) ToProto() *trip.Route {
	route := o.Routes[0]
	geometry := route.Geometry.Coordinates
	coordinates := make([]*trip.Coordinate, len(geometry))
	for i, coord := range geometry {
		coordinates[i] = &trip.Coordinate{
			Latitude:  coord[0],
			Longitude: coord[1],
		}
	}

	return &trip.Route{
		Geometry: []*trip.Geometry{
			{
				Coordinates: coordinates,
			},
		},
		Distance: route.Distance,
		Duration: route.Duration,
	}
}
