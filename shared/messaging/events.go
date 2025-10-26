package messaging

import (
	pb "github.com/atcheri/ride-booking-grpc-proto/golang/trip"
)

const (
	FindAvailableDriversQueue = "find_available_drivers"
	DriverCmdTripRequestQueue = "driver_cmd_trip_request"
	DriverTripResponseQueue   = "driver_trip_response"
)

type TripEventData struct {
	Trip *pb.Trip `json:"trip"`
}
