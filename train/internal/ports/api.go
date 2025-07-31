package ports

import (
	"context"
	"train/internal/application/core/domain"
)

type APIPort interface {
	CreateTrain(ctx context.Context, name string,capacity uint) error
	GetTrainByID(ctx context.Context, ID uint) (*domain.Train, error)
	ListTrains(ctx context.Context) ([]domain.Train, error)
	ListTrainsFiltered(ctx context.Context, filters map[string]string) ([]domain.Train, error)
	UpdateTrain(ctx context.Context, ID uint, name string) error
	UpdateTrainTravelDetails(ctx context.Context, TrainID uint, destination, origin, departureTime, arrivalTime string) error
	DeleteTrain(ctx context.Context, ID uint) error
	CreateSeat(ctx context.Context, trainID uint, seatNumber uint) error
	UpdateSeatNumber(ctx context.Context, ID uint, seatNumber uint) error
	GetSeatByID(ctx context.Context, ID uint) (*domain.Seat, error)
	SeatBooked(ctx context.Context, seatID, trainID,userID uint) error
	CancelSeatBooking(ctx context.Context, seatID, trainID uint) error
	ListSeatsByTrainID(ctx context.Context, trainID uint) ([]domain.Seat, error)
	DeleteSeat(ctx context.Context,ID uint)error
}
