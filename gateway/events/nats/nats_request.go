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

	SubjectRequestListUsers = "request.list.users"
	SubjectRequestGetUserByID = "request.get.user.byID"
	SubjectRequestCreateUser = "request.create.user"
	SubjectRequestGetUserTickets = "request.get.user.tickets"
	SubjectRequestUpdateUser = "request.update.user"
	SubjectRequestDeleteUser = "request.delete.user"
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

	userDTOs,err := utils.UnmarshalListUsersReplay(replay.Data)
	if err != nil {
		return nil, err
	}

	return userDTOs,nil
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
	panic("not implemented") // TODO: Implement
}

func (s *NatsRequestSender) ListTrainsFiltered(ctx context.Context, travelDetails *dto.TravelDetailDTO) ([]dto.TrainDTO, error) {
	panic("not implemented") // TODO: Implement
}

func (s *NatsRequestSender) GetTrainByID(ctx context.Context) (*dto.TrainDTO, error) {
	panic("not implemented") // TODO: Implement
}

func (s *NatsRequestSender) CreateTrain(ctx context.Context, name string, capacity uint) error {
	panic("not implemented") // TODO: Implement
}

func (s *NatsRequestSender) UpdateTrainName(ctx context.Context, name string) error {
	panic("not implemented") // TODO: Implement
}

func (s *NatsRequestSender) UpdateTrainTravelDetails(ctx context.Context, travelDetailDTO *dto.TravelDetailDTO) error {
	panic("not implemented") // TODO: Implement
}

func (s *NatsRequestSender) DeleteTrainByID(ctx context.Context, trainID uint) error {
	panic("not implemented") // TODO: Implement
}

func (s *NatsRequestSender) GetSeatByID(ctx context.Context, seatID uint) (*dto.SeatDTO, error) {
	panic("not implemented") // TODO: Implement
}

func (s *NatsRequestSender) ListSeatsByTrainID(ctx context.Context, trainID uint) ([]*dto.SeatDTO, error) {
	panic("not implemented") // TODO: Implement
}

func (s *NatsRequestSender) CreateSeat(ctx context.Context, trainID uint, seatNumber uint) error {
	panic("not implemented") // TODO: Implement
}

func (s *NatsRequestSender) UpdateSeatNumberBySeatID(ctx context.Context, seatID uint, seatNumber uint) {
	panic("not implemented") // TODO: Implement
}

func (s *NatsRequestSender) DeleteSeatBySeatID(ctx context.Context, seatID uint) error {
	panic("not implemented") // TODO: Implement
}

func (s *NatsRequestSender) GetTicketByID(ctx context.Context, ticketID uint) (*dto.TicketDTO, error) {
	panic("not implemented") // TODO: Implement
}

func (s *NatsRequestSender) BookTicket(ctx context.Context, userID uint, trainID uint, TicketsNumber uint) ([]dto.TicketDTO, error) {
	panic("not implemented") // TODO: Implement
}

func (s *NatsRequestSender) CancelTicket(ctx context.Context, ticketID uint) error {
	panic("not implemented") // TODO: Implement
}
