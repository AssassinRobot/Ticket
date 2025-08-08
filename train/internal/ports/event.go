package ports

import (
	"context"
)

type TrainEventResponder interface {
	ReplyToListTrains(ctx context.Context) error
	ReplyToListTrainsFiltered(ctx context.Context) error
	ReplyToGetTrainByID(ctx context.Context) error
	ReplyToCreateTrain(ctx context.Context) error
	ReplyToUpdateTrainName(ctx context.Context) error
	ReplyToUpdateTrainTravelDetails(ctx context.Context) error
	ReplyToDeleteTrainByID(ctx context.Context) error
	ReplyToListSeatsByTrainID(ctx context.Context) error

	ReplyToGetSeatByID(ctx context.Context) error
	ReplyToCreateSeat(ctx context.Context) error
	ReplyToUpdateSeatNumberBySeatID(ctx context.Context) error
	ReplyToDeleteSeatBySeatID(ctx context.Context) error
}

type SeatEventsConsumer interface {
	ConsumerSeatEvents(ctx context.Context) <-chan error
}
