package nats

import (
	"context"
	"ticket/internal/application/core/domain"
	"ticket/utils"

	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
)

const (
	SeatStream                  = "SEAT"
	SeatBookedSubjectName       = "seat.book.booked"
	SeatBookCanceledSubjectName = "seat.book.canceled"
	SeatSubjectsName            = "seat.book.*"
)

type TicketEventPublisherAdapter struct {
	jetStream jetstream.JetStream
}

func NewEventPublisherAdapter(nats *nats.Conn) (*TicketEventPublisherAdapter, error) {
	jetStream, err := jetstream.New(nats)
	if err != nil {
		return nil, err
	}

	ticketEventPublisherAdapter := &TicketEventPublisherAdapter{
		jetStream: jetStream,
	}

	return ticketEventPublisherAdapter, nil
}

func (p *TicketEventPublisherAdapter) PublishSeatBooked(ctx context.Context, seat *domain.Seat) error {
	_, err := p.jetStream.CreateOrUpdateStream(ctx, jetstream.StreamConfig{
		Name:     SeatStream,
		Subjects: []string{SeatSubjectsName},
		Storage:  jetstream.FileStorage,
	})

	if err != nil {
		return err
	}

	data, err := utils.MarshalSeat(seat)
	if err != nil {
		return err
	}

	_, err = p.jetStream.Publish(ctx, SeatBookedSubjectName, data)

	return err
}

func (p *TicketEventPublisherAdapter) PublishSeatBookingCanceled(ctx context.Context, seat *domain.Seat) error {
	_, err := p.jetStream.CreateOrUpdateStream(ctx, jetstream.StreamConfig{
		Name:     SeatStream,
		Subjects: []string{SeatSubjectsName},
		Storage:  jetstream.FileStorage,
	})

	if err != nil {
		return err
	}

	data, err := utils.MarshalSeat(seat)
	if err != nil {
		return err
	}
	
	_, err = p.jetStream.Publish(ctx, SeatBookCanceledSubjectName, data)

	return err
}
