package ports

import (
	"context"
	"user/internal/application/core/domain"
)

type UserEventPublisher interface {
	PublishUserRegistered(ctx context.Context, data *domain.User) error
	PublishUserUpdated(ctx context.Context, data *domain.User) error
}

type UserEventResponder interface {
	ReplyToGetUser(ctx context.Context, handler func(*domain.User)) error
}
