package messaging

import (
	pb "github.com/atcheri/ride-booking-grpc-proto/golang/trip"
)

const (
	FindAvailableDriversQueue = "find_available_drivers"
)

type TripEventData struct {
	Trip *pb.Trip `json:"trip"`
}
