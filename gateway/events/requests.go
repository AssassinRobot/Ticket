package events

import (
	"context"
	dto "gateway/DTO"
)

type RequestSender interface {
	ListUsers(ctx context.Context) ([]dto.UserDTO, error)
	GetUserByID(ctx context.Context, userID uint) (*dto.UserDTO, error)
	CreateUser(ctx context.Context, firstName, lastName, email string) error
	UpdateUserByID(ctx context.Context, userID uint, firstName, lastName string) error
	DeleteUserByID(ctx context.Context, userID uint) error

	ListTrains(ctx context.Context) ([]dto.TrainDTO, error)
	ListTrainsFiltered(ctx context.Context, trainFilters map[string]string) ([]dto.TrainDTO, error)
	GetTrainByID(ctx context.Context, trainID uint) (*dto.TrainDTO, error)
	CreateTrain(ctx context.Context, name string, capacity uint) error
	UpdateTrainName(ctx context.Context, trainID uint, name string) error
	UpdateTrainTravelDetails(ctx context.Context, trainID uint, origin, destination, departureTime, arrivalTime string) error
	DeleteTrainByID(ctx context.Context, trainID uint) error

	GetSeatByID(ctx context.Context, seatID uint) (*dto.SeatDTO, error)
	ListSeatsByTrainID(ctx context.Context, trainID uint) ([]dto.SeatDTO, error)
	CreateSeat(ctx context.Context, trainID, seatNumber uint) error
	UpdateSeatNumberBySeatID(ctx context.Context, seatID, seatNumber uint) error
	DeleteSeatBySeatID(ctx context.Context, seatID uint) error

	GetTicketByID(ctx context.Context, ticketID uint) (*dto.TicketDTO, error)
	ListTicketsByUserID(ctx context.Context,userID uint)([]dto.TicketDTO, error)
	ListTicketsByTrainID(ctx context.Context,trainID uint)([]dto.TicketDTO,error)
	BookTicket(ctx context.Context, userID, trainID, TicketsNumber uint) ([]dto.TicketDTO, error)
	CancelTicket(ctx context.Context, ticketID uint) error
}
