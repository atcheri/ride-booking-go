package models

import (
	pb "github.com/atcheri/ride-booking-grpc-proto/golang/trip"
)

type TripRoute struct {
	Distance float64
	Duration float64
	Geometry struct {
		Coordinates [][]float64
	}
}

func (r *TripRoute) ToProto() *pb.Route {
	geometry := r.Geometry.Coordinates
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
		Distance: r.Distance,
		Duration: r.Duration,
	}
}
