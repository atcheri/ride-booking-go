package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"

	pb "github.com/atcheri/ride-booking-grpc-proto/golang/trip"
)

type RideFareModel struct {
	ID                primitive.ObjectID
	UserID            string
	PackageSlug       string
	TotalPriceInCents float64
	ExpiresAt         time.Time
	Route             *TripRoute
}

func (r *RideFareModel) ToProto() *pb.RideFare {
	return &pb.RideFare{
		Id:                r.ID.Hex(),
		UserID:            r.UserID,
		PackageSlug:       r.PackageSlug,
		TotalPriceInCents: r.TotalPriceInCents,
	}
}
