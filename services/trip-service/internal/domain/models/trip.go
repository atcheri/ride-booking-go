package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"

	pb "github.com/atcheri/ride-booking-grpc-proto/golang/trip"
)

type TripModel struct {
	ID       primitive.ObjectID `bson:"_id,omitempty"`
	UserID   string             `bson:"userID"`
	Status   string             `bson:"status"`
	RideFare *RideFareModel     `bson:"rideFare"`
	Driver   *pb.TripDriver     `bson:"driver"`
}

func (t *TripModel) ToProto() *pb.Trip {
	return &pb.Trip{
		Id:           t.ID.Hex(),
		UserID:       t.UserID,
		SelectedFare: t.RideFare.ToProto(),
		Route:        t.RideFare.Route.ToProto(),
		Status:       t.Status,
		Driver:       t.Driver,
	}
}
