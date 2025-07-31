package utils

import (
	"train/gen"
	"train/internal/application/core/domain"

	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func MarshalTrains(trains []domain.Train) ([]byte, error) {
	bytes := []byte{}
	for _, train := range trains {
		data, err := MarshalTrain(&train)
		if err != nil {
			return nil, err
		}

		bytes = append(bytes, data...)
	}

	return bytes, nil
}

func MarshalTrain(train *domain.Train) ([]byte, error) {
	protoTrain := &gen.Train{
		ID:            uint32(train.ID),
		Name:          train.Name,
		Origin:        train.Origin,
		Destination:   train.Destination,
		Capacity:      train.Capacity,
		IsFull:        train.IsFull,
		DepartureTime: timestamppb.New(train.DepartureTime),
		ArrivalTime:   timestamppb.New(train.ArrivalTime),
	}

	for _, seat := range train.Seats {
		protoSeat := &gen.Seat{
			ID:         uint32(seat.ID),
			TrainId:    uint32(train.ID),
			SeatNumber: uint32(seat.ID),
		}

		protoTrain.Seats = append(protoTrain.Seats, protoSeat)
	}

	data, err := proto.Marshal(protoTrain)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func UnMarshalFilters(data []byte) (*domain.TrainFilters, error) {
	protoTrainFilter := gen.TrainFilter{}

	err := proto.Unmarshal(data, &protoTrainFilter)
	if err != nil {
		return nil, err
	}

	trainFilter := domain.NewFilterTrains(
		protoTrainFilter.Name,
		protoTrainFilter.Origin,
		protoTrainFilter.Destination,
		protoTrainFilter.DepartureTime.AsTime(),
		protoTrainFilter.ArrivalTime.AsTime(),
	)
	
	return trainFilter, nil
}

func UnmarshalSeat(data []byte) (*domain.Seat,error){
	var protoSeat = gen.Seat{} 

	err := proto.Unmarshal(data,&protoSeat)
	if err != nil {
		return nil, err
	}

	seat := &domain.Seat{
		ID:         uint(protoSeat.ID),
		TrainID:    uint(protoSeat.TrainId),
		UserID: uint(protoSeat.UserId),
		SeatNumber: uint(protoSeat.SeatNumber),
	}

	return seat, nil
}