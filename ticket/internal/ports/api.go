package ports

import (
	"context"
	"ticket/internal/application/core/domain"
)

type APIPort interface {
	GetTicketByID(ctx context.Context, ticketID uint) (*domain.Ticket, error)
	GetTicketsByUserID(ctx context.Context, userID uint) ([]domain.Ticket, error)
	BookTicket(ctx context.Context, userID,trainID,ticketNumber uint) ([]domain.Ticket, error)
	CancelTicket(ctx context.Context, ticketID uint) error
}
