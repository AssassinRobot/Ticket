package ports

import (
	"context"
	"train/internal/application/core/domain"
)

type DatabasePort interface {
	CreateTrain(ctx context.Context, train *domain.Train) error

	GetTrainByID(ctx context.Context, ID uint) (*domain.Train, error)
	ListTrains(ctx context.Context) ([]domain.Train, error)
	ListTrainsFiltered(ctx context.Context, filter *domain.TrainFilters) ([]domain.Train, error)
	IsTrainAvailable(ctx context.Context, trainID uint) (bool, error)

	UpdateTrain(ctx context.Context, ID uint, name string) error
	UpdateTrainTravelDetails(ctx context.Context, trainID uint, travelDetails *domain.TrainTravelDetails) error
	MinusTrainAvailableSeats(ctx context.Context, trainID uint) error
	PlusTrainAvailableSeats(ctx context.Context, trainID uint) error

	DeleteTrain(ctx context.Context, ID uint) error

	CreateSeat(ctx context.Context, seat *domain.Seat) error
	UpdateSeatNumber(ctx context.Context, ID uint, seatNumber uint) error
	GetSeatByID(ctx context.Context, ID uint) (*domain.Seat, error)
	ListSeatsByTrainID(ctx context.Context, trainID uint) ([]domain.Seat, error)
	UpdateSeatBookingStatus(ctx context.Context, seatID uint, booked bool) error
	UpdateSeatUser(ctx context.Context, seatID,userID uint) error
	DeleteSeat(ctx context.Context, ID uint) error
}
