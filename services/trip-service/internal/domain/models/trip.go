package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"

	pb "github.com/atcheri/ride-booking-grpc-proto/golang/trip"
)

type TripModel struct {
	ID       primitive.ObjectID
	UserId   string
	Status   string
	RideFare *RideFareModel
	Driver   *pb.Driver
}
