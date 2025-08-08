package ports

import (
	"context"
	"ticket/internal/application/core/domain"
)

type EventPublisherPort interface {
	PublishSeatBooked(ctx context.Context, seat *domain.Seat) error
	PublishSeatBookingCanceled(ctx context.Context, seat *domain.Seat) error
}

type RequestPort interface {
	RequestGetTrainByID(ctx context.Context, trainID uint) (*domain.Train, error)
}