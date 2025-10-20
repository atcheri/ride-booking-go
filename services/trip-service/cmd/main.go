package main

import (
	"context"
	"log"
	"time"

	"github.com/atcheri/ride-booking-go/services/trip-service/internal/domain/models"
	"github.com/atcheri/ride-booking-go/services/trip-service/internal/infrastructure/repository"
	"github.com/atcheri/ride-booking-go/services/trip-service/internal/service"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func main() {
	ctx := context.Background()
	inMemoryRepository := repository.NewInMemoryRepository()
	tripService := service.NewTripService(inMemoryRepository)

	fare := &models.RideFareModel{
		ID:     primitive.NewObjectID(),
		UserID: "fake-user-id",
	}

	t, err := tripService.CreateTrip(ctx, fare)
	if err != nil {
		log.Println(err)
		return
	}

	log.Println(t)

	for {
		time.Sleep(time.Second)
	}

}
