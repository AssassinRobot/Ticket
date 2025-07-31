package api

import (
	"context"
	"errors"
	"fmt"
	"log"
	"ticket/internal/application/core/domain"
	"ticket/internal/ports"
	"time"
)

type APIAdapter struct {
	databasePort       ports.DatabasePort
	eventPublisherPort ports.EventPublisherPort
	requestPort        ports.RequestPort
}

var ErrNoAvailableTrain = errors.New("no available train")
var ErrTicketHaveAlreadyCanceled = errors.New("ticket have already canceled")

func NewAPIAdapter(dbPort ports.DatabasePort, eventPublisherPort ports.EventPublisherPort, requestPort ports.RequestPort) *APIAdapter {
	return &APIAdapter{
		databasePort:       dbPort,
		eventPublisherPort: eventPublisherPort,
		requestPort:        requestPort,
	}
}

func (a *APIAdapter) GetTicketByID(ctx context.Context, ticketID uint) (*domain.Ticket, error) {
	return a.databasePort.GetTicketByID(ctx, ticketID)
}

func (a *APIAdapter) GetTicketsByUserID(ctx context.Context, userID uint) ([]domain.Ticket, error) {
	return a.databasePort.GetTicketsByUserID(ctx, userID)
}

func (a *APIAdapter) BookTicket(ctx context.Context, userID, trainID, ticketNumber uint) ([]domain.Ticket, error) {
	availableTrain, err := a.requestPort.RequestGetTrainByID(ctx, trainID)
	if err != nil {
		return nil, err
	}

	trainSeatsLen := len(availableTrain.Seats)

	if trainSeatsLen < int(ticketNumber) {
		return nil, fmt.Errorf("%w for %d tickets", ErrNoAvailableTrain, ticketNumber)
	}

	bookedTickets := []domain.Ticket{}

	if ticketNumber == 1 {
		selectedSeat := availableTrain.Seats[0]

		selectedSeat.UserID = userID

		ticket := domain.NewTicket(userID,
			availableTrain.ID, selectedSeat.ID,
			selectedSeat.SeatNumber,
			availableTrain.DepartureTime,
			availableTrain.TravelDetails)

		tx, err := a.databasePort.CreateTicket(ctx, ticket)
		if err != nil {
			return nil, err
		}

		err = a.eventPublisherPort.PublishSeatBooked(ctx, &selectedSeat)
		if err != nil {
			rollBackError := tx.Rollback()
			log.Printf("roll back error:%v\n", rollBackError)

			return nil, err
		}

		bookedTickets = append(bookedTickets, *ticket)

		err = tx.Commit()
		if err != nil {
			return nil, err
		}

	} else if ticketNumber > 1 {
		for i := 0; i < trainSeatsLen; i++ {
			selectedSeat := availableTrain.Seats[i]

			selectedSeat.UserID = userID

			ticket := domain.NewTicket(userID,
				availableTrain.ID, selectedSeat.ID,
				selectedSeat.SeatNumber,
				availableTrain.DepartureTime,
				availableTrain.TravelDetails)

			tx, err := a.databasePort.CreateTicket(ctx, ticket)
			if err != nil {
				return nil, err
			}

			err = a.eventPublisherPort.PublishSeatBooked(ctx, &selectedSeat)
			if err != nil {
				rollBackError := tx.Rollback()
				log.Printf("roll back error:%v\n", rollBackError)

				return nil, err
			}

			bookedTickets = append(bookedTickets, *ticket)

			err = tx.Commit()
			if err != nil {
				return nil, err
			}
		}
	}

	return bookedTickets, nil
}

func (a *APIAdapter) CancelTicket(ctx context.Context, ticketID uint) error {
	ticket, err := a.databasePort.GetTicketByID(ctx, ticketID)
	if err != nil {
		return err
	}

	if ticket.CanceledAt != nil {
		return ErrTicketHaveAlreadyCanceled
	}

	seat := &domain.Seat{
		ID:         ticket.SeatID,
		TrainID:    ticket.TrainID,
		SeatNumber: ticket.SeatNumber,
	}

	err = a.eventPublisherPort.PublishSeatBookingCanceled(ctx, seat)
	if err != nil {
		return err
	}

	err = a.databasePort.UpdateCanceledAt(ctx, ticketID, time.Now())

	return err
}
