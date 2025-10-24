package service

import (
	pb "github.com/atcheri/ride-booking-grpc-proto/golang/driver"
)

type DriverService interface {
	RegisterDriver(driverId string, packageSlug string) (*pb.Driver, error)
	UnregisterDriver(driverId string)
}
