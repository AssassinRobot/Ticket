package postgres

import (
	"context"
	"fmt"
	"ticket/internal/application/core/domain"
	"ticket/internal/ports"
	"time"

	"gorm.io/gorm"
)

type PostgresDBAdapter struct {
	db *gorm.DB
}

func NewPostgresDBAdapter(db *gorm.DB) *PostgresDBAdapter {
	return &PostgresDBAdapter{
		db: db,
	}
}

func (r *PostgresDBAdapter) CreateTicket(ctx context.Context, ticket *domain.Ticket) (ports.Tx, error) {
	tx, err := begin(r.db)
	if err != nil {
		return nil, err
	}

	err = tx.db.WithContext(ctx).Create(ticket).Error
	if err != nil {
		return nil, fmt.Errorf("failed to create ticket: %w", err)
	}

	return tx, nil
}

func (r *PostgresDBAdapter) GetTicketByID(ctx context.Context, ID uint) (*domain.Ticket, error) {
	var ticket domain.Ticket
	err := r.db.WithContext(ctx).First(&ticket, ID).Error

	if err != nil {
		return nil, fmt.Errorf("failed to get ticket by ID: %w", err)
	}

	return &ticket, nil
}

func (r *PostgresDBAdapter) GetTicketsByUserID(ctx context.Context, userID uint) ([]domain.Ticket, error) {
	var tickets []domain.Ticket

	err := r.db.WithContext(ctx).Where("user_id = ?", userID).Find(&tickets).Error
	if err != nil {
		return nil, fmt.Errorf("failed to list tickets by user id %d: %w", userID, err)
	}

	return tickets, nil
}

func (r *PostgresDBAdapter) UpdateCanceledAt(ctx context.Context, ticketID uint, canceledAt time.Time) error {
	tx, err := begin(r.db)
	if err != nil {
		return err
	}

	err = tx.db.WithContext(ctx).Model(&domain.Ticket{}).Where("ID = ?", ticketID).Update("canceled_at", canceledAt).Error
	if err != nil {
		return fmt.Errorf("failed to update canceled_at field by ticket id %d: %w", ticketID, err)
	}

	return nil
}
