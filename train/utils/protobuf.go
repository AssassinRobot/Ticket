package utils

import (
	"train/gen"
	"train/internal/application/core/domain"

	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
)
func UnmarshalCreateSeatRequest(data []byte) (uint, uint, error) {
	protoSeat := &gen.CreateSeatRequest{}
	err := proto.Unmarshal(data, protoSeat)
	if err != nil {
		return 0, 0, err
	}

	return uint(protoSeat.TrainId), uint(protoSeat.Number), nil
}

func UnmarshalUpdateSeatNumberRequest(data []byte) (uint, uint, error) {
	protoSeat := &gen.UpdateSeatNumberRequest{}
	err := proto.Unmarshal(data, protoSeat)
	if err != nil {
		return 0, 0, err
	}

	return uint(protoSeat.ID), uint(protoSeat.Number), nil
}

func UnmarshalCreateTrainRequest(data []byte) (string, uint, error) {
	protoTrain := &gen.CreateTrainRequest{}
	err := proto.Unmarshal(data, protoTrain)
	if err != nil {
		return "", 0, err
	}

	return protoTrain.Name, uint(protoTrain.Capacity), nil
}

func UnmarshalUpdateTrainNameRequest(data []byte) (uint, string, error) {
	protoTrain := &gen.UpdateTrainNameRequest{}
	err := proto.Unmarshal(data, protoTrain)
	if err != nil {
		return 0, "", err
	}

	return uint(protoTrain.ID), protoTrain.Name, nil
}

func UnmarshalUpdateTrainTravelDetailsRequest(data []byte) (uint, *domain.TrainTravelDetails, error) {
	protoTrain := &gen.UpdateTrainTravelDetailsRequest{}
	err := proto.Unmarshal(data, protoTrain)
	if err != nil {
		return 0, nil, err
	}

	travelDetails := &domain.TrainTravelDetails{
		Origin:        protoTrain.Origin,
		Destination:   protoTrain.Destination,
		DepartureTime: protoTrain.DepartureTime.AsTime(),
		ArrivalTime:   protoTrain.ArrivalTime.AsTime(),
	}

	return uint(protoTrain.ID), travelDetails, nil
}

func MarshalListTrainsReplay(trains []domain.Train) ([]byte, error) {
	protoListTrainsReplay := &gen.ListTrainsReplay{}

	for _, train := range trains {
		protoTrain := convertTrainToProtoTrain(&train)

		protoListTrainsReplay.Trains = append(protoListTrainsReplay.Trains, protoTrain)
	}

	data, err := proto.Marshal(protoListTrainsReplay)
	if err != nil {
		return nil, err
	}
	
	return data, nil
}

func MarshalTrain(train *domain.Train) ([]byte, error) {
	protoTrain := convertTrainToProtoTrain(train)

	data, err := proto.Marshal(protoTrain)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func UnmarshalFilters(data []byte) (*domain.TrainFilters, error) {
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

func MarshalSeats(seats []domain.Seat) ([]byte, error) {
	protoSeats := &gen.ListSeatsReplay{}

	for _, seat := range seats {
		protoSeat := convertSeatToProtoSeat(&seat)
		protoSeats.Seats = append(protoSeats.Seats, protoSeat)
	}

	data, err := proto.Marshal(protoSeats)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func MarshalSeat(seat *domain.Seat) ([]byte, error) {
	protoSeat := convertSeatToProtoSeat(seat)

	data, err := proto.Marshal(protoSeat)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func UnmarshalSeat(data []byte) (*domain.Seat, error) {
	var protoSeat = gen.Seat{}

	err := proto.Unmarshal(data, &protoSeat)
	if err != nil {
		return nil, err
	}

	seat := &domain.Seat{
		ID:         uint(protoSeat.ID),
		TrainID:    uint(protoSeat.TrainId),
		UserID:     uint(protoSeat.UserId),
		SeatNumber: uint(protoSeat.SeatNumber),
	}

	return seat, nil
}

func convertTrainToProtoTrain(train *domain.Train) *gen.Train {
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
		protoSeat := convertSeatToProtoSeat(&seat)
		protoTrain.Seats = append(protoTrain.Seats, protoSeat)
	}

	return protoTrain
}

func convertSeatToProtoSeat(seat *domain.Seat) *gen.Seat {
	return &gen.Seat{
		ID:         uint32(seat.ID),
		TrainId:    uint32(seat.TrainID),
		UserId:     uint32(seat.UserID),
		SeatNumber: uint32(seat.SeatNumber),
	}
}
