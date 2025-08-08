package utils

import (
	"ticket/gen"
	"ticket/internal/application/core/domain"

	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
)


func UnmarshalBookTicketRequest(data []byte)(uint,uint,uint,error){
	protoBookTicketRequest := &gen.BookTicketRequest{}

	err := proto.Unmarshal(data,protoBookTicketRequest)
	if err != nil {
		return 0,0,0, err
	}

	return uint(protoBookTicketRequest.UserId),uint(protoBookTicketRequest.TrainId),uint(protoBookTicketRequest.TicketNumber),nil
}

func UnmarshalTrain(data []byte) (*domain.Train, error) {
	protoTrain := gen.Train{}
	err := proto.Unmarshal(data, &protoTrain)
	if err != nil {
		return nil, err
	}

	train := &domain.Train{
		ID:   uint(protoTrain.ID),
		Name: protoTrain.Name,
		TravelDetails: domain.TravelDetails{
			Origin:        protoTrain.Origin,
			Destination:   protoTrain.Destination,
			DepartureTime: protoTrain.DepartureTime.AsTime(),
			ArrivalTime:   protoTrain.ArrivalTime.AsTime(),
		},
	}

	for _, protoSeat := range protoTrain.Seats {
		seat := domain.Seat{
			ID:         uint(protoSeat.ID),
			TrainID:    uint(protoSeat.TrainId),
			SeatNumber: uint(protoSeat.SeatNumber),
		}

		train.Seats = append(train.Seats, seat)
	}

	return train, nil
}

func MarshalSeat(seat *domain.Seat) ([]byte, error) {
	var protoSeat = gen.Seat{
		ID:         uint32(seat.ID),
		TrainId:    uint32(seat.TrainID),
		SeatNumber: uint32(seat.SeatNumber),
		UserId:     uint32(seat.UserID),
	}

	data, err := proto.Marshal(&protoSeat)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func MarshalTickets(tickets []domain.Ticket) ([]byte, error) {
	protoListTickets := &gen.ListTickets{}

	for _, ticket := range tickets {
		protoTicket := convertTicketToProtoTicket(&ticket)

		protoListTickets.Tickets = append(protoListTickets.Tickets, protoTicket)
	}

	data, err := proto.Marshal(protoListTickets)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func MarshalTicket(ticket *domain.Ticket) ([]byte, error) {
	protoTicket := convertTicketToProtoTicket(ticket)

	data, err := proto.Marshal(protoTicket)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func convertTicketToProtoTicket(ticket *domain.Ticket) *gen.Ticket {
	protoTicket := &gen.Ticket{
		ID:            uint32(ticket.ID),
		UserId:        uint32(ticket.UserID),
		TrainId:       uint32(ticket.TrainID),
		SeatNumber:    uint32(ticket.SeatNumber),
		Origin:        ticket.Origin,
		Destination:   ticket.Destination,
		DepartureTime: timestamppb.New(ticket.DepartureTime),
		ArrivalTime:   timestamppb.New(ticket.ArrivalTime),
		ExpiresAt:     timestamppb.New(ticket.ExpiresAt),
	}

	if ticket.CanceledAt != nil {
		protoTicket.CanceledAt = timestamppb.New(*ticket.CanceledAt)
	}

	return protoTicket
}
