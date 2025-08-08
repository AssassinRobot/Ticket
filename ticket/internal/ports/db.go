package ports

import (
	"context"
	"ticket/internal/application/core/domain"
	"time"
)

type DatabasePort interface {
	GetTicketByID(ctx context.Context,ticketID uint) (*domain.Ticket, error)
	GetTicketsByUserID(ctx context.Context,userID uint) ([]domain.Ticket, error)
	GetTicketsByTrainID(ctx context.Context,trainID uint) ([]domain.Ticket, error)
	CreateTicket(ctx context.Context,ticket *domain.Ticket) (Tx,error)
	UpdateCanceledAt(ctx context.Context,ticketID uint, canceledAt time.Time) error
}

type Tx interface {
	Commit() error
	Rollback() error
}