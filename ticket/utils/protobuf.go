package utils

import (
	"ticket/gen"
	"ticket/internal/application/core/domain"

	"google.golang.org/protobuf/proto"
)

func UnMarshalTrain(data []byte) (*domain.Train, error) {
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
		UserId: uint32(seat.UserID),
	}

	data, err := proto.Marshal(&protoSeat)
	if err != nil {
		return nil, err
	}

	return data, nil
}
