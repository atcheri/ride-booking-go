package service

import (
	"context"

	"github.com/atcheri/ride-booking-go/services/payment-service/pkg/types"
)

type PaymentService interface {
	CreatePaymentSession(ctx context.Context, tripID, userID, driverID string, amount int64, currency string) (*types.PaymentIntent, error)
}

type PaymentProcessor interface {
	CreatePaymentSession(ctx context.Context, amount int64, currency string, metadata map[string]string) (string, error)
}
