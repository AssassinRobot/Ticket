package nats

import (
	"context"
	"errors"
	"fmt"
	"train/utils"

	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
)

type NatsEventConsumer struct {
	jetStream jetstream.JetStream
	SeatHandlers
}

type SeatHandlers interface {
	SeatBooked(ctx context.Context, seatID, trainID uint) error
	CancelSeatBooking(ctx context.Context, seatID, trainID uint) error
}

const (
	SeatStream                  = "SEAT"
	SeatBookedSubjectName       = "seat.booked"
	SeatBookCanceledSubjectName = "seat.book.canceled"
)

var (
	ErrInvalidSubject = errors.New("invalid seat subject")
)

func NewNatsEventConsumer(nats *nats.Conn, seatHandlers SeatHandlers) (*NatsEventConsumer, error) {
	jetStream, err := jetstream.New(nats)
	if err != nil {
		return nil, err
	}

	natsEventConsumerInstance := &NatsEventConsumer{
		jetStream: jetStream,
		SeatHandlers: seatHandlers,
	}

	return natsEventConsumerInstance, nil
}

func (c *NatsEventConsumer) ConsumerSeatEvents(ctx context.Context) <-chan error {
	var errorStream = make(chan error)

	seatConsumer, err := c.jetStream.CreateOrUpdateConsumer(ctx, SeatStream, jetstream.ConsumerConfig{
		Durable:       "seat_book_consumer",
		MaxAckPending: 5,
	})

	if err != nil {
		errorStream <- err

		close(errorStream)

		return errorStream
	}

	go func() {
		defer close(errorStream)

		for {
			msgs, err := seatConsumer.Fetch(5)
			if err != nil {
				errorStream <- err
			}

			for msg := range msgs.Messages() {
				switch msg.Subject() {
				case SeatBookedSubjectName:
					seat, err := utils.UnmarshalSeat(msg.Data())
					if err != nil {
						errorStream <- err
						ack(msg, errorStream)
						continue
					}

					err = c.SeatHandlers.SeatBooked(ctx, seat.ID, seat.TrainID)
					if err != nil {
						errorStream <- err
						ack(msg, errorStream)
						continue
					}

					ack(msg, errorStream)

				case SeatBookCanceledSubjectName:
					seat, err := utils.UnmarshalSeat(msg.Data())
					if err != nil {
						errorStream <- err
						ack(msg, errorStream)
						continue
					}

					err = c.CancelSeatBooking(ctx, seat.ID, seat.TrainID)
					if err != nil {
						errorStream <- err
						ack(msg, errorStream)
						continue
					}

					ack(msg, errorStream)
				default:
					errorStream <- fmt.Errorf("%w:%s", ErrInvalidSubject, msg.Subject())
				}

			}

		}
	}()

	return errorStream
}

func ack(msg jetstream.Msg, errorStream chan<- error) {
	err := msg.Ack()
	if err != nil {
		errorStream <- err
	}
}
