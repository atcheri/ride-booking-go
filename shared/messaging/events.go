package messaging

import (
	pbd "github.com/atcheri/ride-booking-grpc-proto/golang/driver"
	pb "github.com/atcheri/ride-booking-grpc-proto/golang/trip"
)

const (
	FindAvailableDriversQueue      = "find_available_drivers"
	DriverCmdTripRequestQueue      = "driver_cmd_trip_request"
	DriverTripResponseQueue        = "driver_trip_response"
	NotifyRiderNoDriversFoundQueue = "notify_rider_no_drivers_found"
)

type TripEventData struct {
	Trip *pb.Trip `json:"trip"`
}

type DriverTripResponseData struct {
	Driver  *pbd.Driver `json:"driver"`
	TripID  string      `json:"tripID"`
	RiderID string      `json:"riderID"`
}
