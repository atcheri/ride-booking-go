package service

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"

	"github.com/atcheri/ride-booking-go/services/payment-service/internal/domain/service"
	"github.com/atcheri/ride-booking-go/services/payment-service/pkg/types"
)

type paymentService struct {
	paymentProcessor service.PaymentProcessor
}

// NewPaymentService creates a new instance of the payment service
func NewPaymentService(processor service.PaymentProcessor) service.PaymentService {
	return &paymentService{
		paymentProcessor: processor,
	}
}

// CreatePaymentSession creates a new payment session for a trip
func (s *paymentService) CreatePaymentSession(
	ctx context.Context,
	tripID string,
	userID string,
	driverID string,
	amount int64,
	currency string,
) (*types.PaymentIntent, error) {
	metadata := map[string]string{
		"trip_id":   tripID,
		"user_id":   userID,
		"driver_id": driverID,
	}

	sessionID, err := s.paymentProcessor.CreatePaymentSession(ctx, amount, currency, metadata)
	if err != nil {
		return nil, fmt.Errorf("failed to create payment session: %w", err)
	}

	paymentIntent := &types.PaymentIntent{
		ID:              uuid.New().String(),
		TripID:          tripID,
		UserID:          userID,
		DriverID:        driverID,
		Amount:          amount,
		Currency:        currency,
		StripeSessionID: sessionID,
		CreatedAt:       time.Now(),
	}

	return paymentIntent, nil
}
