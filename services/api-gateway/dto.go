package main

import (
	"github.com/atcheri/ride-booking-go/shared/types"
	pb "github.com/atcheri/ride-booking-grpc-proto/golang/trip"
)

type tripPreviewRequest struct {
	UserId      string           `json:"userID"`
	PickUp      types.Coordinate `json:"pickup"`
	Destination types.Coordinate `json:"destination"`
}

func (r *tripPreviewRequest) ToProto() *pb.PreviewTripRequest {
	return &pb.PreviewTripRequest{
		UserID: r.UserId,
		StartLocation: &pb.Coordinate{
			Latitude:  r.PickUp.Latitude,
			Longitude: r.PickUp.Longitude,
		},
		EndLocation: &pb.Coordinate{
			Latitude:  r.Destination.Latitude,
			Longitude: r.Destination.Longitude,
		},
	}
}

type startTripRequest struct {
	RideFareID string `json:"rideFareID"`
	UserID     string `json:"userID"`
}

func (r *startTripRequest) ToProto() *pb.CreateTripRequest {
	return &pb.CreateTripRequest{
		RideFareID: r.RideFareID,
		UserID:     r.UserID,
	}
}
