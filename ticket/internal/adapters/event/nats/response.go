package nats

import (
	"context"
	"fmt"
	"strconv"
	"ticket/internal/ports"
	"ticket/utils"

	"github.com/nats-io/nats.go"
)

const (
	SubjectRequestGetTicketByID        = "request.get.ticket.byID"
	SubjectRequestBookTicket           = "request.book.ticket"
	SubjectRequestCancelTicket         = "request.cancel.ticket"
	SubjectRequestListTicketsByUserID  = "request.list.tickets.byUserID"
	SubjectRequestListTicketsByTrainID = "request.list.tickets.byTrainID"
)

type TicketEventResponderAdapter struct {
	natsConn   *nats.Conn
	APIAdapter ports.APIPort
}

func NewTicketEventResponderAdapter(natsConn *nats.Conn, apiAdapter ports.APIPort) *TicketEventResponderAdapter {
	return &TicketEventResponderAdapter{
		natsConn:   natsConn,
		APIAdapter: apiAdapter,
	}
}

func (r *TicketEventResponderAdapter) ReplyToGetTicketByID(ctx context.Context) error {
	subscription, err := r.natsConn.Subscribe(SubjectRequestGetTicketByID, func(msg *nats.Msg) {
		ticketIDstr := string(msg.Data)
		ticketID, err := strconv.Atoi(ticketIDstr)
		if err != nil {
			errString := fmt.Errorf("error invalid ticketID: %d, %w", ticketID, err).Error()

			msg.Respond([]byte(errString))
			return
		}

		ticket, err := r.APIAdapter.GetTicketByID(ctx, uint(ticketID))
		if err != nil {
			errString := fmt.Errorf("error list tickets: %w", err).Error()
			msg.Respond([]byte(errString))
			return
		}

		serializedTicketData, err := utils.MarshalTicket(ticket)
		if err != nil {
			fmt.Println()
			errString := fmt.Errorf("error serialize ticket: %v,%w", ticket, err).Error()
			msg.Respond([]byte(errString))
			return
		}

		msg.Respond(serializedTicketData)
	})

	if err != nil {
		return err
	}

	go func() {
		<-ctx.Done()
		subscription.Unsubscribe()
	}()

	return nil
}

func (r *TicketEventResponderAdapter) ReplayToListTicketsByUserID(ctx context.Context) error {
	subscription, err := r.natsConn.Subscribe(SubjectRequestListTicketsByUserID, func(msg *nats.Msg) {
		userIDstr := string(msg.Data)
		userID, err := strconv.Atoi(userIDstr)
		if err != nil {
			errString := fmt.Errorf("error invalid userID: %d, %w", userID, err).Error()

			msg.Respond([]byte(errString))
			return
		}

		tickets, err := r.APIAdapter.GetTicketsByUserID(ctx, uint(userID))
		if err != nil {
			errString := fmt.Errorf("error list tickets: %w", err).Error()
			msg.Respond([]byte(errString))
			return
		}

		serializedTicketsData, err := utils.MarshalTickets(tickets)
		if err != nil {
			fmt.Println()
			errString := fmt.Errorf("error serialize tickets: %v,%w", tickets, err).Error()
			msg.Respond([]byte(errString))
			return
		}

		msg.Respond(serializedTicketsData)
	})

	if err != nil {
		return err
	}

	go func() {
		<-ctx.Done()
		subscription.Unsubscribe()
	}()

	return nil
}

func (r *TicketEventResponderAdapter) ReplayToListTicketsByTrainID(ctx context.Context) error {
	subscription, err := r.natsConn.Subscribe(SubjectRequestListTicketsByTrainID, func(msg *nats.Msg) {
		trainIDstr := string(msg.Data)
		trainID, err := strconv.Atoi(trainIDstr)
		if err != nil {
			errString := fmt.Errorf("error invalid trainID: %d, %w", trainID, err).Error()

			msg.Respond([]byte(errString))
			return
		}

		tickets, err := r.APIAdapter.GetTicketsByTrainID(ctx, uint(trainID))
		if err != nil {
			errString := fmt.Errorf("error list tickets: %w", err).Error()
			msg.Respond([]byte(errString))
			return
		}

		serializedTicketsData, err := utils.MarshalTickets(tickets)
		if err != nil {
			fmt.Println()
			errString := fmt.Errorf("error serialize tickets: %v,%w", tickets, err).Error()
			msg.Respond([]byte(errString))
			return
		}

		msg.Respond(serializedTicketsData)
	})

	if err != nil {
		return err
	}

	go func() {
		<-ctx.Done()
		subscription.Unsubscribe()
	}()

	return nil
}

func (r *TicketEventResponderAdapter) ReplayToBookTicket(ctx context.Context) error {
	subscription, err := r.natsConn.Subscribe(SubjectRequestBookTicket, func(msg *nats.Msg) {
		userID, trainID, ticketNumber, err := utils.UnmarshalBookTicketRequest(msg.Data)

		if err != nil {
			errString := fmt.Errorf("error unmarshal request(BookTicket):%w", err).Error()
			msg.Respond([]byte(errString))
			return
		}

		tickets, err := r.APIAdapter.BookTicket(ctx, userID, trainID, ticketNumber)
		if err != nil {
			errString := fmt.Errorf("error book ticket:, %w", err).Error()
			msg.Respond([]byte(errString))
			return
		}

		serializedTicketsData, err := utils.MarshalTickets(tickets)
		if err != nil {
			errString := fmt.Errorf("error unmarshal request(BookTicket):%w", err).Error()
			msg.Respond([]byte(errString))
			return
		}

		msg.Respond(serializedTicketsData)
	})

	if err != nil {
		return err
	}

	go func() {
		<-ctx.Done()
		subscription.Unsubscribe()
	}()

	return nil
}

func (r *TicketEventResponderAdapter) ReplayToCancelTicket(ctx context.Context) error {
	subscription, err := r.natsConn.Subscribe(SubjectRequestCancelTicket, func(msg *nats.Msg) {
		ticketIDstr := string(msg.Data)
		ticketID, err := strconv.Atoi(ticketIDstr)
		if err != nil {
			errString := fmt.Errorf("error invalid ticketID: %d", ticketID).Error()

			msg.Respond([]byte(errString))
			return
		}

		err = r.APIAdapter.CancelTicket(ctx, uint(ticketID))
		if err != nil {
			errString := fmt.Errorf("error cancel ticket:%d, %w", ticketID, err).Error()
			msg.Respond([]byte(errString))
			return
		}

		msg.Respond([]byte("ticket canceled successfully"))
	})

	if err != nil {
		return err
	}

	go func() {
		<-ctx.Done()
		subscription.Unsubscribe()
	}()

	return nil
}
