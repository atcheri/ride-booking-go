package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"

	pb "github.com/atcheri/ride-booking-grpc-proto/golang/trip"
)

type RideFareModel struct {
	ID                primitive.ObjectID `bson:"_id,omitempty"`
	UserID            string             `bson:"userID"`
	PackageSlug       string             `bson:"packageSlug"`
	TotalPriceInCents float64            `bson:"totalPriceInCents"`
	ExpiresAt         time.Time          `bson:"expiresAt"`
	Route             *TripRoute         `bson:"route"`
}

func (r *RideFareModel) ToProto() *pb.RideFare {
	return &pb.RideFare{
		Id:                r.ID.Hex(),
		UserID:            r.UserID,
		PackageSlug:       r.PackageSlug,
		TotalPriceInCents: r.TotalPriceInCents,
	}
}
