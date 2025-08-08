package nats

import (
	"context"
	dto "gateway/DTO"
	"gateway/utils"
	"strconv"

	"github.com/nats-io/nats.go"
)

const (
	APIGatewayRequestSubjects = "Gateway.Request.>"

	SubjectRequestListUsers      = "request.list.users"
	SubjectRequestGetUserByID    = "request.get.user.byID"
	SubjectRequestCreateUser     = "request.create.user"
	SubjectRequestGetUserTickets = "request.get.user.tickets"
	SubjectRequestUpdateUser     = "request.update.user"
	SubjectRequestDeleteUser     = "request.delete.user"

	SubjectRequestListTrains              = "request.list.trains"
	SubjectRequestListTrainsFiltered      = "request.list.trains.filtered"
	SubjectRequestGetTrainByID            = "request.get.train.byID"
	SubjectRequestCreateTrain             = "request.create.train"
	SubjectRequestUpdateTrain             = "request.update.train"
	SubjectRequestUpdateTrainTravelDetail = "request.update.train.travel.detail"
	SubjectRequestDeleteTrainByID         = "request.delete.train.byID"

	SubjectRequestGetSeatByID        = "request.get.seat.byID"
	SubjectRequestListSeatsByTrainID = "request.list.seats.byTrainID"
	SubjectRequestCreateSeat         = "request.create.seat"
	SubjectRequestUpdateSeatNumber   = "request.update.seat.number"
	SubjectRequestDeleteSeatByID     = "request.delete.seat.byID"

	SubjectRequestGetTicketByID        = "request.get.ticket.byID"
	SubjectRequestBookTicket           = "request.book.ticket"
	SubjectRequestCancelTicket         = "request.cancel.ticket"
	SubjectRequestListTicketsByUserID  = "request.list.tickets.byUserID"
	SubjectRequestListTicketsByTrainID = "request.list.tickets.byTrainID"
)

type NatsRequestSender struct {
	nats *nats.Conn
}

func NewNatsRequestSender(nats *nats.Conn) *NatsRequestSender {
	return &NatsRequestSender{
		nats: nats,
	}
}

func (s *NatsRequestSender) ListUsers(ctx context.Context) ([]dto.UserDTO, error) {
	replay, err := s.nats.RequestWithContext(ctx, SubjectRequestListUsers, nil)
	if err != nil {
		return nil, err
	} else if err := utils.HandleError(replay.Data); err != nil {
		return nil, err
	}

	userDTOs, err := utils.UnmarshalListUsersReplay(replay.Data)
	if err != nil {
		return nil, err
	}

	return userDTOs, nil
}

func (s *NatsRequestSender) GetUserByID(ctx context.Context, userID uint) (*dto.UserDTO, error) {
	userIDstr := strconv.Itoa(int(userID))

	replay, err := s.nats.RequestWithContext(ctx, SubjectRequestGetUserByID, []byte(userIDstr))
	if err != nil {
		return nil, err
	} else if err := utils.HandleError(replay.Data); err != nil {
		return nil, err
	}

	dtoUser, err := utils.UnmarshalUser(replay.Data)
	if err != nil {
		return nil, err
	}

	return dtoUser, nil
}

func (s *NatsRequestSender) CreateUser(ctx context.Context, firstName string, lastName string, email string) error {
	requestData, err := utils.MarshalCreateUserRequest(firstName, lastName, email)
	if err != nil {
		return err
	}

	replay, err := s.nats.RequestWithContext(ctx, SubjectRequestCreateUser, requestData)
	if err != nil {
		return err
	} else if err := utils.HandleError(replay.Data); err != nil {
		return err
	}

	return nil
}

func (s *NatsRequestSender) UpdateUserByID(ctx context.Context, userID uint, firstName string, lastName string) error {
	requestData, err := utils.MarshalUpdateUserRequest(userID, firstName, lastName)
	if err != nil {
		return err
	}

	replay, err := s.nats.RequestWithContext(ctx, SubjectRequestUpdateUser, requestData)
	if err != nil {
		return err
	} else if err := utils.HandleError(replay.Data); err != nil {
		return err
	}

	return nil
}

func (s *NatsRequestSender) DeleteUserByID(ctx context.Context, userID uint) error {
	userIDstr := strconv.Itoa(int(userID))

	replay, err := s.nats.RequestWithContext(ctx, SubjectRequestDeleteUser, []byte(userIDstr))
	if err != nil {
		return err
	} else if err := utils.HandleError(replay.Data); err != nil {
		return err
	}

	return nil
}

