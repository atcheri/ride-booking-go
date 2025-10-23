package dto

import (
	"log"

	"github.com/atcheri/ride-booking-go/services/trip-service/internal/domain/models"
	pb "github.com/atcheri/ride-booking-grpc-proto/golang/trip"
)

func FareModelToProto(fare *models.RideFareModel) *pb.RideFare {
	log.Printf("fare model: %v", fare)
	return &pb.RideFare{
		Id:                fare.ID.Hex(),
		UserID:            fare.UserID,
		PackageSlug:       fare.PackageSlug,
		TotalPriceInCents: fare.TotalPriceInCents,
	}
}

func FareModelsToProto(fares []*models.RideFareModel) []*pb.RideFare {
	protoFares := make([]*pb.RideFare, 0)
	for _, fare := range fares {
		protoFares = append(protoFares, FareModelToProto(fare))
	}

	return protoFares
}
