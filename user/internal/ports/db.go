package ports

import (
	"context"
	"user/internal/application/core/domain"
)

type DatabasePort interface {
	SaveUser(ctx context.Context, user *domain.User) error
	GetUserByID(ctx context.Context, ID uint) (*domain.User, error)
	ListUsers(ctx context.Context) ([]domain.User, error)
	UpdateUser(ctx context.Context, ID uint, firstName, lastName string) error
	DeleteUser(ctx context.Context, ID uint) error
}