func (s *NatsRequestSender) ListTrains(crx context.Context) ([]dto.TrainDTO, error) {
	replay, err := s.nats.RequestWithContext(crx, SubjectRequestListTrains, nil)
	if err != nil {
		return nil, err
	} else if err := utils.HandleError(replay.Data); err != nil {
		return nil, err
	}

	trainDTOs, err := utils.UnmarshalListTrainsReplay(replay.Data)
	if err != nil {
		return nil, err
	}

	return trainDTOs, nil
}

func (s *NatsRequestSender) ListTrainsFiltered(ctx context.Context, trainFilters map[string]string) ([]dto.TrainDTO, error) {
	requestData, err := utils.MarshalListTrainsFilteredRequest(trainFilters)
	if err != nil {
		return nil, err
	}

	replay, err := s.nats.RequestWithContext(ctx, SubjectRequestListTrainsFiltered, requestData)
	if err != nil {
		return nil, err
	} else if err := utils.HandleError(replay.Data); err != nil {
		return nil, err
	}

	trainDTOs, err := utils.UnmarshalListTrainsReplay(replay.Data)
	if err != nil {
		return nil, err
	}

	return trainDTOs, nil
}

func (s *NatsRequestSender) GetTrainByID(ctx context.Context, trainID uint) (*dto.TrainDTO, error) {
	trainIDstr := strconv.Itoa(int(trainID))

	replay, err := s.nats.RequestWithContext(ctx, SubjectRequestGetTrainByID, []byte(trainIDstr))
	if err != nil {
		return nil, err
	} else if err := utils.HandleError(replay.Data); err != nil {
		return nil, err
	}

	dtoTrain, err := utils.UnmarshalTrain(replay.Data)
	if err != nil {
		return nil, err
	}

	return dtoTrain, nil
}

func (s *NatsRequestSender) CreateTrain(ctx context.Context, name string, capacity uint) error {
	requestData, err := utils.MarshalCreateTrainRequest(name, capacity)
	if err != nil {
		return err
	}

	replay, err := s.nats.RequestWithContext(ctx, SubjectRequestCreateTrain, requestData)
	if err != nil {
		return err
	} else if err := utils.HandleError(replay.Data); err != nil {
		return err
	}

	return nil
}

func (s *NatsRequestSender) UpdateTrainName(ctx context.Context, trainID uint, name string) error {
	requestData, err := utils.MarshalUpdateTrainNameRequest(trainID, name)
	if err != nil {
		return err
	}

	replay, err := s.nats.RequestWithContext(ctx, SubjectRequestUpdateTrain, requestData)
	if err != nil {
		return err
	} else if err := utils.HandleError(replay.Data); err != nil {
		return err
	}

	return nil
}

func (s *NatsRequestSender) UpdateTrainTravelDetails(ctx context.Context, trainID uint, origin, destination, departureTime, arrivalTime string) error {

	travelDetailsMap := map[string]string{
		"origin":         origin,
		"destination":    destination,
		"departure_time": departureTime,
		"arrival_time":   arrivalTime,
	}

	requestData, err := utils.MarshalUpdateTrainTravelDetailsRequest(trainID, travelDetailsMap)
	if err != nil {
		return err
	}

	replay, err := s.nats.RequestWithContext(ctx, SubjectRequestUpdateTrainTravelDetail, requestData)
	if err != nil {
		return err
	} else if err := utils.HandleError(replay.Data); err != nil {
		return err
	}

	return nil
}

func (s *NatsRequestSender) DeleteTrainByID(ctx context.Context, trainID uint) error {
	trainIDstr := strconv.Itoa(int(trainID))

	replay, err := s.nats.RequestWithContext(ctx, SubjectRequestDeleteTrainByID, []byte(trainIDstr))
	if err != nil {
		return err
	} else if err := utils.HandleError(replay.Data); err != nil {
		return err
	}

	return nil
}

func (s *NatsRequestSender) GetSeatByID(ctx context.Context, seatID uint) (*dto.SeatDTO, error) {
	seatIDstr := strconv.Itoa(int(seatID))

	replay, err := s.nats.RequestWithContext(ctx, SubjectRequestGetSeatByID, []byte(seatIDstr))
	if err != nil {
		return nil, err
	} else if err := utils.HandleError(replay.Data); err != nil {
		return nil, err
	}

	dtoSeat, err := utils.UnmarshalSeat(replay.Data)
	if err != nil {
		return nil, err
	}

	return dtoSeat, nil
}

func (s *NatsRequestSender) ListSeatsByTrainID(ctx context.Context, trainID uint) ([]dto.SeatDTO, error) {
	trainIDstr := strconv.Itoa(int(trainID))

	replay, err := s.nats.RequestWithContext(ctx, SubjectRequestListSeatsByTrainID, []byte(trainIDstr))
	if err != nil {
		return nil, err
	} else if err := utils.HandleError(replay.Data); err != nil {
		return nil, err
	}

	dtoSeats, err := utils.UnmarshalListSeatsByTrainIDReplay(replay.Data)
	if err != nil {
		return nil, err
	}

	return dtoSeats, nil
}

