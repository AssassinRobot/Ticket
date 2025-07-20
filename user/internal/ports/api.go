package ports

import (
	"context"
	"user/internal/application/core/domain"
)

type APIPort interface {
	Register(ctx context.Context, firstName, lastName, email string) (*domain.User, error)
	GetUserByID(ctx context.Context, id uint) (*domain.User, error)
	ListUsers(ctx context.Context) ([]domain.User, error)
	UpdateUser(ctx context.Context, id uint, firstName, lastName string) error
	DeleteUser(ctx context.Context, id uint) error
}
