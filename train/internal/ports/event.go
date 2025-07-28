package ports

import (
	"context"
)

type TrainEventResponder interface {
	ReplyToListTrains(ctx context.Context) error
	ReplyToGetTrainByID(ctx context.Context) error
	ReplyToListTrainsFiltered(ctx context.Context) error
}

type SeatEventsConsumer interface {
	ConsumerSeatEvents(ctx context.Context) <-chan error
}