func (s *NatsRequestSender) CreateSeat(ctx context.Context, trainID uint, seatNumber uint) error {
	requestData, err := utils.MarshalCreateSeatRequest(trainID, seatNumber)
	if err != nil {
		return err
	}

	replay, err := s.nats.RequestWithContext(ctx, SubjectRequestCreateSeat, requestData)
	if err != nil {
		return err
	} else if err := utils.HandleError(replay.Data); err != nil {
		return err
	}

	return nil
}

func (s *NatsRequestSender) UpdateSeatNumberBySeatID(ctx context.Context, seatID uint, seatNumber uint) error {
	requestData, err := utils.MarshalUpdateSeatNumberRequest(seatID, seatNumber)
	if err != nil {
		return err
	}

	replay, err := s.nats.RequestWithContext(ctx, SubjectRequestUpdateSeatNumber, requestData)
	if err != nil {
		return err
	} else if err := utils.HandleError(replay.Data); err != nil {
		return err
	}

	return nil
}

func (s *NatsRequestSender) DeleteSeatBySeatID(ctx context.Context, seatID uint) error {
	seatIDstr := strconv.Itoa(int(seatID))

	replay, err := s.nats.RequestWithContext(ctx, SubjectRequestDeleteSeatByID, []byte(seatIDstr))
	if err != nil {
		return err
	} else if err := utils.HandleError(replay.Data); err != nil {
		return err
	}

	return nil
}

func (s *NatsRequestSender) GetTicketByID(ctx context.Context, ticketID uint) (*dto.TicketDTO, error) {
	ticketIDstr := strconv.Itoa(int(ticketID))

	replay, err := s.nats.RequestWithContext(ctx, SubjectRequestGetTicketByID, []byte(ticketIDstr))
	if err != nil {
		return nil, err
	} else if err := utils.HandleError(replay.Data); err != nil {
		return nil, err
	}

	dtoTicket, err := utils.UnmarshalTicket(replay.Data)
	if err != nil {
		return nil, err
	}

	return dtoTicket, nil
}

func (s *NatsRequestSender) ListTicketsByUserID(ctx context.Context, userID uint) ([]dto.TicketDTO, error) {
	userIDstr := strconv.Itoa(int(userID))

	replay, err := s.nats.RequestWithContext(ctx, SubjectRequestListTicketsByUserID, []byte(userIDstr))
	if err != nil {
		return nil, err
	} else if err := utils.HandleError(replay.Data); err != nil {
		return nil, err
	}

	dtoTicket, err := utils.UnmarshalTickets(replay.Data)
	if err != nil {
		return nil, err
	}

	return dtoTicket, nil
}

func (s *NatsRequestSender) ListTicketsByTrainID(ctx context.Context, trainID uint) ([]dto.TicketDTO, error) {
	trainIDstr := strconv.Itoa(int(trainID))

	replay, err := s.nats.RequestWithContext(ctx, SubjectRequestListTicketsByTrainID, []byte(trainIDstr))
	if err != nil {
		return nil, err
	} else if err := utils.HandleError(replay.Data); err != nil {
		return nil, err
	}

	dtoTicket, err := utils.UnmarshalTickets(replay.Data)
	if err != nil {
		return nil, err
	}

	return dtoTicket, nil
}

func (s *NatsRequestSender) BookTicket(ctx context.Context, userID uint, trainID uint, ticketsNumber uint) ([]dto.TicketDTO, error) {
	requestData, err := utils.MarshalBookTicketRequest(trainID, userID, ticketsNumber)
	if err != nil {
		return nil, err
	}

	replay, err := s.nats.RequestWithContext(ctx, SubjectRequestBookTicket, requestData)
	if err != nil {
		return nil, err
	} else if err := utils.HandleError(replay.Data); err != nil {
		return nil, err
	}

	ticketsDTO, err := utils.UnmarshalTickets(replay.Data)
	if err != nil {
		return nil, err
	}

	return ticketsDTO, nil

}

func (s *NatsRequestSender) CancelTicket(ctx context.Context, ticketID uint) error {
	ticketIDstr := strconv.Itoa(int(ticketID))

	replay, err := s.nats.RequestWithContext(ctx, SubjectRequestCancelTicket, []byte(ticketIDstr))
	if err != nil {
		return err
	} else if err := utils.HandleError(replay.Data); err != nil {
		return err
	}

	return nil
}
